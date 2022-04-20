package config

import (
	"github.com/jelliflix/jelliflix/infrastructure/logger"
	"github.com/spf13/viper"
)

type (
	configuration struct {
		Redis        redis
		Services     services
		Integrations integrations
	}

	redis struct {
		Addr string
		Pass string
	}

	services struct {
		IMDB     imdb
		OMDB     omdb
		Torrent  torrent
		Jellyfin jellyfin
	}

	integrations struct {
		Telegram telegram
	}

	imdb struct {
		User string

		RefreshRate int
	}

	omdb struct {
		Token string
	}

	torrent struct {
		Port    int
		Clients []string
		Quality struct {
			Movies string
			Series string
		}
	}

	jellyfin struct {
		Path string
	}

	telegram struct {
		User  int
		Token string
	}
)

var Cfg = getConfig()

func getConfig() (c configuration) {
	viper.SetConfigName("jelliflix")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		logger.Log.WithError(err).Fatal("failed reading config")
	}
	if err := viper.UnmarshalExact(&c); err != nil {
		logger.Log.WithError(err).Fatal("failed parsing config")
	}

	return
}
