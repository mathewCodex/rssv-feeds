package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration,){
	log.Printf("Scraping on %v goroutines every %s  duration", concurrency, timeBetweenRequest,)
	ticker :=  time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Printf("error fetching feeds", err)
			continue
		}
		wg := &sync.WaitGroup{}//using a sync mechanism to fetch each group at a time
		for _, feed := range feeds {//iterating thru feeds
			wg.Add(1)//adding one to the wait group for every feed
				//adding a new goroutine
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries,wg *sync.WaitGroup, feed database.Feed){
	defer wg.Done()//decrement the counter by one 

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetch:", err)
		return
	}
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching Feed:", err)
		return
	}
	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{
			
		}
		if item.Description != ""{
			description.String = item.Description
			description.Valid = true
		}
		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Couldn't parse date %v with err %v",  item.PubDate, err)
			continue
		}
		_, err = db.CreatePost(context.Background(),  database.CreatePostParams{
			ID: 	uuid.New(),
			CreatedAt: 		time.Now().UTC(),
			UpdatedAt: 		time.Now().UTC(),
			Title: 			item.title,
			Description: description,
			PublishedAt: pubAt,
			Url: item.Link,
			FeedID: feed.ID,

		})
		if err != nil{
			if strings.Contains(err.Error(), "duplicate Key"){
				continue
			}
			log.Println("Failed to create Post :", err)
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}