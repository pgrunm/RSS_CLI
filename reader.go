package main

import (
	"fmt"
	"log"
	"net/http"

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Log to console that site was accessed
		log.Println("Site was accessed.")

		for _, feed := range feeds {
			rssFeeds := ParseFeeds(feed, proxy)

			for _, rss := range rssFeeds.Items {
				// Needs some more formatting!
				fmt.Fprintf(w, "%s: %s\n", rss.Title, rss.Link)
			}
		}
	})
	http.ListenAndServe(":80", nil)

}
