package main

import (
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// this function will extract the
func extractChannelId(msg *gotgbot.Message) (channelId int64, err error) {

	args := strings.Split(msg.Text, " ")

	if msg.ReplyToMessage != nil && msg.ReplyToMessage.SenderChat != nil && len(args) == 1 {
		channelId = msg.ReplyToMessage.SenderChat.Id
	} else {
		if len(args) > 1 {
			if strings.HasPrefix(args[1], "-100") {
				channelId, err = strconv.ParseInt(args[1], 10, 64)
				if err != nil {
					return 0, err
				}
			}
		} else {
			return -1, err
		}

	}

	return channelId, err
}

func isUserAdmin(bot *gotgbot.Bot, chatID, userId int64) bool {
	// Placing this first would not make additional queries if check is success!
	if userId == 1087968824 {
		return true
	}

	chat, err := bot.GetChat(chatID, nil)
	if err != nil {
		log.Errorf("[isUserAdmin]: %v", err)
		return false
	}

	if chat.Type == "private" {
		return true
	}

	var adminlist = make([]int64, 0)

	adminsL, _ := chat.GetAdministrators(bot, nil)

	for _, admin := range adminsL {
		adminlist = append(adminlist, admin.GetUser().Id)
	}

	return findInInt64Slice(adminlist, userId)
}

func findInInt64Slice(slice []int64, val int64) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
