package main

import (
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// this function will extract the 
func extractChannelId(msg *gotgbot.Message) (channelId int64, err error) {

	args := strings.Split(msg.Text, " ")

	if msg.ReplyToMessage != nil && msg.ReplyToMessage.SenderChat != nil && len(args) == 1 {
		channelId = msg.ReplyToMessage.SenderChat.Id
	} else {
		if strings.HasPrefix(args[1], "-100") {
			return msg.SenderChat.Id, err
		}

		channelId, err = strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return 0, err
		}
	}

	return channelId, err
}
