# INSTAGRAM FOLLOWING BOT

## WTH is this???

Okay this is a program that connects to instagram and does the following thing:

1. Will check all the people you are following already.
2. Will check all the media of those people and the likes within.
3. Will start following them as well.
4. Thats it, so at the end of the day you be following lot of people.
 
## Steps

1. Clone ```goinstadownload``` repository
2. Go to your $GOPATH 
3. Install the following dependencies 
``` 
go get -v -u github.com/ahmdrz/goinsta
```
4. Create a config file with any name and path and with this structure
```
{ 
	"InstaUser":"your_instagram_user",
	"InstaPass":"your_instagram_pwd",
	"Blacklist":["names_to_exclude"]
}
```
5. Build and run passing the correct parameters
```
./main <Path of your config File> <Time between users> 
> Path of your config File: This is the path to the config file, where user and pwd are stored.
> Time between users: To check each user we should wait some time, this is because we will overload the Instagram Api and maybe we will be disconnected.
```
