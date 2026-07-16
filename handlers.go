package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func interactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
	if !isAdmin(s, i.GuildID, i.Member.User.ID, i.ChannelID) {
		return
	}
	if i.Member == nil {
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
	if _, err := s.ChannelMessageSend(channelID, fmt.Sprintf("Banned user: %s", getUsername(s, userID))); err != nil {
		sendLog(s, guildID, fmt.Sprintf("Failed to send ban message: Banned user: %s", getUsername(s, userID)))
	}

}

func guildCreate(s *discordgo.Session, g *discordgo.GuildCreate) {
	err := s.RequestGuildMembers(g.ID, "", 0, "", false)
	if err != nil {
		fmt.Println("Failed to request guild members:", err)
	}
}
