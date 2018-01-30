package main

import (
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/urShadow/go-vk-api"
)

// User type from VK api
type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UserList type from VK api
type UserList struct {
	Response []User `json:"responce"`
}

// Variables used for command line parameters
var (
	DiscordToken string
	VkToken      string
	ChannelID    string
)

func init() {
	log.Println("Init")

	flag.StringVar(&DiscordToken, "dt", "", "Discord authentication token")
	flag.StringVar(&VkToken, "vt", "", "Vk token")
	flag.StringVar(&ChannelID, "dcid", "", "Channel ID in Discord")
	flag.Parse()

	log.Println(DiscordToken)
}

func main() {
	log.Println("Bot starting ...")

	discord, err := discordgo.New("Bot " + DiscordToken)
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

	api := vk.New("ru")
	vkerr := api.Init(VkToken)
	if vkerr != nil {
		log.Errorln(vkerr)
	}

	api.OnNewMessage(func(msg *vk.LPMessage) {
		if msg.Flags&vk.FlagMessageOutBox == 0 {

			userData, err := api.CallMethod("users.get", vk.RequestParams{
				"user_ids": strconv.FormatInt(msg.FromID, 10),
			})
			if err == nil {
				res := UserList{}
				if err = json.Unmarshal(userData, &res); err != nil {
					log.Errorf("Decoding JSON. %s", err)
				}

				log.Print("Response: ")
				log.Println(res)

				user := res.Response[0]

				message := user.FirstName
				message += " " + user.LastName
				message += " [" + user.ID + "]\n"
				message += msg.Text

				_, err = discord.ChannelMessageSend(ChannelID, message)
				if err != nil {
					log.Errorln(err)
				}
			} else {
				log.Errorf("Don't get user: %s", err)
			}

			// if msg.Text == "ping" {
			// 	api.Messages.Send(vk.RequestParams{
			// 		"peer_id":          strconv.FormatInt(msg.FromID, 10),
			// 		"message":          "Pong!",
			// 		"forward_messages": strconv.FormatInt(msg.ID, 10),
			// 	})

			// }
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
	if m.Content == "ping" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
		if err != nil {
			log.Errorln(err)
		}
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
