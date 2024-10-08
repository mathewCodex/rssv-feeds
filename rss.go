package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func urlToFeed(url string) (RSSFeed, error) {
	httpClient := http.Client{//creating a http client to make a get req
		Timeout: 10 * time.Second,
	}
	resp, err := httpClient.Get(url)
	if (err != nil){
		return RSSFeed{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		
		return RSSFeed{},  err
	}
	rssFeed := RSSFeed{}
	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		
		return RSSFeed{},  err
	}
	return rssFeed, nil
}