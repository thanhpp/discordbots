package discordbot

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type BotConfig struct {
	AppID          string `mapstructure:"AppID"`
	PublicKey      string `mapstructure:"PublicKey"`
	BotToken       string `mapstructure:"BotToken"`
	WeatherChannel string `mapstructure:"WeatherChannel"`
	LogChannel     string `mapstructure:"LogChannel"`
}

func NewConfigFromFile(p string) (*BotConfig, error) {
	file := filepath.Base(p)

	// 0: filename, 1: file ext
	fileSl := strings.Split(file, ".")
	if len(fileSl) != 2 {
		return nil, errors.Errorf("invalid config file: %s", file)
	}

	// read the config file using viper
	v := viper.New()
	v.SetConfigName(fileSl[0])
	v.SetConfigType(fileSl[1])
	v.AddConfigPath(".")     // current folder
	v.AddConfigPath("..")    // depth 1 - test file
	v.AddConfigPath("../..") // depth 2 - test file

	if err := v.ReadInConfig(); err != nil {
		return nil, errors.WithMessage(err, "read config")
	}

	c := new(BotConfig)
	if err := v.Unmarshal(c); err != nil {
		return nil, errors.WithMessage(err, "unmarshal")
	}

	return c, nil
}
