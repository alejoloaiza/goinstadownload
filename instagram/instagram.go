package instagram

import (
	"encoding/json"
	"fmt"
	"goinstadownload/config"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/ahmdrz/goinsta"
)

type FollowingUser struct {
	ID         int64
	Username   string
	Fullname   string
	Preference bool
}

var (
	myUsers        []FollowingUser
	myInboxUsers   = make(map[string]int)
	RateLimit      int
	SleepTime      int
	FollowCounter  = 0
	MessageCounter = 0
	FollowingList  = make(map[string]int)
	BlacklistNames = make(map[string]int)
	BlacklistUsers = make(map[string]int)
	FemaleNames    = make(map[string]int)
	TownPreference = make(map[int]string)
	Insta          *goinsta.Instagram
	InChan         chan string
	OutChan        chan string
)

func InstaLogin(in chan string, out chan string) {
	OutChan = out
	InChan = in
	Insta = goinsta.New(config.Localconfig.InstaUser, config.Localconfig.InstaPass)
	err := Insta.Login()
	if err != nil {
		ValidateErrors(err, "Login")
		return
	}
	InChan <- "Connected ok to Instagram"

}
func ListAllFollowing() map[int]string {
	users, err := Insta.UserFollowing(Insta.InstaType.LoggedInUser.ID, "")
	var response = make(map[int]string)
	if err != nil {
		ValidateErrors(err, "UserFollowing")
		return nil
	}
	for i, user := range users.Users {
		response[i] = user.Username
		FollowingList[user.Username] = 1
	}

	return response
}
func Uploadlists() {

	blacklistnameraw := config.Localconfig.BlacklistNames
	for _, bname := range blacklistnameraw {
		BlacklistNames[bname] = 1
	}
	blacklistuserraw := config.Localconfig.BlacklistUsers
	for _, bname2 := range blacklistuserraw {
		BlacklistUsers[bname2] = 1
	}
	femalelistraw := config.Localconfig.FemaleNames
	for _, bname3 := range femalelistraw {
		FemaleNames[bname3] = 1
	}
	townpreferenceraw := config.Localconfig.TownPreference
	for i, bname4 := range townpreferenceraw {
		TownPreference[i] = bname4
	}
}
func InstaDirectMessage(UserId string, Message string) {
	user, err := Insta.GetUserByUsername(UserId)
	id := strconv.FormatInt(user.User.ID, 10)
	if err != nil {
		ValidateErrors(err, "GetUserByUsername")
		return
	}
	resp, err := Insta.DirectMessage(id, Message)

	if err != nil {
		ValidateErrors(err, "DirectMessage")
		return
	}
	fmt.Println(resp)
}
func ValidateErrors(err error, addinfo string) {
	log.Println(addinfo + " " + err.Error())
	InChan <- addinfo + " " + err.Error()
	if strings.Contains(strings.ToLower(err.Error()), "logout") {
		_ = Insta.Login()
		time.Sleep(2 * time.Second)
	}
}
func InstaShowComments(InUserToFollow string) {
	following := make(map[int]string)
	if InUserToFollow != "" {
		following[1] = InUserToFollow
	} else {
		following = ListAllFollowing()
	}

	for _, UserToFollow := range following {
		log.Printf("Checking user: %s ", UserToFollow)
		InChan <- "Checking user: " + UserToFollow
		r, err := Insta.GetUserByUsername(UserToFollow)
		if err != nil {
			ValidateErrors(err, "GetUserByUsername")
			return
		}

		resp, err := Insta.LatestUserFeed(r.User.ID)
		if err != nil {
			ValidateErrors(err, "LatestUserFeed")
			return
		}
		for _, item := range resp.Items {
			time.Sleep(2 * time.Second)
			resp2, _ := Insta.MediaComments(item.ID, "")
			if err != nil {
				ValidateErrors(err, "MediaComments")
			}
			for _, comment := range resp2.Comments {

				fullname := strings.Split(comment.User.FullName, " ")
				firstname := strings.ToLower(fullname[0])
				var gender string
				//log.Printf("Checking comment of: %s ", comment.User.FullName)
				if FemaleNames[firstname] == 1 {
					gender = "female"
				}
				if gender == "female" && BlacklistNames[firstname] != 1 && BlacklistUsers[comment.User.Username] != 1 && UserToFollow != comment.User.Username && FollowingList[comment.User.Username] != 1 {
					time.Sleep(3 * time.Second)
					tofollow, err := Insta.GetUserByID(comment.User.ID)
					jsoninbytes, err := json.Marshal(tofollow)
					jsonuserprofile := string(jsoninbytes)
					for _, preflocation := range TownPreference {
						if strings.Contains(jsonuserprofile, preflocation) {
							_, err = Insta.Follow(comment.User.ID)
							FollowCounter++
							log.Printf(">> #%v Following-> Name:%s \t|User:%s \t|Comment:%s \n", FollowCounter, comment.User.FullName, comment.User.Username, comment.Text)
							InChan <- "Following #" + strconv.Itoa(FollowCounter) + " -> " + comment.User.Username
							FollowingList[comment.User.Username] = 1
							if FollowCounter >= RateLimit {
								//	time.Sleep(12 * time.Hour)
								log.Printf("End of process, #%v Follow requests sent\n", FollowCounter)
								InChan <- "End of process " + strconv.Itoa(FollowCounter)
								return
							}
							if err != nil {
								ValidateErrors(err, "Follow")
							}
						}
					}
				}
				select {
				case msg := <-OutChan:
					if msg == "stop" {
						InChan <- "Stopped process on #" + strconv.Itoa(FollowCounter)
						log.Printf("Stopped, #%v Follow requests sent\n", FollowCounter)
						return
					}
				default:

				}

			}
		}
		time.Sleep(5 * time.Second)
	}
	log.Printf("End of process, #%v Follow requests sent\n", FollowCounter)
	InChan <- "End of process " + strconv.Itoa(FollowCounter)

}

