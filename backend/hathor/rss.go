package hathor

import (
	"fmt"
	rss "github.com/jteeuwen/go-pkg-rss"
)

type Episode struct {
	Podcast string
	Title   string
	Url     string
}

func NewEpisode(rssfeed RssFeed, item *rss.Item) Episode {
	return Episode{rssfeed.Key, item.Title, item.Enclosures[0].Url}
}

var episodes = make(chan Episode)

type RssFeed struct {
	Key    string
	Config Feed
}

func NewRssFeed(key string, config Feed) RssFeed {
	return RssFeed{key, config}
}

func (r RssFeed) Channels(feed *rss.Feed, newchannels []*rss.Channel) {
	fmt.Printf(" - %d new channel(s) in %s\n", len(newchannels), feed.Url)
}

func (r RssFeed) Items(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	fmt.Printf(" - %d new item(s) in %s\n", len(newitems), feed.Url)

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

func ProcessEpisodes(process func(Episode)) {
	for {
		res := <-episodes
		go process(res)
	}
}