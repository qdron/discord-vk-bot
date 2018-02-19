package main

import (
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/urShadow/go-vk-api"
)

// User type from VK api
type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Photo50   string `json:"photo_50"`
	FromID    int64  `json:"from_id"`
}

// UserList type from VK api
type UserList struct {
	Response []User `json:"response"`
}

// GroupData type from VK API
type GroupData struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Photo50 string `json:"photo_50"`
}

// GroupListData type from VK API
type GroupListData struct {
	Response []GroupData `json:"response"`
}

// Variables used for command line parameters
var (
	discordToken string
	VkToken      string
	channelID    string
	api          *vk.VK
	groupID      string
)

var log *logrus.Logger

func newLogger() *logrus.Logger {
	if log != nil {
		return log
	}

	path := "/var/log/discord-vk-bot/bot.log"
	writer, err := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(604800)*time.Second),
	)
	if err != nil {
		return log
	}

	log = logrus.New()
	log.Hooks.Add(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  writer,
			logrus.ErrorLevel: writer,
		},
		&logrus.TextFormatter{},
	))

	return log
}

func init() {
	log = newLogger()

	log.Println("Init")

	flag.StringVar(&VkToken, "vk_token", "", "Vk token")
	flag.StringVar(&groupID, "vk_groupid", "", "VK group id")
	flag.StringVar(&discordToken, "discord_token", "", "Discord authentication token")
	flag.StringVar(&channelID, "discord_channelid", "", "Channel ID in Discord")
	flag.Parse()

	log.Println(discordToken)
}

// Get User stuct by user ID
func getUser(userID string) (User, error) {
	var err error
	var userData []byte
	if userID == "" {
		userData, err = api.CallMethod("users.get", vk.RequestParams{
			"fields": "photo_50, from_id",
		})
	} else {
		userData, err = api.CallMethod("users.get", vk.RequestParams{
			"user_ids": userID,
			"fields":   "photo_50, from_id",
		})
	}

	if err != nil {
		return User{}, err
	}

	log.Println(string(userData))
	res := UserList{}
	if err = json.Unmarshal(userData, &res); err != nil {
		log.Errorf("Decoding JSON. %s", err)
		return User{}, err
	}

	return res.Response[0], err
}

func getGroupByID(groupID string) (GroupData, error) {
	var err error
	var data []byte
	data, err = api.CallMethod("groups.getById", vk.RequestParams{
		"group_id": groupID,
	})
	if err != nil {
		return GroupData{}, err
	}

	log.Debug(string(data))
	res := GroupListData{}
	if err = json.Unmarshal(data, &res); err != nil {
		log.Errorf("Decoding JSON. %s", err)
		return GroupData{}, err
	}

	if len(res.Response) <= 0 {
		return GroupData{}, err
	}

	return res.Response[0], err
}

func main() {
	log.Println("Bot starting ...")

	discord, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Errorln(err)
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	discord.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		log.Errorln("error opening connection,", err)
		return
	}

	api = vk.New("ru")
	vkerr := api.Init(VkToken)
	if vkerr != nil {
		log.Errorln(vkerr)
	}

	group, err := getGroupByID(groupID)

	log.Print(group.Name)

	api.OnNewMessage(func(msg *vk.LPMessage) {
		user, err := getUser(strconv.FormatInt(msg.FromID, 10))
		if err != nil {
			return
		}

		userName := user.FirstName + " " + user.LastName + " [" + strconv.FormatInt(user.ID, 10) + "]"

		var author discordgo.MessageEmbedAuthor
		var footer discordgo.MessageEmbedFooter
		if msg.Flags&vk.FlagMessageOutBox == 0 {
			author = discordgo.MessageEmbedAuthor{
				Name:    userName,
				IconURL: user.Photo50,
			}
		} else {
			author = discordgo.MessageEmbedAuthor{
				Name:    group.Name,
				IconURL: group.Photo50,
			}
			footer = discordgo.MessageEmbedFooter{Text: "для " + userName, IconURL: user.Photo50}
		}

		embed := discordgo.MessageEmbed{Author: &author, Description: msg.Text, Footer: &footer}

		_, err = discord.ChannelMessageSendEmbed(channelID, &embed)
		if err != nil {
			log.Errorln(err)
		}
	})

	// Cleanly close down the Discord session.
	defer discord.Close()

	go api.RunLongPoll()

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Println("Bot finished")
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If the message is "ping" reply with "Pong!"
	if strings.HasPrefix(m.Content, "!ответ") {
		data := strings.Split(m.Content, " ")

		if len(data) < 3 {
			log.Error("No message or profileID")
			return
		}

		userID := data[1]
		message := strings.Join(data[2:], " ")

		api.Messages.Send(vk.RequestParams{
			"peer_id": userID,
			"message": message,
		})
	}
}
