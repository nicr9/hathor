package main

import (
	"fmt"
	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/nicr9/hathor/backend/hathor"
	"os"
	"time"
)

func download(episode hathor.Episode) {
	fmt.Printf("%s - Downloading '%s'\n", episode.Key, episode.Title)
	err := os.MkdirAll(episode.DirPath, 0755)
	if err != nil {
		fmt.Printf("[e] Couldn't make dir: %s\n -> %s", episode.DirPath, err)
	}
}

func main() {
	feeds, err := hathor.GetFeeds()
	if err != nil {
		fmt.Printf("[e] %s\n", err)
		return
	}
	for key := range feeds {

		go func(uri string, timeout int) {
			rssfeed := hathor.NewRssFeed(key, feeds[key])
			feed := rss.New(timeout, true, rssfeed.Channels, rssfeed.Items)
			for {
				if err := feed.Fetch(uri, nil); err != nil {
					fmt.Printf("[e] %s: %s\n", err, uri)
					return
				}
				update := feed.SecondsTillUpdate()
				fmt.Printf("%s - Updating again in %d seconds\n", key, update)
				<-time.After(time.Duration(update) * time.Second)
			}
		}(feeds[key].Source, 5)

		hathor.ProcessEpisodes(download)
	}

}
