package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func main() {
	b, err := gotgbot.NewBot(
		botToken,
		&gotgbot.BotOpts{
			DefaultRequestOpts: &gotgbot.RequestOpts{
				Timeout: gotgbot.DefaultTimeout,
				APIURL:  apiUrl,
			},
			Client: http.Client{},
		},
	)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	// Create updater and dispatcher.
	updater := ext.NewUpdater(nil)
	dispatcher := updater.Dispatcher

	// Handlers for running commands.
	dispatcher.AddHandler(handlers.NewCommand("start", start))
	dispatcher.AddHandler(handlers.NewCommand("help", help))
	dispatcher.AddHandler(handlers.NewCommand("source", source))
	dispatcher.AddHandler(handlers.NewCommand("ignore", ignoreChannel))
	dispatcher.AddHandler(handlers.NewCommand("unignore", unignoreChannel))
	dispatcher.AddHandler(handlers.NewCommand("ignorelist", ignoreList))
	dispatcher.AddHandler(handlers.NewCommand("start", start))
	dispatcher.AddHandlerToGroup(
		handlers.NewMessage(
			func(msg *gotgbot.Message) bool {
				return msg.GetSender().IsAnonymousChannel()
			},
			restrictChannels,
		),
		-1,
	)

	if enableWebhook {
		log.Println("[Webhook] Starting webhook...")

		// Set Webhook
		ok, err := b.SetWebhook(
			webhookUrl,
			&gotgbot.SetWebhookOpts{
				MaxConnections: 50,
			},
		)

		if !ok && err != nil {
			log.Fatalf("[Webhook] Failed to set webhook: %s", err.Error())
		}

		log.Printf("[Webhook] Set Webhook to: %s\n", webhookUrl)

		// Start the webhook
		err = updater.StartWebhook(b,
			ext.WebhookOpts{
				Listen:  "0.0.0.0",
				Port:    webhookPort,
				URLPath: botToken,
			},
		)
		if err != nil {
			log.Fatalf("[Webhook] Failed to start webhook: %s", err.Error())
		}

		log.Println("[Webhook] Webhook started Successfully!")
	} else {
		err = updater.StartPolling(b, &ext.PollingOpts{DropPendingUpdates: false})
		if err != nil {
			log.Fatalf("[Polling] Failed to start polling: %s\n", err.Error())
		}
		log.Println("[Polling] Started Polling...!")
	}

	// log msg telling that bot has started
	log.Printf("%s has been started...!\nMade with ❤️ by @DivideProjects\n", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}

func start(bot *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	user := ctx.EffectiveSender.User
	chat := ctx.EffectiveChat

	var text string
	var kb gotgbot.InlineKeyboardMarkup

	// stay silent in group chats
	if chat.Type != "private" {
		return nil
	}

	text = fmt.Sprintf(
		"<b>Hi %s</b>,\n"+
			"I am <b>%s</b>, a Telegram Bot made using <a href=\"https://go.dev\">Go</a>\n\n"+
			"Send /help for getting info on how on use me!\n"+
			"Also you can send /source to get my source code to know how i'm built ;) and make sure to give a star to it; that makes my Devs to work more on O.S. projects like me :)\n\n"+
			"Hope you liked it !\n"+
			"Brought to You with ❤️ By @DivideProjects\n"+
			"Head towards @DivideProjectsDiscussion for any queries!",
		user.FirstName, bot.FirstName,
	)
	kb = gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{

				{
					Text: "Support",
					Url:  "https://t.me/DivideProjectsDiscussion",
				},
				{
					Text: "Channel",
					Url:  "https://t.me/DivideProjects",
				},
			},
			{
				{
					Text: "Source",
					Url:  "https://github.com/DivideProjects/RestrictChannelRobot",
				},
			},
		},
	}

	_, err := msg.Reply(
		bot,
		text,
		&gotgbot.SendMessageOpts{
			ParseMode:             "HTML",
			DisableWebPagePreview: true,
			ReplyMarkup:           kb,
		},
	)
	if err != nil {
		log.Println("[Start] Failed to reply:", err.Error())
		return err
	}

	return nil
}

func help(bot *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	chat := ctx.EffectiveChat

	var text string

	// stay silent in group chats
	if chat.Type != "private" {
		return nil
	}

	text = fmt.Sprint(
		"Add me to your group with the following permissions and I'll handle the rest!",
		"\n - Ban Permissions: To ban the channels",
		"\n - Delete Message Permissions: To delete the messages sent by channel",

		"\n\n<b>Some Tips:</b>",
		"\n1. To ignore a channel use /ignore by replying a message from that channel or you can pass a channel id. for more help type /ignore.",
		"\n2. To unignore a channel use /unignore by replying a message from that channel or you can pass a channel id. for more help type /unignore.",
		"\n3. To get the list of all ignored channel use /ignorelist.",

		"\n\n<b>Available Commands:</b>",
		"\n/start - ✨ display start message.",
		"\n/ignore - ✅ unban and allow that user to sending message as channel (admin only).",
		"\n/ignorelist - 📋 get list ignored channel.",
		"\n/unignore - ⛔️ ban an unallow that user to sending message as channel (admin only).",
		"\n/source - 📚 get source code.",
	)

	_, err := msg.Reply(
		bot,
		text,
		&gotgbot.SendMessageOpts{
			ParseMode:             "HTML",
			DisableWebPagePreview: true,
		},
	)
	if err != nil {
		log.Println("[Start] Failed to reply:", err.Error())
		return err
	}

	return nil
}

