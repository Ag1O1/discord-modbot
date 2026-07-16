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
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	dg.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMembers

	dg.AddHandler(messageCreate)

	dg.AddHandler(interactionHandler)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}
	defer func() {
		err := dg.Close()
		if err != nil {
			fmt.Println("Error closing connection,", err)
		}

	}()

	registerCommands(dg)
	fmt.Println("Bot is now running. Press CTRL+C to exit.")

	select {}
}
