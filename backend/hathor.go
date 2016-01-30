package main

import (
	"fmt"
	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/nicr9/hathor/backend/hathor"
	"time"
)

var log, errlog, episodes = make(chan string), make(chan error), make(chan Episode)

type Episode struct {
	Podcast string
	Title   string
	Url     string
}

func NewEpisode(rssfeed RssFeed, item *rss.Item) Episode {
	return Episode{rssfeed.Key, item.Title, item.Enclosures[0].Url}
}

type RssFeed struct {
	Key    string
	Config hathor.Feed
}

func (r RssFeed) chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
	log <- fmt.Sprintf(" - %d new channel(s) in %s", len(newchannels), feed.Url)
}

func (r RssFeed) itemHandler(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	log <- fmt.Sprintf(" - %d new item(s) in %s", len(newitems), feed.Url)

	switch r.Config.Download {
	case "latest":
		newitems = []*rss.Item{newitems[0]}
	}

	for _, item := range newitems {
		episodes <- NewEpisode(r, item)
		// TODO: What happens if there's no enclosures?
		// TODO: What happens if there's more than one enclosure?
	}
}

func download(epi Episode) {
	fmt.Printf(" - %+v\n", epi)
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
			rssfeed := RssFeed{key, feeds[key]}
			feed := rss.New(timeout, true, rssfeed.chanHandler, rssfeed.itemHandler)
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
		case res := <-episodes:
			go download(res)
		}
	}
}
