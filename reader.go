package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
)

var (
	// Adding a cache with 5 min expiration and deletion after 10 mins
	c = cache.New(5*time.Minute, 10*time.Minute)
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
		log.Printf("Site was accessed from %s.", r.RemoteAddr)

		for _, feed := range feeds {
			rssFeeds := ParseFeeds(feed, proxy)

			// Print the title of the news site
			fmt.Fprintf(w, "<p>%s </p>", rssFeeds.Title)
			for _, rss := range rssFeeds.Items {
				// Needs some more formatting!
				fmt.Fprintf(w, "<a href=%s>%s</a> <br>", rss.Link, rss.Title)
			}
		}
	})
	http.ListenAndServe(":80", nil)

}
