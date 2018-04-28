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
	var sleeptime int
	if len(os.Args) > 2 {
		arg1 = os.Args[1]
		sleeptime, _ = strconv.Atoi(os.Args[2])
	} else {
		fmt.Println("Error: Usage command <configpath> <wait time in minutes>")
	}

	_ = config.GetConfig(arg1)
	instagram.Uploadlists()
	for {
		instagram.InstaLogin()
		r := instagram.ListAllFollowing()
		for _, user := range r {
			fmt.Printf("Checking user: %s \n", user)
			instagram.InstaShowComments(user)
			time.Sleep(time.Second * time.Duration(sleeptime))
		}
		instagram.InstaLogout()
		fmt.Printf("============== WAITING FOR NEXT CYCLE ===============")
		time.Sleep(time.Minute * 360)
	}
}
