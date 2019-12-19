package main

import (
	"log"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Parse configuration
	feeds := viper.GetStringSlice("Feeds")
	proxy := viper.GetString("Proxy")

	for _, feed := range feeds {
		ParseFeeds(feed, proxy)
	}

}
