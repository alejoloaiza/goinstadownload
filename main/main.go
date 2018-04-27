package main

import (
	"fmt"
	"goinstadownload/config"
	"goinstadownload/instagram"
	"os"
	"time"
)

func main() {
	arg := "../config/config.json"
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}
	_ = config.GetConfig(arg)
	instagram.InstaLogin()
	r := instagram.ListAllFollowing()
	for _, user := range r {
		fmt.Printf("Checking user: %s ", user)
		time.Sleep(5 * time.Minute)
		instagram.InstaShowComments(user)
	}
	defer instagram.InstaLogout()
}
