package irc

import (
	"bufio"
	"fmt"
	"goinstadownload/config"
	"goinstadownload/extra"
	"goinstadownload/instagram"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	Context    string
	Connection net.Conn
	InChan     chan string
	OutChan    chan string
)

func StartIRCprocess() {
	//MsgChan := make(chan string)
	//allconfig := config.GetConfig(configpath)
	InChan = make(chan string)
	OutChan = make(chan string)
MainCycle:
	for {
		Connection, err := net.Dial("tcp", config.Localconfig.IRCServerPort)

		if err != nil {
			fmt.Println(err)
			time.Sleep(2000 * time.Millisecond)
			continue MainCycle
		}

		fmt.Fprintln(Connection, "NICK "+config.Localconfig.IRCNick)
		fmt.Fprintln(Connection, "USER "+config.Localconfig.IRCUser)
		fmt.Fprintln(Connection, "JOIN "+config.Localconfig.IRCChannels)
		go RoutineWriter(Connection)
		MyReader := bufio.NewReader(Connection)
	ReaderCycle:
		for {

			message, err := MyReader.ReadString('\n')
			// atomixxx: To handle if connection is closed, and jump to next execution.
			if err != nil {
				fmt.Println(time.Now().Format(time.Stamp) + ">>>" + err.Error())
				if io.EOF == err {
					Connection.Close()
					fmt.Println("server closed connection")
				}
				time.Sleep(2000 * time.Millisecond)
				break ReaderCycle
			}

			fmt.Print(time.Now().Format(time.Stamp) + ">>" + message)

			// atomixxx: Split the message into words to better compare between different commands
			text := strings.Split(message, " ")
			//fmt.Println("Number of objects in text: "+ strconv.Itoa(len(text)))
			var respond bool = false
			var response string
			// atomixxx: Logic to detect messages, BOT logic should go inside this
			if len(text) >= 4 && text[1] == "PRIVMSG" {
				respond = true
				var repeat bool = true
				var respondTo string
				//atomixxx logic to differ if message is channel or private from user
				if text[2][0:1] == "#" {
					// logic to respond the same thing to a channel / repeater BOT
					respondTo = text[2]
					Context = respondTo
				} else {
					userto := strings.Split(text[0], "!")
					respondTo = userto[0][1:]
					Context = respondTo
					// logic to respond the same thing to a user / repeater BOT
				}
				// If its a command BOT will execute the command given
				if text[3] == ":!cmd" {
					repeat = false
					commandresponse := ProcessCommand(text[4:])
					response = "PRIVMSG " + respondTo + " :" + commandresponse

				}
				// If is not a command BOT will repeat the same thing
				if repeat == true {
					response = "PRIVMSG " + respondTo + " " + strings.Join(text[3:], " ")

				}
			}
			// atomixxx: Ping/Pong handler to avoid timeout disconnect from the irc server
			if len(text) == 2 && text[0] == "PING" {
				response = "PONG " + text[1]
				respond = true
			}
			// This checks if the received text requires response or not, and respond according to the above logic

			if respond == true {
				fmt.Fprintln(Connection, response)
				fmt.Println(time.Now().Format(time.Stamp) + "<<" + response)
			}

		}
		// atomixxx: If connection is closed, will try to reconnect after 2 seconds
		time.Sleep(2000 * time.Millisecond)
	}

}

func ProcessCommand(command []string) string {
	var bodyString string
	var UserToFollow string = ""
	if strings.TrimSpace(command[0]) == "stop" {
		OutChan <- "stop"
		bodyString = "Command received... processing"
	}
	if len(command) >= 3 && strings.TrimSpace(command[0]) == "set" {
		if extra.RemoveEnds(command[1]) == "proxy" {
			config.Localconfig.UseProxy, _ = strconv.ParseBool(extra.RemoveEnds(command[2]))
			bodyString = "Proxy set to " + strconv.FormatBool(config.Localconfig.UseProxy)
		}
		if extra.RemoveEnds(command[1]) == "logout" {
			instagram.InstaLogout()
			bodyString = "Logged out"

		}
	}
	if len(command) >= 3 && strings.TrimSpace(command[0]) == "init" {
		var arg2, arg3 int
		var err error
		arg1 := command[1]
		if extra.IsInteger(extra.RemoveEnds(command[2])) {
			arg2, err = strconv.Atoi(extra.RemoveEnds(command[2]))
		} else {
			UserToFollow = extra.RemoveEnds(command[2])
		}

		if len(command) >= 4 {
			arg3, err = strconv.Atoi(extra.RemoveEnds(command[3]))
		}
		if err != nil {
			return ""
		}
		if arg1 == "follow" {
			go ExecuteFollowProcess(UserToFollow, arg2)
		}
		if arg1 == "message" {
			go ExecuteMessageProcess(arg3)
		}
		if arg1 == "auto" {
			go ExecuteAutomaticMode(arg3, arg2)

		}
		bodyString = "Command received... processing"
	}

	return bodyString
}
func ExecuteFollowProcess(UserToFollow string, Limit int) {
	instagram.InstaLogin(InChan, OutChan)
	instagram.InstaShowComments(UserToFollow, Limit)
	//	defer instagram.InstaLogout()
}

func ExecuteMessageProcess(SleepTime int) {
	instagram.InstaLogin(InChan, OutChan)
	//instagram.InstaRandomMessages(SleepTime)
	//	defer instagram.InstaLogout()
}

func ExecuteAutomaticMode(SleepTime int, Limit int) {
	instagram.InstaLogin(InChan, OutChan)
	instagram.InstaTimeLineMessages(SleepTime, Limit)
	//	defer instagram.InstaLogout()
}
func RoutineWriter(Response net.Conn) {
	for {
		var err error
		select {
		case msg := <-InChan:
			if Context != "" {
				_, err = fmt.Fprintln(Response, "PRIVMSG "+Context+" :"+msg)
			} else {
				_, err = fmt.Fprintln(Response, "PRIVMSG "+config.Localconfig.IRCChannels+" :"+msg)
			}
		}
		if err != nil {
			log.Println("Error during Write in RoutineWriter")
			break
		}
	}
}
