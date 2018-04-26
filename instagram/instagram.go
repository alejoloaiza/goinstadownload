package instagram

import (
	"fmt"
	"goinstadownload/config"
	"log"

	"github.com/ahmdrz/goinsta"
)

var Insta *goinsta.Instagram

func InstaLogin() {
	Insta = goinsta.New(config.Localconfig.InstaUser, config.Localconfig.InstaPass)
	fmt.Println(config.Localconfig.InstaUser, config.Localconfig.InstaPass)
	if err := Insta.Login(); err != nil {
		panic(err)
	}

}

func InstaShowComments(userIDToSpy string) {
	r, err := Insta.GetUserByUsername(userIDToSpy)
	if err != nil {
		fmt.Println(err)
		log.Println(err)
		return
	}
	resp, err := Insta.LatestUserFeed(r.User.ID)

	for _, item := range resp.Items {
		resp2, _ := Insta.MediaComments(item.ID, "")

		for _, comment := range resp2.Comments {

			fmt.Printf("Name:%s |User:%s |Comment:%s \n", comment.User.FullName, comment.User.Username, comment.Text)
		}
	}
}

func InstaLogout() {
	Insta.Logout()
}
