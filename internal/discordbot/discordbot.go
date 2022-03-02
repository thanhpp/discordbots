package discordbot

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/thanhpp/discordbots/pkg/logger"
)

type Bot struct {
	name          string
	token         string
	session       *discordgo.Session
	mainCtx       context.Context
	mainCtxCancel context.CancelFunc
	receiveChan   chan *Message
	alerts        []*Alert
}

func NewBot(name, token string) (*Bot, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	mainCtx, mainCtxCancel := context.WithCancel(context.Background())

	bot := &Bot{
		name:          name,
		token:         token,
		session:       s,
		mainCtx:       mainCtx,
		mainCtxCancel: mainCtxCancel,
		receiveChan:   make(chan *Message),
	}

	return bot, nil
}

func (b Bot) Name() string {
	return b.name
}

func (b *Bot) Start() error {
	err := b.session.Open()
	if err != nil {
		return err
	}

	for i := range b.alerts {
		go b.alerts[i].Start()
	}

	go func() {
		for {
			select {
			case msg := <-b.receiveChan:
				if err := b.sendMessage(msg.ChannelID, msg.Stringtify()); err != nil {
					logger.Get().Errorf("err sending alert message", err, msg)
					continue
				}
				// logger.Get().Debug("Sent msg from alert") - Otherwise, the discord log channel will be flooded with duplicate messages

			case <-b.mainCtx.Done():
				logger.Get().Info("Stop receiving alert message")
				return
			}
		}
	}()

	defer func(er *error) {
		logger.Get().Infof("Bot stopping....")
		b.mainCtxCancel()
		*er = b.session.Close()
		close(b.receiveChan)
		logger.Get().Debugf("Bot stopped")
	}(&err)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	return err
}

func (b *Bot) AddHandlers(handlers ...interface{}) {
	for i := range handlers {
		if handlers[i] == nil {
			continue
		}

		b.session.AddHandler(handlers[i])
	}
}
