package main

import (
	"os"

	server "github.com/WildEgor/notifier/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	// HINT: set logger output
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	// HINT: load global .env
	viper.SetConfigName(".env")
	viper.SetConfigType("dotenv")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Info("The .env file has not been found in the current directory")
		} else {
			log.Error(err.Error())
		}
	}
}

func main() {
	server, _ := server.New()
	port := viper.Get("APP_PORT")
	log.Fatal(server.Listen(port))
}
