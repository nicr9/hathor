package main

import (
	"fmt"
	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/nicr9/hathor/backend/hathor"
	"time"
)

func download(episode hathor.Episode) {
	fmt.Printf(" - %+v\n", episode)
}

func main() {
	feeds, err := hathor.GetFeeds()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("--- Feeds ---")
	for key := range feeds {
		fmt.Printf(" * %s\n", key)

		go func(uri string, timeout int) {
			rssfeed := hathor.NewRssFeed(key, feeds[key])
			feed := rss.New(timeout, true, rssfeed.Channels, rssfeed.Items)
			for {
				if err := feed.Fetch(uri, nil); err != nil {
					fmt.Printf("[e] %s: %s\n", err, uri)
					return
				}
				<-time.After(time.Duration(feed.SecondsTillUpdate() * 1e9))
			}
		}(feeds[key].Source, 5)

		hathor.ProcessEpisodes(download)
	}

}
