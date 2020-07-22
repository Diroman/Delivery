package config

import (
	"github.com/spf13/viper"
	"log"
	"path/filepath"
)

type Config struct {
	Server struct {
		GrpcPort int
	}
	Database struct {
		Port     int
		Host     string
		User     string
		Password string
		DBName   string
	}
}

func Read() Config {
	var conf Config
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(filepath.Join("config"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}

	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatalln(err)
	}
	return conf
}
