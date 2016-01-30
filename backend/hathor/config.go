package hathor

import (
	"gopkg.in/yaml.v2"
)

// TEST CONFIG
var config string = `
podcastinit:
  source: "http://podcastinit.podbean.com/feed/"
  download: latest
`

type Feed struct {
	Source   string "source"
	Download string "download,omitempty"
}

func GetFeeds() (feeds map[string]Feed, err error) {
	feeds = make(map[string]Feed)
	err = yaml.Unmarshal([]byte(config), &feeds)

	return
}
