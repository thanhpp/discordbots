package main

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/thanhpp/discordbots/internal/discordbot"
	"github.com/thanhpp/discordbots/pkg/timehelper"
	"github.com/thanhpp/discordbots/pkg/weathercrawler"
)

func main() {
	// read config
	botCfg, err := discordbot.NewConfigFromFile("discordbot-secret.yml")
	if err != nil {
		panic(errors.WithMessage(err, "Read discord bot config"))
	}

	// setup
	bot, err := discordbot.NewBot("thanhpp's bot", botCfg.BotToken)
	if err != nil {
		panic(errors.WithMessage(err, "Create a new discord bot"))
	}
	addWeatherAlert(bot, botCfg.WeatherChannel)

	// start the bot
	log.Printf("Bot %s is starting... \n", bot.Name())
	if err := bot.Start(); err != nil {
		panic(errors.WithMessage(err, "Start a bot"))
	}
	log.Printf("Bot %s stopped \n", bot.Name())
}

func addWeatherAlert(bot *discordbot.Bot, channelID string) {
	bot.NewAlert(
		"Weather Alert",
		channelID,
		func(alertCtx context.Context) chan *discordbot.Message {
			returnC := make(chan *discordbot.Message)

			go func() {
				// wait
				timeCond1, err := timehelper.NewTimeCond(0, 0, 0)
				if err != nil {
					panic(err)
				}

				timeCond2, err := timehelper.NewTimeCond(6, 0, 0)
				if err != nil {
					panic(err)
				}

				timeCond3, err := timehelper.NewTimeCond(12, 0, 0)
				if err != nil {
					panic(err)
				}

				timeCond4, err := timehelper.NewTimeCond(18, 0, 0)
				if err != nil {
					panic(err)
				}

				waitDur, err := timehelper.ShortestUntilNextTimeCond(time.Now(), timeCond1, timeCond2, timeCond3, timeCond4)
				if err != nil {
					panic(err)
				}

				<-time.After(waitDur)

				var (
					toSend       = false
					toSendTicker = time.NewTicker(time.Hour * 6)
					crawlTicker  = time.NewTicker(time.Minute * 5)
					weatherCrw   = weathercrawler.NewWeatherCrawler()
				)

				for {
					select {
					case <-alertCtx.Done():
						toSendTicker.Stop()
						crawlTicker.Stop()
						close(returnC)
						return

					case <-toSendTicker.C:
						toSend = true

					case <-crawlTicker.C:
						if !toSend {
							continue
						}
						weatherNow := weatherCrw.GetInfoNow()
						msg := new(discordbot.Message)
						msg.ChannelID = channelID
						msg.Topic = "HANOI WEATHER INFO"
						msg.AddContent("TMP", weatherNow.Temperature)
						msg.AddContent("STT", weatherNow.Status)
						msg.AddContent("HMD", weatherNow.Humidity)
						msg.AddContent("UPD", weatherNow.LastUpdated)
						toSend = true
					}
				}
			}()

			return returnC
		},
	)
}
