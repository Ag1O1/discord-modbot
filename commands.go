package main

import (
	"fmt"
	"log"
	"strings"
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
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "count-channel",
				Description: "Set the channel where the counter appers",
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
	config := getGuildConfig(i.GuildID)
	var msg strings.Builder

	for _, option := range data.Options {
		switch option.Name {
		case "channel-to-monitor":
			channel := option.ChannelValue(s)
			config.ChannelToMonitor = channel.ID
			fmt.Fprintf(&msg, "Set channel to monitor to %s\n", channel.Name)
		case "log-channel":
			channel := option.ChannelValue(s)
			config.LogChannelID = channel.ID
			fmt.Fprintf(&msg, "Set channel to monitor to %s\n", channel.Name)
		case "count-channel":
			channel := option.ChannelValue(s)
			config.BanCountChannelID = channel.ID
			if _, err := s.ChannelMessageSend(config.BanCountChannelID, fmt.Sprintf("Ban count: %v", config.BanCount.Load())); err != nil {
				sendLog(s, i.GuildID, "Unable to create count message")
				return
			}
			fmt.Fprintf(&msg, "Set channel to put count message to %s\n", channel.Name)
		}
	}

	respond(s, i, msg.String())
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
