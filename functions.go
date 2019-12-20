package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/mmcdole/gofeed"
	"github.com/patrickmn/go-cache"
)

// ParseFeeds allows to get feeds from a site.
func ParseFeeds(siteURL, proxyURL string) *gofeed.Feed {
	// Proxy URL see https://stackoverflow.com/questions/14661511/setting-up-proxy-for-http-client
	var client http.Client

	// Proxy URL is given
	if len(proxyURL) > 0 {
		proxyURL, err := url.Parse(proxyURL)
		if err != nil {
			fmt.Println(err)
		}

		client = http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	} else {
		client = http.Client{}
	}


	item, found := c.Get(siteURL)
	if found {
		log.Printf("Cache hit for site: %s", siteURL)

		//  Type assertion see: https://golangcode.com/convert-interface-to-number/
		return item.(*gofeed.Feed)
	}
	log.Printf("No cache hit for site: %s", siteURL)
	// Get the Feed of the particular website
	resp, err := client.Get(siteURL)

	if err != nil {
		fmt.Println(err)
	} else {
		defer resp.Body.Close()
		// Read the response and parse it as string
		body, _ := ioutil.ReadAll(resp.Body)
		fp := gofeed.NewParser()
		feed, _ := fp.ParseString(string(body))

		c.Set(siteURL, feed, cache.DefaultExpiration)

		// Return the feed with all its items.
		return feed

	}

	return nil
}
