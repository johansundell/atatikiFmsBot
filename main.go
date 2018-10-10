package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
	BotID string
)

type commandFunc struct {
	command   string
	helpText  string
	extracter string
	category  category
}

type category string

const filenameSettings = "settings.json"

var settings Config

const (
	categoryStats  category = "=== Stats ==="
	categoryAdmin  category = "=== Admin ==="
	catgoryHelp    category = "=== Help ==="
	categoryHidden category = "=== Hidden ==="
	categorySearch category = "=== Search ==="
	categoryFun    category = "=== Fun ==="
)

var botFuncs map[commandFunc]func(ctx *context.Context, command string) (string, error) = make(map[commandFunc]func(ctx *context.Context, command string) (string, error))
var lockMap = sync.RWMutex{}

func main() {
	ex, err := os.Executable()
	if err != nil {
		fmt.Println(err)
		return
	}
	dir, _ := filepath.Split(ex)
	dat, err := ioutil.ReadFile(dir + filenameSettings)
	if err != nil {
		data, _ := json.Marshal(settings)
		ioutil.WriteFile(dir+filenameSettings, data, 0664)
		fmt.Println("settings.json missing, " + err.Error())
		return
	}

	if err := json.Unmarshal(dat, &settings); err != nil {
		fmt.Println(err)
		return
	}

	Token = settings.Token
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Get the account information.
	u, err := dg.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
		return
	}

	// Store the account ID for later use.
	BotID = u.ID

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println(m.Author.ID, m.Author.Username, m.Author.Username+"#"+m.Author.Discriminator, m.Content)

	// Ignore all messages created by the bot itself
	if m.Author.ID == BotID {
		fmt.Println("me")
		return
	}
	msg := ""
	c, _ := s.Channel(m.ChannelID)
	_ = c
	_ = msg
	if strings.HasPrefix(m.Content, "!") {
		command := strings.ToLower(m.Content)
		if strings.Contains(command, "! ") {
			command = strings.Replace(command, "! ", "!", 1)
		}
		lockMap.RLock()
		defer lockMap.RUnlock()

		ctx := context.WithValue(context.Background(), "sess", s)
		ctx = context.WithValue(ctx, "msg", m)
		for _, v := range botFuncs {
			if str, err := v(&ctx, command); err != nil {
				fmt.Println(err)
			} else {
				msg += str
			}
		}
		s.ChannelMessageSend(m.ChannelID, msg)
	}
}
