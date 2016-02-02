package main

import (
	"fmt"
	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/nicr9/hathor/backend/hathor"
	"io"
	"net/http"
	"os"
	"time"
)

func download(episode hathor.Episode) {
	err := os.MkdirAll(episode.DirPath, 0755)
	if err != nil {
		fmt.Printf("[e] Couldn't make dir: %s\n -> %s", episode.DirPath, err)
		return
	}

	// Only download if it's not already cached
	if _, err := os.Stat(episode.FilePath); os.IsNotExist(err) {
		fmt.Printf("%s - Downloading '%s'\n", episode.Key, episode.Title)

		// Create the file
		out, err := os.Create(episode.FilePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer out.Close()

		// Get the data
		resp, err := http.Get(episode.Url)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		// Writer the body to file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%s - Finished '%s'\n", episode.Key, episode.Title)
	}
}

func main() {
	config, err := hathor.GetConfig()
	if err != nil {
		fmt.Printf("[e] %s\n", err)
		return
	}
	for key := range config {

		go func(uri string, timeout int) {
			rssfeed := hathor.NewRssFeed(key, config[key])
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
		}(config[key].Source, 5)
	}

	// Wait for episodes to arrive and download 'em
	go hathor.ProcessEpisodes(download)
	hathor.ServeFeeds()
}
