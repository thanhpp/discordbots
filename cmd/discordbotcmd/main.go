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
	log.Printf("Bot %s stopped \n\n\n", bot.Name())
}

func addWeatherAlert(bot *discordbot.Bot, channelID string) {
	bot.NewAlert(
		"Weather Alert",
		channelID,
		func(alertCtx context.Context) chan *discordbot.Message {
			returnC := make(chan *discordbot.Message)

			go func() {
				// wait
				timeCond1 := timehelper.MustNewTimeCond(0, 0, 0)
				timeCond2 := timehelper.MustNewTimeCond(6, 0, 0)
				timeCond3 := timehelper.MustNewTimeCond(12, 0, 0)
				timeCond4 := timehelper.MustNewTimeCond(18, 0, 0)

				waitDur, err := timehelper.ShortestUntilNextTimeCond(time.Now(), timeCond1, timeCond2, timeCond3, timeCond4)
				if err != nil {
					panic(err)
				}

				log.Printf("[Weather alert] Wait %f mins to start \n", waitDur.Minutes())
				<-time.After(waitDur)
				log.Printf("[Weather alert] Start\n")

				var (
					toSend       = false
					toSendTicker = time.NewTicker(time.Hour * 6)
					crawlTicker  = time.NewTicker(time.Minute * 10)
					weatherCrw   = weathercrawler.NewWeatherCrawler()
				)

				// first message, due to the ticker starts after its interval
				msg := parseWeatherMsg(weatherCrw)
				msg.ChannelID = channelID
				log.Printf("[Weather alert] Send message: %+v \n", msg)
				returnC <- msg

				for {
					select {
					case <-alertCtx.Done():
						toSendTicker.Stop()
						crawlTicker.Stop()
						close(returnC)
						return

					case <-toSendTicker.C:
						log.Println("[DEBUG] Update toSend ticker")
						toSend = true

					case <-crawlTicker.C:
						if !toSend {
							continue
						}
						msg := parseWeatherMsg(weatherCrw)
						msg.ChannelID = channelID
						log.Printf("[Weather alert] Send message: %+v \n", msg)
						returnC <- msg
						toSend = false
					}
				}
			}()

			return returnC
		},
	)
}

func parseWeatherMsg(crwl *weathercrawler.WeatherCrawler) *discordbot.Message {
	log.Println("[DEBUG] Crawling weather....")
	weatherNow := crwl.GetInfoNow()
	msg := new(discordbot.Message)
	msg.Topic = "HANOI WEATHER INFO"
	msg.AddContent("TMP", weatherNow.Temperature)
	msg.AddContent("STT", weatherNow.Status)
	msg.AddContent("HMD", weatherNow.Humidity)
	msg.AddContent("UPD", weatherNow.LastUpdated)

	return msg
}
