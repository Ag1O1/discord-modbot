package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "Pings the bot and returns the latency",
	},

	{
		Name:        "config",
		Description: "Manage config options",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel-to-monitor",
				Description: "Set the honeypot channel",
			},
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "log-channel",
				Description: "Set the log channel",
			},
		},
	},
}

func registerCommands(s *discordgo.Session) {
	for _, guild := range s.State.Guilds {
		_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, guild.ID, commands)
		if err != nil {
			log.Printf("Failed to register Commands: %v", err)
		}
	}
}

func handlePing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	gatewayLatency := s.HeartbeatLatency()
	start := time.Now()
	respond(s, i, "Pinging...")
	apiLatency := time.Since(start)
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: new(fmt.Sprintf("Pong!\nGateway latency:%v\nAPI latency:%v", gatewayLatency, apiLatency)),
	})
	if err != nil {
		fmt.Println("Failed to edit message", i.ApplicationCommandData().Name)
	}
}

func handleConfig(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	option := data.Options[0]
	var msg string

	switch option.Name {
	case "channel-to-monitor":
		channel := option.ChannelValue(s)
		getGuildConfig(i.GuildID).ChannelToMonitor = channel.ID
		msg = fmt.Sprintf("Set channel to monitor to %s", channel.Name)
	case "log-channel":
		channel := option.ChannelValue(s)
		getGuildConfig(i.GuildID).LogChannelID = channel.ID
		msg = fmt.Sprintf("Set log channel to %s", channel.Name)
	}
	respond(s, i, msg)
	saveConfig()
}

func handleUnknown(s *discordgo.Session, i *discordgo.InteractionCreate) {
	respond(s, i, "Error: unkown command")
}

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	if err != nil {
		fmt.Println("Failed to respond to interaction: ", i.ApplicationCommandData().Name)
	}
}
