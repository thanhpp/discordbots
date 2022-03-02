package discordbot

import (
	"context"
	"sync"

	"github.com/thanhpp/discordbots/pkg/logger"
)

type AlertMsgGenerator func(alertCtx context.Context) chan *Message

type Alert struct {
	lock        sync.Mutex
	name        string
	channelID   string
	msgGen      AlertMsgGenerator
	ctx         context.Context
	cancelFn    context.CancelFunc
	isCanceled  bool
	receiveChan chan<- *Message
}

func (b *Bot) NewAlert(name, channelID string, msgGen AlertMsgGenerator) {
	newCtx, cancel := context.WithCancel(b.mainCtx)
	b.alerts = append(b.alerts, &Alert{
		name:        name,
		channelID:   channelID,
		msgGen:      msgGen,
		ctx:         newCtx,
		cancelFn:    cancel,
		isCanceled:  false,
		receiveChan: b.receiveChan,
	})
}

func (a *Alert) Cancel() {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.isCanceled {
		return
	}

	a.cancelFn()
	a.isCanceled = true
}

func (a *Alert) Ctx() context.Context {
	return a.ctx
}

func (a *Alert) Start() {
	msgC := a.msgGen(a.ctx)
	for {
		select {
		case <-a.ctx.Done():
			logger.Get().Infof("[Alert] %s is stopped", a.name)
			a.Cancel()
			return
		case msg := <-msgC:
			a.receiveChan <- msg
		}
	}

}
