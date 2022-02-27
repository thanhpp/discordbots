package discordbot

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/thanhpp/discordbots/pkg/weathercrawler"
)

type Bot struct {
	name    string
	token   string
	session *discordgo.Session
	mainCtx context.Context
}

func NewBot(name, token string) (*Bot, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		name:    name,
		token:   token,
		session: s,
		mainCtx: context.Background(),
	}

	return bot, nil
}

func (b *Bot) Start() error {
	err := b.session.Open()
	if err != nil {
		return err
	}

	weatherCrw := weathercrawler.NewWeatherCrawler()
	weatherInfo := weatherCrw.GetInfoNow()
	msg := new(Message)
	msg.Topic = "[HANOI WEATHER INFO]"
	msg.AddContent("TMP", weatherInfo.Temperature)
	msg.AddContent("STA", weatherInfo.Status)
	msg.AddContent("HMD", weatherInfo.Humidity)
	msg.AddContent("UPD", weatherInfo.LastUpdated)
	b.sendMessage("947185195575021669", msg.Stringtify())

	defer func(er *error) {
		*er = b.session.Close()
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
