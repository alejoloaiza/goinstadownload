package main

import (
	"goinstadownload/config"
	"goinstadownload/instagram"
	"goinstadownload/irc"
	"os"
)

func main() {
	var configpath string
	configpath = "../config/config.json"
	if len(os.Args) >= 2 {
		configpath = os.Args[1]
	}
	_ = config.GetConfig(configpath)
	instagram.Uploadlists()

	irc.StartIRCprocess()

}
