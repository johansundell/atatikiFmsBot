package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/johansundell/atatikiFmsBot/fmsadmin"
	"github.com/kardianos/service"
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
var server fmsadmin.Server

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

	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	svcConfig := &service.Config{
		Name:        "fmsBot",
		DisplayName: "fmsBot",
		Description: "fmsBot is a simple bot",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
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
