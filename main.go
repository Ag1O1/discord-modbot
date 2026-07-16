package main

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/bwmarrin/discordgo"
)

type GuildConfig struct {
	ChannelToMonitor  string
	BanCountChannelID string
	BanCountMessageID string
	BanCount          atomic.Uint32
	LogChannelID      string
}

var configs = make(map[string]*GuildConfig)
var configsMu sync.RWMutex

func main() {
	loadConfig()
	token := os.Getenv("DISCORD_TOKEN")

	if token == "" {
		fmt.Println("No token provided. Set the DISCORD_TOKEN environment variable.")
		return
	}
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	s.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMembers

	s.AddHandler(messageCreate)

	s.AddHandler(interactionHandler)

	s.AddHandler(guildCreate)

	err = s.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}
	defer func() {
		err := s.Close()
		if err != nil {
			fmt.Println("Error closing connection,", err)
		}

	}()

	registerCommands(s)
	fmt.Println("Bot is now running. Press CTRL+C to exit.")

	select {}
}
