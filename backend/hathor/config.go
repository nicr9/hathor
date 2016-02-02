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

func GetConfig() (config map[string]Feed, err error) {
	config = make(map[string]Feed)
	err = yaml.Unmarshal([]byte(config), &config)

	return
}
