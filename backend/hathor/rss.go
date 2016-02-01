package hathor

import (
	"fmt"
	rss "github.com/jteeuwen/go-pkg-rss"
	"net/http"
	"net/url"
	"path"
)

var CachePath string = "/tmp/hathor"

type Episode struct {
	Key      string
	Title    string
	Url      string
	DirPath  string
	FilePath string
}

func NewEpisode(rssfeed RssFeed, item *rss.Item) Episode {
	key := rssfeed.Key
	title := item.Title
	uri := item.Enclosures[0].Url

	// Parse uri and determine path to save file locally
	dirPath, filePath := "", ""
	u, err := url.Parse(uri)
	if err != nil {
		fmt.Printf("[e] Failed to parse uri (podcast:%s, episode:%s)\n", key, title)
	} else {
		resourcePath := u.Path
		fileName := path.Base(resourcePath)
		dirPath = path.Join(CachePath, key)
		filePath = path.Join(dirPath, fileName)
	}

	result := Episode{rssfeed.Key, item.Title, uri, dirPath, filePath}
	return result
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
	fmt.Printf("%s - %d new channel(s)\n", r.Key, len(newchannels))
}

func (r RssFeed) Items(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	fmt.Printf("%s - %d new item(s)\n", r.Key, len(newitems))

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

func ServeFeeds() {
	fmt.Println("hathor - starting server...")
	http.ListenAndServe(":8080", http.FileServer(http.Dir(CachePath)))
}
