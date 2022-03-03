package main

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/thanhpp/discordbots/internal/discordbot"
	"github.com/thanhpp/discordbots/pkg/logger"
	"github.com/thanhpp/discordbots/pkg/timehelper"
	"github.com/thanhpp/discordbots/pkg/weathercrawler"
)

var (
	enableDiscordLog bool = true
)

func main() {
	// read config
	botCfg, err := discordbot.NewConfigFromFile("discordbot-secret.yml")
	if err != nil {
		panic(errors.WithMessage(err, "Read discord bot config"))
	}

	// setup
	logger.Set(&logger.LogConfig{
		Color:      true,
		LoggerName: "thanhpp's discordbot",
		Level:      "DEBUG",
	}, enableDiscordLog)

	bot, err := discordbot.NewBot("thanhpp", botCfg.BotToken)
	if err != nil {
		panic(errors.WithMessage(err, "Create a new discord bot"))
	}

	// adding alert
	addWeatherAlert(bot, botCfg.WeatherChannel)
	if enableDiscordLog {
		addLogAlert(bot, botCfg.LogChannel)
	}

	// start the bot
	logger.Get().Infof("Bot %s is starting...", bot.Name())
	if err := bot.Start(); err != nil {
		logger.Get().Fatalf("%+v", errors.WithMessage(err, "Start a bot"))
	}
	logger.Get().Infof("Bot %s stopped", bot.Name())
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
				testTimeCond := timehelper.MustNewTimeCond(0, 0, 0)

				waitDur, err := timehelper.ShortestUntilNextTimeCond(time.Now(), timeCond1, timeCond2, timeCond3, timeCond4, testTimeCond)
				if err != nil {
					panic(err)
				}

				logger.Get().Infof("[Weather alert] Wait %.2f mins to start", waitDur.Minutes())
				<-time.After(waitDur)
				logger.Get().Infof("[Weather alert] Start")

				var (
					toSend       = false
					toSendTicker = time.NewTicker(time.Hour * 6)
					crawlTicker  = time.NewTicker(time.Minute * 10)
					weatherCrw   = weathercrawler.NewWeatherCrawler()
				)

				// first message, due to the ticker starts after its interval
				msg, _ := parseWeatherMsg(weatherCrw)
				msg.ChannelID = channelID
				logger.Get().Debugf("[Weather alert] Send message: %+v", msg)
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
						msg, ok := parseWeatherMsg(weatherCrw)
						if !ok {
							logger.Get().Warn("[Weather alert] Failed to parse weather info")
							continue
						}
						msg.ChannelID = channelID
						logger.Get().Debugf("[Weather alert] Send message: %+v", msg)
						returnC <- msg
						toSend = false
					}
				}
			}()

			return returnC
		},
	)
}

func parseWeatherMsg(crwl *weathercrawler.WeatherCrawler) (*discordbot.Message, bool) {
	logger.Get().Debugf("[DEBUG] Crawling weather....")
	weatherNow := crwl.GetInfoNow()
	if len(weatherNow.Temperature) == 0 {
		log.Println("[ERROR] Empty temparature - 1")
		return nil, false
	}

	if weatherNow.Temperature == "°C" {
		logger.Get().Warn("[ERROR] Empty temparature - 2")
		return nil, false
	}

	msg := new(discordbot.Message)
	msg.Topic = "HANOI WEATHER INFO"
	msg.AddContent("TMP", weatherNow.Temperature)
	msg.AddContent("STT", weatherNow.Status)
	msg.AddContent("HMD", weatherNow.Humidity)
	msg.AddContent("UPD", weatherNow.LastUpdated)

	return msg, true
}

func addLogAlert(bot *discordbot.Bot, channelID string) {
	bot.NewAlert(
		"Log Alert",
		channelID,
		func(alertCtx context.Context) chan *discordbot.Message {
			var logC = make(chan *discordbot.Message)

			go func() {
				for {
					select {
					case logMsg := <-logger.Get().OutputC():
						// create new log message
						msg := new(discordbot.Message)
						msg.ChannelID = channelID
						msg.Topic = "LOG"
						msg.AddContent("Content", logMsg)
						logC <- msg

					case <-alertCtx.Done():
						return
					}
				}
			}()

			return logC
		},
	)
}
