package main

import (
	"fmt"
	"goinstadownload/config"
	"goinstadownload/instagram"
	"os"
	"strconv"
	"time"
)

func main() {

	var arg1 string
	var arg2 int
	if len(os.Args) > 2 {
		arg1 = os.Args[1]
		arg2, _ = strconv.Atoi(os.Args[2])
	} else {
		fmt.Println("Error: Usage command <configpath> <wait time in minutes>")
	}

	_ = config.GetConfig("../config/config.json")
	instagram.RateLimit = arg2
	instagram.Uploadlists()
	instagram.InstaLogin()

	if arg1 == "--follow" {
		r := instagram.ListAllFollowing()
		for _, user := range r {
			fmt.Printf("Checking user: %s \n", user)
			instagram.InstaShowComments(user)
			time.Sleep(time.Second * 15)
		}
	}
	if arg1 == "--message" {
		instagram.InstaRandomMessages()
	}
	defer instagram.InstaLogout()

}
