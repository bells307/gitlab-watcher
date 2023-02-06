package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Listen             string          `mapstructure:"listen"`
	Telegram           *TelegramConfig `mapstructure:"telegram"`
	SendTemplateErrors bool            `mapstructure:"send_template_errors"`
}

type TelegramConfig struct {
	Token     string            `mapstructure:"token"`
	Debug     bool              `mapstructure:"debug"`
	ParseMode string            `mapstructure:"parse_mode"`
	Templates string            `mapstructure:"templates"`
	Chats     []int64           `mapstructure:"chats"`
	Users     map[string]string `mapstructure:"users"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetDefault("listen", "0.0.0.0:8888")

	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
