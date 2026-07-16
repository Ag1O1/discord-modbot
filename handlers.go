package main

import (
	"github.com/bwmarrin/discordgo"
)

func interactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch i.ApplicationCommandData().Name {
	case "ping":
		handlePing(s, i)
	case "config":
		handleConfig(s, i)
	default:
		handleUnknown(s, i)
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	guildID := m.GuildID
	channelID := m.ChannelID
	userID := m.Author.ID

	if userID == s.State.User.ID {
		return
	}
	if channelID != getGuildConfig(guildID).ChannelToMonitor {
		return
	}
	if isAdmin(s, guildID, userID, channelID) {
		return
	}

	banUser(s, guildID, userID, 1)
}
