package main

import (
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
func ParseFeeds(siteURL, proxyURL string) (*gofeed.Feed, error) {

	// Measure the execution time of this function
	defer duration(track("ParseFeeds for site " + siteURL))

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
		//  Type assertion see: https://golangcode.com/convert-interface-to-number/
		return item.(*gofeed.Feed), nil
	}
	// Get the Feed of the particular website
	resp, err := client.Get(siteURL)

	if err != nil {
		fmt.Println(err)
	} else {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			// Read the response and parse it as string
			body, _ := ioutil.ReadAll(resp.Body)
			fp := gofeed.NewParser()
			feed, _ := fp.ParseString(string(body))

			c.Set(siteURL, feed, cache.DefaultExpiration)

			// Return the feed with all its items.
			return feed, nil
		}
		return nil, fmt.Errorf("Return code for site %s was different than 200: %d", siteURL, resp.StatusCode)
	}
	return nil, nil

}

// Source: https://yourbasic.org/golang/measure-execution-time/
func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
	log.Printf("%v: %v\n", msg, time.Since(start))
}
