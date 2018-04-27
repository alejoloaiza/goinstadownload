package instagram

import (
	"fmt"
	"goinstadownload/api"
	"goinstadownload/config"
	"log"
	"strings"

	"github.com/ahmdrz/goinsta"
)

var Insta *goinsta.Instagram

func InstaLogin() {
	Insta = goinsta.New(config.Localconfig.InstaUser, config.Localconfig.InstaPass)
	if err := Insta.Login(); err != nil {
		panic(err)
	}

}
func ListAllFollowing() map[int]string {
	users, err := Insta.UserFollowing(Insta.InstaType.LoggedInUser.ID, "")
	var response = make(map[int]string)
	if err != nil {
		return nil
	}
	for i, user := range users.Users {
		response[i] = user.Username
	}
	return response
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
			fullname := strings.Split(comment.User.FullName, " ")
			gender := api.GetGender(fullname[0])
			if gender == "female" {
				fmt.Printf(">> Following-> Name:%s \t|User:%s \t|Comment:%s \n", comment.User.FullName, comment.User.Username, comment.Text)
				_, err = Insta.Follow(comment.User.ID)
			}

		}
	}
}

func InstaLogout() {
	Insta.Logout()
}
