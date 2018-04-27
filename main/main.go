package main

import (
	"fmt"
	"goinstadownload/config"
	"goinstadownload/instagram"
	"os"
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
		fmt.Println(user)
	}
	//instagram.InstaShowComments(config.Localconfig.UserToSpy)
	defer instagram.InstaLogout()
}
