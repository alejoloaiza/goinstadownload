package instagram

import (
	"fmt"
	"goinstadownload/config"
	"log"
	"strings"
	"time"

	"github.com/ahmdrz/goinsta"
)

var Insta *goinsta.Instagram
var BlacklistNames = make(map[string]int)
var BlacklistUsers = make(map[string]int)
var FemalelistNames = make(map[string]int)

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
func Uploadlists() {
	blacklistraw := config.Localconfig.BlacklistNames
	for _, bname := range blacklistraw {
		BlacklistNames[bname] = 1
	}
	blacklistraw = config.Localconfig.BlacklistUsers
	for _, bname2 := range blacklistraw {
		BlacklistUsers[bname2] = 1
	}
	femalelistraw := config.Localconfig.FemaleNames
	for _, bname3 := range femalelistraw {
		FemalelistNames[bname3] = 1
	}
}
func InstaDirectMessage() {
	resp, err := Insta.DirectMessage("some_girl", "hello")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}
func InstaShowComments(userIDToSpy string) {

	r, err := Insta.GetUserByUsername(userIDToSpy)
	if err != nil {
		fmt.Println(err)
		log.Println(err)
		return
	}

	resp, err := Insta.LatestUserFeed(r.User.ID)
	if err != nil {
		fmt.Println(err)
		log.Println(err)
		return
	}
	for _, item := range resp.Items {
		time.Sleep(2 * time.Second)
		resp2, _ := Insta.MediaComments(item.ID, "")
		if err != nil {

			fmt.Println(err)
			log.Println(err)
		}
		for _, comment := range resp2.Comments {
			fullname := strings.Split(comment.User.FullName, " ")
			firstname := strings.ToLower(fullname[0])
			var gender string
			if FemalelistNames[firstname] == 1 {
				gender = "female"
			}
			/*if len(fullname) > 1 {
				gender = api.GetGender(fullname[0] + "/" + fullname[1])
			}*/
			//fmt.Printf(">> COMMENT-> Name:%s \t|User:%s \t|Comment:%s \n", comment.User.FullName, comment.User.Username, comment.Text)
			//log.Printf("%s %d %d %s \n", gender, BlacklistNames[firstname], BlacklistUsers[comment.User.Username], comment.User.Username)
			if gender == "female" && BlacklistNames[firstname] != 1 && BlacklistUsers[comment.User.Username] != 1 && userIDToSpy != comment.User.Username {
				time.Sleep(3 * time.Second)
				log.Printf(">> Following-> Name:%s \t|User:%s \t|Comment:%s \n", comment.User.FullName, comment.User.Username, comment.Text)
				_, err = Insta.Follow(comment.User.ID)
				if err != nil {
					fmt.Println(err)
					log.Println(err)
				}
			}

		}
	}
}

func InstaLogout() {
	Insta.Logout()

}
