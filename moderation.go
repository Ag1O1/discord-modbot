package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func banUser(s *discordgo.Session, guildID, userID string, days int) {
	userName := getUsername(s, userID)
	if err := s.GuildBanCreate(guildID, userID, days); err != nil {
		sendLog(s, guildID, fmt.Sprintf("Error: unable to ban user %s: %v", userName, err))
		return
	}
	if err := s.GuildBanDelete(guildID, userID); err != nil {
		sendLog(s, guildID, fmt.Sprintf("Error: unable to unban user %s: %v", userName, err))
		return
	}
	updateCounter(s, guildID)
}

func isAdmin(s *discordgo.Session, guildID, userID, channelID string) bool {
	perms, err := s.State.UserChannelPermissions(userID, channelID)
	if err != nil {
		userName := getUsername(s, userID)
		sendLog(s, guildID, fmt.Sprintf("Error: failed to check permissions for user %s: %v", userName, err))
		return false
	}
	return perms&discordgo.PermissionAdministrator != 0
}

func getUsername(s *discordgo.Session, userID string) string {
	user, err := s.User(userID)
	if err != nil {
		return "Unknown"
	}
	return user.Username
}

func getGuildConfig(guildID string) *GuildConfig {
	configsMu.RLock()
	config, ok := configs[guildID]
	configsMu.RUnlock()

	if ok {
		return config
	}

	configsMu.Lock()
	defer configsMu.Unlock()

	config, ok = configs[guildID]
	if !ok {
		config = &GuildConfig{}
		configs[guildID] = config
	}
	return config
}

func sendLog(s *discordgo.Session, guildID, message string) {
	if _, err := s.ChannelMessageSend(getGuildConfig(guildID).LogChannelID, message); err != nil {
		log.Printf("Failed to send discord log: %v", err)
	}
}
func updateCounter(s *discordgo.Session, guildID string) {
	config := getGuildConfig(guildID)
	config.BanCount.Add(1)

	if config.BanCountChannelID == "" {
		sendLog(s, guildID, "Ban count channel ID not set")
		return
	}
	if config.BanCountMessageID == "" {
		if countMessage, err := s.ChannelMessageSend(config.BanCountChannelID, "Ban count: 0"); err != nil {
			sendLog(s, guildID, "Unable to create count message")
			return
		} else {
			config.BanCountMessageID = countMessage.ID
		}
	}
	if _, err := s.ChannelMessageEdit(config.BanCountChannelID, config.BanCountMessageID, fmt.Sprintf("Ban count: %v", config.BanCount.Load())); err != nil {
		sendLog(s, guildID, "Unable to edit count message")
	}
}