func InstaLogout() {
	Insta.Logout()

}

func PrepareMessage(Message string, NameOfUser string) string {
	resp := ""
	if FemaleNames[strings.ToLower(NameOfUser)] == 1 {
		resp = strings.Replace(Message, "{name}", strings.ToLower(NameOfUser), 1)
	} else {
		resp = strings.Replace(Message, "{name}", "", 1)
	}
	return resp

}
func DirectMessage(To string, Name string, Id int64, Pref bool) {
	max := len(config.Localconfig.Sentences)
	Message := config.Localconfig.Sentences[Random(0, max)]
	newMessage := PrepareMessage(Message, Name)

	_, err := Insta.DirectMessage(strconv.FormatInt(Id, 10), newMessage)
	if err != nil {
		ValidateErrors(err, "DirectMessage")
		return
	}
	MessageCounter++
	if Pref {
		log.Printf("Message #%v with PREFERENCE to %s:%s >> %s \n", MessageCounter, Name, To, newMessage)
	} else {
		log.Printf("Message #%v to %s:%s >> %s \n", MessageCounter, Name, To, newMessage)
	}
	InChan <- "Message sent to " + Name + " User " + To

}
func Random(min int, max int) int {
	return rand.Intn(max-min) + min
}

func InstaRandomMessages() {
	rand.Seed(time.Now().UnixNano())

	var response FollowingUser
	if SleepTime == 0 {
		SleepTime = 10
	}
	inbox, err := Insta.GetV2Inbox("")
	if err != nil {
		ValidateErrors(err, "GetV2Inbox")
		return
	}
	for _, thread := range inbox.Inbox.Threads {
		for _, userthreads := range thread.Users {
			myInboxUsers[userthreads.Username] = 1
		}
	}
	var timeLineCounter int
	var nextMaxID string
	for timeLineCounter < 5 {
		preferences, err := Insta.Timeline(nextMaxID)
		if err != nil {
			ValidateErrors(err, "Timeline")
			return
		}
		nextMaxID = preferences.NextMaxID
		for _, item := range preferences.Items {
			timelocation := strings.ToLower(item.Location.City)
			if timelocation == "" {
				timelocation = strings.ToLower(item.Location.Name)
			}
			if timelocation == "" {
				continue
			}
			for _, preflocation := range TownPreference {

				if strings.Contains(timelocation, preflocation) && myInboxUsers[item.User.Username] != 1 {
					fullname := strings.Split(item.User.FullName, " ")
					firstname := strings.ToLower(fullname[0])
					response.ID = item.User.ID
					response.Username = item.User.Username
					response.Fullname = firstname
					myInboxUsers[item.User.Username] = 1
					response.Preference = true
					myUsers = append(myUsers, response)

				}
			}

		}
		timeLineCounter++
	}

	users, err := Insta.UserFollowing(Insta.InstaType.LoggedInUser.ID, "")
	if err != nil {
		ValidateErrors(err, "UserFollowing")
		return
	}

	for _, user := range users.Users {
		if myInboxUsers[user.Username] != 1 {
			fullname := strings.Split(user.FullName, " ")
			firstname := strings.ToLower(fullname[0])
			response.Username = user.Username
			response.ID = user.ID
			response.Fullname = firstname
			response.Preference = false
			myUsers = append(myUsers, response)
		}
	}

	for _, dmuser := range myUsers {
		//fmt.Println(dmuser.Username, dmuser.Fullname, dmuser.ID)
		DirectMessage(dmuser.Username, dmuser.Fullname, dmuser.ID, dmuser.Preference)

		if MessageCounter >= RateLimit {
			log.Printf("End of process, #%v Messages sent\n", MessageCounter)
			InChan <- "End of process"
			break
		}
		select {
		case msg := <-OutChan:
			if msg == "stop" {
				InChan <- "Stopped process on #" + strconv.Itoa(MessageCounter)
				log.Printf("Stopped, #%v Follow requests sent\n", FollowCounter)
				return
			}
		default:
			time.Sleep(time.Duration(SleepTime) * time.Minute)
		}

	}

}

