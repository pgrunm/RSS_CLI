package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/mmcdole/gofeed"
)

// ParseFeeds allows to get feeds from a site.
func ParseFeeds(siteURL, proxyURL string) {
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

		fmt.Println(feed.Link)
		for _, e := range feed.Items {
			fmt.Printf("%s: %s\n", e.Title, e.Link)
		}
	}

}
