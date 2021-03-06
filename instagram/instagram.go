package instagram

import (
	"encoding/json"
	"fmt"
	"goinstadownload/config"
	"goinstadownload/extra"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ahmdrz/goinsta"
	"golang.org/x/net/proxy"
)

type FollowingUser struct {
	ID         int64
	Username   string
	Fullname   string
	Preference bool
}

var (
	myInboxUsers = make(map[string]int)

	FollowingList  = make(map[string]int)
	BlacklistNames = make(map[string]int)
	BlacklistUsers = make(map[string]int)
	PreferredNames = make(map[string]int)
	TownPreference = make(map[int]string)
	Insta          *goinsta.Instagram
	InChan         chan string
	OutChan        chan string
)

func InstaLogin(in chan string, out chan string) {
	OutChan = out
	InChan = in

	if Insta == nil {
		if config.Localconfig.UseProxy {
			dialer, err := proxy.SOCKS5("tcp", config.Localconfig.Proxy, nil, proxy.Direct)
			if err != nil {
				log.Fatal(err)
			}
			proxyTransport := &http.Transport{}
			proxyTransport.Dial = dialer.Dial
			Insta = goinsta.New(config.Localconfig.InstaUser, config.Localconfig.InstaPass)
			Insta.Transport = *proxyTransport
		} else {
			Insta = goinsta.New(config.Localconfig.InstaUser, config.Localconfig.InstaPass)

		}

	}
	if !Insta.InstaType.IsLoggedIn {
		err := Insta.Login()

		if err != nil {
			ValidateErrors(err, "Login")
			return
		}
		InChan <- "Connected ok to Instagram"
	}
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
	femalelistraw := config.Localconfig.PreferredNames
	for _, bname3 := range femalelistraw {
		PreferredNames[bname3] = 1
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
	log.Println(addinfo + " - " + err.Error())
	InChan <- addinfo + " - " + err.Error()
	if strings.Contains(strings.ToLower(err.Error()), "logout") {
		_ = Insta.Login()
		time.Sleep(2 * time.Second)
	}
}
func InstaShowComments(InUserToFollow string, Limit int) {
	var FollowCounter int = 0
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
			//resp2, _ := Insta.MediaComments(item.ID, "")
			resp2, _ := Insta.MediaLikers(item.ID)
			if err != nil {
				ValidateErrors(err, "MediaLikers")
			}
			for _, comment := range resp2.Users {
				//for _, comment := range resp2.Comments {

				fullname := strings.Split(comment.FullName, " ")
				firstname := strings.ToLower(fullname[0])
				var preferred bool
				//log.Printf("Checking comment of: %s ", comment.User.FullName)
				if PreferredNames[firstname] == 1 {
					preferred = true
				}
				//log.Println(firstname)
				if preferred && BlacklistNames[firstname] != 1 && BlacklistUsers[comment.Username] != 1 && UserToFollow != comment.Username && FollowingList[comment.Username] != 1 {
					time.Sleep(3 * time.Second)
					tofollow, err := Insta.GetUserByID(comment.ID)
					if err != nil {
						ValidateErrors(err, "GetUserByID")
					}
					jsoninbytes, err := json.Marshal(tofollow)
					jsonuserprofile := strings.ToLower(string(jsoninbytes))
				LocationLoop:
					for _, preflocation := range TownPreference {
						if strings.Contains(jsonuserprofile, preflocation) {
							_, err = Insta.Follow(comment.ID)
							FollowCounter++
							log.Printf(">> #%v Following-> Name:%s \t|User:%s \n", FollowCounter, comment.FullName, comment.Username)
							InChan <- "Following #" + strconv.Itoa(FollowCounter) + " -> " + comment.Username
							FollowingList[comment.Username] = 1
							if FollowCounter >= Limit {
								//	time.Sleep(12 * time.Hour)
								log.Printf("End of process, #%v Follow requests sent\n", FollowCounter)
								InChan <- "End of process " + strconv.Itoa(FollowCounter)
								return
							}
							if err != nil {
								ValidateErrors(err, "Follow")
							}
							break LocationLoop
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
	if Insta != nil {
		Insta.Logout()
		Insta = nil
	}
}

func PrepareMessage(Message string, NameOfUser string) string {
	resp := ""
	if PreferredNames[strings.ToLower(NameOfUser)] == 1 {
		resp = strings.Replace(Message, "{name}", strings.ToLower(NameOfUser), 1)
	} else {
		resp = strings.Replace(Message, "{name}", "", 1)
	}
	return resp

}
func DirectMessage(To string, Name string, Id int64, Pref bool) bool {
	max := len(config.Localconfig.Sentences)
	Message := config.Localconfig.Sentences[Random(0, max)]
	newMessage := PrepareMessage(Message, Name)
	/*
		resp, err := Insta.DirectMessage(strconv.FormatInt(Id, 10), "Hola")
		if err != nil {
			ValidateErrors(err, "DirectMessage")
			return
		}
		time.Sleep(1 * time.Second)
		resp2, err := Insta.GetDirectThread(resp.Threads[0].ThreadID)
		if err != nil {
			ValidateErrors(err, "GetDirectThread")
			return
		}
		time.Sleep(5 * time.Second)
		fmt.Println(len(resp2.Thread.Items))
		if len(resp2.Thread.Items) <= 1 {
	*/
	resp3, err := Insta.DirectMessage(strconv.FormatInt(Id, 10), newMessage)
	jsoninbytes, _ := json.Marshal(resp3)
	jsontimeline := string(jsoninbytes)
	jsontimeline = jsontimeline + ""
	//fmt.Println(jsontimeline)
	if err != nil {
		ValidateErrors(err, "DirectMessage")
		return false
	}
	//}

	log.Printf("Message to %s:%s >> %s \n", Name, To, newMessage)

	InChan <- "Message sent to " + Name + " User " + To
	return true
}
func Random(min int, max int) int {
	return rand.Intn(max-min) + min
}

func InstaTimeLineMessages(SleepTime int, Limit int) {
	var myUsers []FollowingUser
	var MessageCounter int = 0
	rand.Seed(time.Now().UnixNano())

	var response FollowingUser
	if SleepTime == 0 {
		SleepTime = 10
	}
StartProcess:
	for {
		log.Println("New cycle started")
		myUsers = make([]FollowingUser, 0)
		/*inbox, err := Insta.GetV2Inbox("")
		if err != nil {
			ValidateErrors(err, "GetV2Inbox")
			return
		}
		for _, thread := range inbox.Inbox.Threads {
			for _, userthreads := range thread.Users {
				myInboxUsers[userthreads.Username] = 1
			}
		}
		*/
		preferences, err := Insta.Timeline("")
		if err != nil {
			ValidateErrors(err, "Timeline")
			time.Sleep(1 * time.Minute)
			continue StartProcess
		}
		for _, item := range preferences.Items {
			/*
				jsoninbytes, _ := json.Marshal(item)
				jsontimeline := strings.ToLower(string(jsoninbytes))
			*/
			if item.Location.Lng != 0 && item.Location.Lat != 0 {
				itemLat := float64(item.Location.Lat)
				itemLng := float64(item.Location.Lng)
				distance := extra.Distance(itemLat, itemLng, config.Localconfig.LocalLat, config.Localconfig.LocalLng)
				log.Printf("Distance in meter is %v", distance)
				if distance < config.Localconfig.MinimumDistance && myInboxUsers[item.User.Username] != 1 {
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

		for _, dmuser := range myUsers {
			//fmt.Println(dmuser.Username, dmuser.Fullname, dmuser.ID)
			if DirectMessage(dmuser.Username, dmuser.Fullname, dmuser.ID, dmuser.Preference) {
				MessageCounter++
			}

			if MessageCounter >= Limit {
				log.Printf("End of process, #%v Messages sent\n", MessageCounter)
				InChan <- "End of process"
				break StartProcess
			}

		}
		select {
		case msg := <-OutChan:
			if msg == "stop" {
				InChan <- "Stopped process on #" + strconv.Itoa(MessageCounter)
				log.Printf("Stopped, #%v Messages requests sent\n", MessageCounter)
				return
			}
		default:
		}
		time.Sleep(time.Duration(SleepTime) * time.Minute)

	}

}

func InstaStoriesMessages(SleepTime int, Limit int) {
	var StoryUser []FollowingUser
	//var myUsers []FollowingUser
	var MessageCounter int = 0
	rand.Seed(time.Now().UnixNano())

	//var response FollowingUser
	if SleepTime == 0 {
		SleepTime = 10
	}
StartProcess:
	for {
		log.Println("New cycle started")
		//myUsers = make([]FollowingUser, 0)

		allstories, err := Insta.GetReelsTrayFeed()
		if err != nil {
			ValidateErrors(err, "GetReelsTrayFeed")
			time.Sleep(1 * time.Minute)
			continue StartProcess
		}
		for _, tray := range allstories.Tray {
			var newuser = FollowingUser{}
			newuser.Username = tray.User.Username
			newuser.Fullname = tray.User.FullName
			StoryUser = append(StoryUser, newuser)
		}
		for _, user := range StoryUser {
			respuser, _ := Insta.GetUserByUsername(user.Username)
			user.ID = respuser.User.ID
			time.Sleep(1 * time.Second)
			_, _ = Insta.GetUserStories(user.ID)

		}
		/*
			jsoninbytes, _ := json.Marshal(item)
			jsontimeline := strings.ToLower(string(jsoninbytes))
		*/
		/*
				if item.Location.Lng != 0 && item.Location.Lat != 0 {
					itemLat := float64(item.Location.Lat)
					itemLng := float64(item.Location.Lng)
					distance := extra.Distance(itemLat, itemLng, config.Localconfig.LocalLat, config.Localconfig.LocalLng)
					log.Printf("Distance in meter is %v", distance)
					if distance < config.Localconfig.MinimumDistance && myInboxUsers[item.User.Username] != 1 {
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

			for _, dmuser := range myUsers {
				//fmt.Println(dmuser.Username, dmuser.Fullname, dmuser.ID)
				if DirectMessage(dmuser.Username, dmuser.Fullname, dmuser.ID, dmuser.Preference) {
					MessageCounter++
				}

				if MessageCounter >= Limit {
					log.Printf("End of process, #%v Messages sent\n", MessageCounter)
					InChan <- "End of process"
					break StartProcess
				}

			}

		*/
		select {
		case msg := <-OutChan:
			if msg == "stop" {
				InChan <- "Stopped process on #" + strconv.Itoa(MessageCounter)
				log.Printf("Stopped, #%v Messages requests sent\n", MessageCounter)
				return
			}
		default:
		}
		time.Sleep(time.Duration(SleepTime) * time.Minute)

	}

}
