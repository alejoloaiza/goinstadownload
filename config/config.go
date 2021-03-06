package config

import "encoding/json"
import "os"
import "fmt"

type Configuration struct {
	InstaUser       string
	InstaPass       string
	BlacklistNames  []string
	BlacklistUsers  []string
	PreferredNames  []string
	Sentences       []string
	TownPreference  []string
	IRCNick         string
	IRCChannels     string
	IRCUser         string
	IRCServerPort   string
	LocalLat        float64
	LocalLng        float64
	MinimumDistance float64
	Proxy           string
	UseProxy        bool
}

var Localconfig *Configuration

func GetConfig(configpath string) *Configuration {
	file, _ := os.Open(configpath)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	Localconfig = &configuration
	return &configuration
}