func InstaTimeLineMessages() {

	rand.Seed(time.Now().UnixNano())

	var response FollowingUser
	if SleepTime == 0 {
		SleepTime = 10
	}
StartProcess:
	for {
		myUsers = make([]FollowingUser, 0)
		inbox, err := Insta.GetV2Inbox("")
		if err != nil {
			ValidateErrors(err, "GetV2Inbox")
			return
		}
		for _, thread := range inbox.Inbox.Threads {
			for _, userthreads := range thread.Users {
				myInboxUsers[userthreads.Username] = 1
			}
		}
		var timeLineCounter int
		var nextMaxID string
		timeLineCounter = 0
		for timeLineCounter < 5 {
			preferences, err := Insta.Timeline(nextMaxID)
			if err != nil {
				ValidateErrors(err, "Timeline")
				time.Sleep(1 * time.Minute)
				continue StartProcess
			}
			nextMaxID = preferences.NextMaxID
			for _, item := range preferences.Items {
				timelocation := strings.ToLower(item.Location.City)
				if timelocation == "" {
					timelocation = strings.ToLower(item.Location.Name)
				}
				if timelocation == "" {
					continue
				}
				for _, preflocation := range TownPreference {

					if strings.Contains(timelocation, preflocation) && myInboxUsers[item.User.Username] != 1 {
						fullname := strings.Split(item.User.FullName, " ")
						firstname := strings.ToLower(fullname[0])
						response.ID = item.User.ID
						response.Username = item.User.Username
						response.Fullname = firstname
						myInboxUsers[item.User.Username] = 1
						response.Preference = true
						myUsers = append(myUsers, response)

					}
				}

			}
			timeLineCounter++
		}

		for _, dmuser := range myUsers {
			//fmt.Println(dmuser.Username, dmuser.Fullname, dmuser.ID)
			DirectMessage(dmuser.Username, dmuser.Fullname, dmuser.ID, dmuser.Preference)

			if MessageCounter >= RateLimit {
				log.Printf("End of process, #%v Messages sent\n", MessageCounter)
				InChan <- "End of process"
				break
			}
			select {
			case msg := <-OutChan:
				if msg == "stop" {
					InChan <- "Stopped process on #" + strconv.Itoa(MessageCounter)
					log.Printf("Stopped, #%v Follow requests sent\n", FollowCounter)
					return
				}
			default:
				time.Sleep(5 * time.Minute)
			}

		}
	}
	time.Sleep(time.Duration(SleepTime) * time.Minute)

}
