package irc

import (
	"bufio"
	"fmt"
	"goinstadownload/config"
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
	InfoChan   chan string
)

func StartIRCprocess() {
	//MsgChan := make(chan string)
	//allconfig := config.GetConfig(configpath)
	InfoChan = make(chan string)
	for {
		Connection, err := net.Dial("tcp", config.Localconfig.IRCServerPort)

		if err != nil {
			fmt.Println(err)
			time.Sleep(2000 * time.Millisecond)
			continue
		}

		fmt.Fprintln(Connection, "NICK "+config.Localconfig.IRCNick)
		fmt.Fprintln(Connection, "USER "+config.Localconfig.IRCUser)
		fmt.Fprintln(Connection, "JOIN "+config.Localconfig.IRCChannels)
		go RoutineWriter(Connection)
		MyReader := bufio.NewReader(Connection)
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
				break
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
	if strings.TrimSpace(command[0]) == "init" && len(command) >= 3 {
		var arg3 int
		arg1 := command[1]
		arg2, err := strconv.Atoi(strings.TrimSuffix(strings.TrimSuffix(command[2], "\n"), "\r"))
		if len(command) >= 4 {
			arg3, err = strconv.Atoi(strings.TrimSuffix(strings.TrimSuffix(command[3], "\n"), "\r"))
		}
		if err != nil {
			log.Println(err)
			return ""
		}
		instagram.RateLimit = arg2
		instagram.SleepTime = arg3
		if arg1 == "-follow" {
			go ExecuteFollowProcess()
		}
		if arg1 == "-message" {
			go ExecuteMessageProcess()
		}

	}
	bodyString = "Executed in background"
	return bodyString
}
func ExecuteFollowProcess() {
	instagram.FollowCounter = 0
	instagram.InstaLogin(InfoChan)
	instagram.InstaShowComments()
	defer instagram.InstaLogout()
}

func ExecuteMessageProcess() {
	instagram.MessageCounter = 0
	instagram.InstaLogin(InfoChan)
	instagram.InstaRandomMessages()
	defer instagram.InstaLogout()
}
func RoutineWriter(Response net.Conn) {
	for {
		select {
		case msg := <-InfoChan:
			if Context != "" {
				fmt.Fprintln(Response, "PRIVMSG "+Context+" :"+msg)
			}
		}
		time.Sleep(time.Second * 2)
	}
}
