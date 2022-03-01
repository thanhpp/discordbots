package discordbot

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// HandleNewMessage is called every time a new message is received.
func (b *Bot) HandleNewMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	fmt.Println(m.Content)
	fmt.Println(m.ChannelID)

	if err := b.sendMessage(m.ChannelID, "Hello World!"); err != nil {
		log.Println(err)
	}
}

func (b *Bot) sendMessage(channelID, message string) error {
	_, err := b.session.ChannelMessageSend(channelID, message)
	return err
}
