package hathor

import (
	"gopkg.in/yaml.v2"
)

// TEST CONFIG
var config string = `
podcast.__init__:
  source: "http://podcastinit.podbean.com/feed/"
`

type Feed struct {
	Source string "source"
}

func GetFeeds() (feeds map[string]Feed, err error) {
	feeds = make(map[string]Feed)
	err = yaml.Unmarshal([]byte(config), &feeds)

	return
}
