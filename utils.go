package main

import (
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// this function will extract the channel id from the message
func extractChannelId(msg *gotgbot.Message) (channelId int64, err error) {

	// split the message into arguments
	args := strings.Split(msg.Text, " ")

	// if the message is a reply to a message from a channel, then the channel id will be extracted from the reply
	// else if the message is a reply to a message from a user, then the channel id will be extracted from the message
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

	// return the channel id and nil error if no error occurs
	return channelId, err
}

// this function will check if the user is an admin or not
func isUserAdmin(bot *gotgbot.Bot, chatID, userId int64) bool {
	// Placing this first would not make additional queries if check is success!
	if userId == 1087968824 {
		return true
	}

	// Get chat info using chatID
	chat, err := bot.GetChat(chatID, nil)
	if err != nil {
		log.Errorf("[isUserAdmin]: %v", err)
		return false
	}

	// If chat is private, then user is admin
	if chat.Type == "private" {
		return true
	}

	// make a list of admins
	var adminlist = make([]int64, 0)

	// get the list of admins
	adminsL, _ := chat.GetAdministrators(bot, nil)

	// append the admin id to the list
	for _, admin := range adminsL {
		adminlist = append(adminlist, admin.GetUser().Id)
	}

	// check if the user is in the list of admins
	return findInInt64Slice(adminlist, userId)
}

// this function checks if the int is there in the slice or not
func findInInt64Slice(slice []int64, val int64) bool {
	// loop through the slice and check if the value is there or not
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	// return false if the value is not there in the slice
	return false
}
