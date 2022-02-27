package discordbot_test

import (
	"testing"

	"github.com/thanhpp/discordbots/internal/discordbot"
)

func NewBot() (*discordbot.Bot, error) {
	var (
		configPath = "./discordbot-secret.yml"
		botName    = "test-bot"
	)

	botCfg, err := discordbot.NewConfigFromFile(configPath)
	if err != nil {
		return nil, err
	}

	bot, err := discordbot.NewBot(botName, botCfg.BotToken)
	if err != nil {
		return nil, err
	}

	return bot, nil
}

func TestSendMessageToChannel(t *testing.T) {
	bot, err := NewBot()
	if err != nil {
		t.Fatal(err)
	}

	bot.AddHandlers(bot.HandleNewMessage)

	if err := bot.Start(); err != nil {
		t.Fatal(err)
	}
}
