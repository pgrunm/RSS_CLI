package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/patrickmn/go-cache"
)

// ParseFeeds allows to get feeds from a site.
func ParseFeeds(siteURL, proxyURL string, news chan<- *gofeed.Feed, ProxyUser string, ProxyPass string) {

	// Measure the execution time of this function
	defer duration(track("ParseFeeds for site " + siteURL))

	// When finished, write it to the channel
	defer wg.Done()

	// Proxy URL see https://stackoverflow.com/questions/14661511/setting-up-proxy-for-http-client
	var client http.Client

	// Proxy URL is given
	if len(proxyURL) > 0 {
		proxyURL, err := url.Parse(proxyURL)
		if err != nil {
			fmt.Println(err)
		}

		if len(ProxyPass) > 0 {
			// Add Header with Proxy Username and Pass if inside config
			auth := fmt.Sprintf("%s:%s", ProxyUser, ProxyPass)
			basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
			header := http.Header{}
			header.Add("Proxy-Authorization", basicAuth)
			client = http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL), ProxyConnectHeader: header}}
		} else {
			client = http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
		}

	} else {
		client = http.Client{}
	}

	item, found := c.Get(siteURL)
	if found {
		//  Type assertion see: https://golangcode.com/convert-interface-to-number/
		news <- item.(*gofeed.Feed)

		// Increase the counter for cache hits
		cacheHits.Inc()
	} else {
		// rate limit the feed parsing
		<-throttle

		rssRequests.Inc()

		// Changed this to NewRequest as the golang docs says you need this for custom headers
		req, err := http.NewRequest("GET", siteURL, nil)
		if err != nil {
			log.Fatalln(err)
		}

		// Set a custom user header because some site block away default crawlers
		req.Header.Set("User-Agent", "Golang/RSS_Reader by Warryz")
		// Get the Feed of the particular website
		resp, err := client.Do(req)

		if err != nil {
			fmt.Println(err)
			exceptions.Inc()
		} else {
			defer resp.Body.Close()
			if resp.StatusCode == 200 {
				// Read the response and parse it as string
				body, _ := ioutil.ReadAll(resp.Body)
				fp := gofeed.NewParser()
				feed, _ := fp.ParseString(string(body))

				// Return the feed with all its items.
				if feed != nil {
					c.Set(siteURL, feed, cache.DefaultExpiration)
					news <- feed
				}
			}
		}
	}
}

// Source: https://yourbasic.org/golang/measure-execution-time/
func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
	log.Printf("%v: %v\n", msg, time.Since(start))
}