func source(bot *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	chat := ctx.EffectiveChat

	var text string

	// stay silent in group chats
	if chat.Type != "private" {
		return nil
	}

	text = fmt.Sprintf(
		"You can find my source code by <a href=\"%s\">here</a> or by clicking the button below.",
		"https://github.com/DivideProjects/RestrictChannelRobot",
	)

	_, err := msg.Reply(
		bot,
		text,
		&gotgbot.SendMessageOpts{
			ParseMode:             "HTML",
			DisableWebPagePreview: true,
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
					{
						{
							Text: "Source Code",
							Url:  "https://github.com/DivideProjects/RestrictChannelRobot",
						},
					},
				},
			},
		},
	)
	if err != nil {
		log.Println("[Start] Failed to reply:", err.Error())
		return err
	}

	return nil
}

func ignoreChannel(bot *gotgbot.Bot, ctx *ext.Context) error {

	msg := ctx.EffectiveMessage
	chat := ctx.EffectiveChat

	if chat.Type != "supergroup" {
		msg.Reply(bot, "This command can only be used in Groups.", nil)
		return nil
	}

	channelId, err := extractChannelId(msg)

	if channelId == -1 {
		msg.Reply(bot, "Please reply to a message from a channel or pass the channel id to add a user to ignore list.", nil)
		return ext.EndGroups
	}

	if err != nil {
		msg.Reply(bot, "Failed to extract channel id: "+err.Error(), nil)
		return err
	}

	ignoreSettings := getIgnoreSettings(chat.Id)
	for _, i := range ignoreSettings.IgnoredChannels {
		if channelId == i {
			msg.Reply(bot, "This channel is already in ignore list.", nil)
		}
	}

	ignoreChat(chat.Id, channelId)
	msg.Reply(bot, "Added this channel to ignore list.", nil)

	return ext.EndGroups
}

func unignoreChannel(bot *gotgbot.Bot, ctx *ext.Context) error {

	msg := ctx.EffectiveMessage
	chat := ctx.EffectiveChat

	if chat.Type != "supergroup" {
		msg.Reply(bot, "This command can only be used in Groups.", nil)
		return nil
	}

	channelId, err := extractChannelId(msg)

	if channelId == -1 {
		msg.Reply(bot, "Please reply to a message from a channel or pass the channel id to add a user to ignore list.", nil)
		return ext.EndGroups
	}

	if err != nil {
		msg.Reply(bot, "Failed to extract channel id: "+err.Error(), nil)
		return err
	}

	ignoreSettings := getIgnoreSettings(chat.Id)
	for _, i := range ignoreSettings.IgnoredChannels {
		if channelId == i {
			unignoreChat(chat.Id, channelId)
			msg.Reply(bot, "Removed this channel from ignore list.", nil)
			return ext.EndGroups
		}
	}

	msg.Reply(bot, "This channel is not in ignore list.", nil)

	return ext.EndGroups
}

func ignoreList(bot *gotgbot.Bot, ctx *ext.Context) error {

	msg := ctx.EffectiveMessage
	chat := ctx.EffectiveChat

	if chat.Type != "supergroup" {
		msg.Reply(bot, "This command can only be used in Groups.", nil)
		return nil
	}

	var text string

	ignoreList := getIgnoreSettings(chat.Id).IgnoredChannels

	if len(ignoreList) > 1 {
		text = fmt.Sprintf(
			"Here is the list of channels currently ignored by me:",
		)
		for _, i := range ignoreList {
			text += fmt.Sprintf("\n - <code>%d</code>", i)
		}
	} else {
		text = "There are no channels in ignore list."
	}

	msg.Reply(bot, text, nil)

	return ext.EndGroups
}

func restrictChannels(bot *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	chat := ctx.EffectiveChat
	sender := ctx.EffectiveSender
	senderId := sender.Id()

	// if channel is in ignorelist, then return
	ignoreList := getIgnoreSettings(chat.Id).IgnoredChannels
	for _, i := range ignoreList {
		if i == senderId {
			return ext.ContinueGroups
		}
	}

	_, err := msg.Delete(bot, nil)
	if err != nil {
		log.Println("[RestrictChannels] Failed to delete message:", err.Error())
		return err
	}

	_, err = chat.BanSenderChat(bot, sender.Id(), nil)
	if err != nil {
		log.Println("[RestrictChannels] Failed to ban sender:", err.Error())
		return err
	}

	log.Printf("[RestrictChannels] Banning %s (%d) in %s (%d)\n", sender.Name(), sender.Id(), chat.Title, chat.Id)

	return ext.ContinueGroups
}
