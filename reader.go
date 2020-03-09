package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/mmcdole/gofeed"
	"github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

var (
	// Adding a cache with 5 min expiration and deletion after 10 mins
	c = cache.New(5*time.Minute, 10*time.Minute)

	// Create a wait group
	wg sync.WaitGroup

	// Rate Limiting
	rate     = time.Second / 10
	throttle = time.Tick(rate)
)

func main() {
	var feeds []string
	var proxy string

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.WatchConfig()

	// If config file is changed update all configuration values
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
		feeds = viper.GetStringSlice("Feeds")
		proxy = viper.GetString("Proxy")
	})

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Parse configuration
	feeds = viper.GetStringSlice("Feeds")
	proxy = viper.GetString("Proxy")

	// Adding the Prmetheus HTTP handler
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2112", nil)

	// Notify about the started website
	log.Println("Service started: open http://127.0.0.1 in browser")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer duration(track("Time for processing all sites "))

		// Creating a map for the returned feeds
		m := make(map[string]chan *gofeed.Feed, len(feeds))

		// Log to console that site was accessed
		log.Printf("Site was accessed from %s.", r.RemoteAddr)

		for _, feed := range feeds {
			wg.Add(1)
			chNews := make(chan *gofeed.Feed, 1)
			m[feed] = chNews
			go ParseFeeds(feed, proxy, m[feed])
		}

		// Stop execution until the wait group is finished
		wg.Wait()

		defer duration(track("Time for rendering"))

		// Get the items for a feed by order they were mentioned in the configuration file.
		for _, feed := range feeds {

			// Close the channel and write the information from the channel to a variable.
			close(m[feed])
			rss := <-m[feed]
			// Print the title of the news site
			// Needs some more formatting!
			fmt.Fprintf(w, "<p>%s </p>", rss.Title)
			for _, rssFeeds := range rss.Items {
				fmt.Fprintf(w, "<a href=%s>%s</a> <br>", rssFeeds.Link, rssFeeds.Title)
			}
		}
	})

	http.ListenAndServe(":80", nil)

}
