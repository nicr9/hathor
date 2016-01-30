package main

import (
	"fmt"
	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/nicr9/hathor/backend/hathor"
	"time"
)

var log, errlog = make(chan string), make(chan error)

func chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
	log <- fmt.Sprintf(" - %d new channel(s) in %s", len(newchannels), feed.Url)
}

func itemHandler(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	log <- fmt.Sprintf(" - %d new item(s) in %s", len(newitems), feed.Url)
}

func main() {
	feeds, err := hathor.GetFeeds()
	if err != nil {
		errlog <- err
	}
	fmt.Println("--- Feeds ---")
	for key := range feeds {
		fmt.Printf(" * %s\n", key)

		go func(uri string, timeout int) {
			feed := rss.New(timeout, true, chanHandler, itemHandler)
			for {
				if err := feed.Fetch(uri, nil); err != nil {
					errlog <- fmt.Errorf("%s: %s", err, uri)
					return
				}
				<-time.After(time.Duration(feed.SecondsTillUpdate() * 1e9))
			}
		}(feeds[key].Source, 5)
	}

	for {
		select {
		case res := <-log:
			fmt.Println(res)
		case res := <-errlog:
			fmt.Printf("[e] %s\n", res)
		}
	}
}
