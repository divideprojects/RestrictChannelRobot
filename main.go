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
	// Create bot from environment value.
	b, err := gotgbot.NewBot(botToken, &gotgbot.BotOpts{
		Client: http.Client{},
		DefaultRequestOpts: &gotgbot.RequestOpts{
			Timeout: gotgbot.DefaultTimeout,
			APIURL:  gotgbot.DefaultAPIURL,
		},
	})
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	// Create updater and dispatcher.
	updater := ext.NewUpdater(&ext.UpdaterOpts{
		Dispatcher: ext.NewDispatcher(&ext.DispatcherOpts{
			// If an error is returned by a handler, log it and continue going.
			Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
				log.Println("an error occurred while handling update:", err.Error())
				return ext.DispatcherActionNoop
			},
			MaxRoutines: 0,
		}),
	})
	dispatcher := updater.Dispatcher

	if enableWebhook {
		log.Println("[Webhook] Starting webhook...")
		webhookOpts := ext.WebhookOpts{
			ListenAddr:  "localhost:8080", // This example assumes you're in a dev environment running ngrok on 8080.
			SecretToken: webhookSecret,    // Setting a webhook secret here allows you to ensure the webhook is set by you (must be set here AND in SetWebhook!).
		}

		// We use the token as the urlPath for the webhook, as using a secret ensures that strangers aren't crafting fake updates.
		err = updater.StartWebhook(b, botToken, webhookOpts)
		if err != nil {
			panic("failed to start webhook: " + err.Error())
		}

		err = updater.SetAllBotWebhooks(webhookDomain, &gotgbot.SetWebhookOpts{
			MaxConnections:     100,
			DropPendingUpdates: true,
			SecretToken:        webhookOpts.SecretToken,
		})

		if err != nil {
			log.Fatalf("failed to set webhook: %s\n", err.Error())
		} else {
			log.Printf("[Webhook] Set Webhook to: %s\n", webhookDomain)
		}

		log.Println("[Webhook] Webhook started Successfully!")
	} else {
		if _, err = b.DeleteWebhook(nil); err != nil {
			log.Fatalf("[Polling] Failed to remove webhook: %v", err)
		}
		log.Println("[Polling] Removed Webhook!")
		err = updater.StartPolling(b, &ext.PollingOpts{DropPendingUpdates: false})
		if err != nil {
			log.Fatalf("[Polling] Failed to start polling: %s\n", err.Error())
		}
		log.Println("[Polling] Started Polling...!")
	}

	// log msg telling that bot has started
	log.Printf("%s has been started...!\nMade with ❤️ by @DivideProjects\n", b.User.Username)

	// Handlers for running commands.
	dispatcher.AddHandler(handlers.NewCommand("start", start))
	dispatcher.AddHandler(handlers.NewCommand("help", help))
	dispatcher.AddHandler(handlers.NewCommand("source", source))
	dispatcher.AddHandler(handlers.NewCommand("ignore", ignoreChannel))
	dispatcher.AddHandler(handlers.NewCommand("unignore", unignoreChannel))
	dispatcher.AddHandler(handlers.NewCommand("ignorelist", ignoreList))
	dispatcher.AddHandlerToGroup(
		handlers.NewMessage(
			func(msg *gotgbot.Message) bool {
				return msg.GetSender().IsAnonymousChannel()
			},
			restrictChannels,
		),
		-1,
	)

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
			"Head towards @DivideSupport for any queries!",
		user.FirstName, bot.FirstName,
	)
	kb = gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
			{

				{
					Text: "Support",
					Url:  "https://t.me/DivideSupport",
				},
				{
					Text: "Channel",
					Url:  "https://t.me/DivideProjects",
				},
			},
			{
				{
					Text: "Source",
					Url:  "https://github.com/divideprojects/RestrictChannelRobot",
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

	return ext.EndGroups
}

func help(bot *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	chat := ctx.EffectiveChat

	var text string

	// stay silent in group chats
	if chat.Type != "private" {
		return ext.EndGroups
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

	return ext.EndGroups
}

func source(bot *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	chat := ctx.EffectiveChat

	var text string

	// stay silent in group chats
	if chat.Type != "private" {
		return ext.EndGroups
	}

	text = fmt.Sprintf(
		"You can find my source code by <a href=\"%s\">here</a> or by clicking the button below.",
		"https://github.com/divideprojects/RestrictChannelRobot",
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
							Url:  "https://github.com/divideprojects/RestrictChannelRobot",
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

	return ext.EndGroups
}

func ignoreChannel(bot *gotgbot.Bot, ctx *ext.Context) error {

	msg := ctx.EffectiveMessage
	chat := ctx.EffectiveChat
	user := ctx.EffectiveSender

	if !isUserAdmin(bot, chat.Id, user.Id()) {
		_, _ = msg.Reply(bot, "This command can only be used by admins!", nil)
		return ext.EndGroups
	}

	if chat.Type != "supergroup" {
		_, _ = msg.Reply(bot, "This command can only be used in Groups.", nil)
		return ext.EndGroups
	}

	channelId, err := extractChannelId(msg)

	if channelId == -1 {
		_, _ = msg.Reply(bot, "Please reply to a message from a channel or pass the channel id to add a user to ignore list.", nil)
		return ext.EndGroups
	}

	if err != nil {
		_, _ = msg.Reply(bot, "Failed to extract channel id: "+err.Error(), nil)
		return err
	}

	ignoreSettings := getIgnoreSettings(chat.Id)
	for _, i := range ignoreSettings.IgnoredChannels {
		if channelId == i {
			_, _ = msg.Reply(bot, "This channel is already in ignore list.", nil)
		}
	}

	ignoreChat(chat.Id, channelId)
	_, _ = msg.Reply(bot, "Added this channel to ignore list.", nil)

	return ext.EndGroups
}

func unignoreChannel(bot *gotgbot.Bot, ctx *ext.Context) error {

	msg := ctx.EffectiveMessage
	chat := ctx.EffectiveChat
	user := ctx.EffectiveSender

	if !isUserAdmin(bot, chat.Id, user.Id()) {
		_, _ = msg.Reply(bot, "This command can only be used by admins!", nil)
		return ext.EndGroups
	}
	if chat.Type != "supergroup" {
		_, _ = msg.Reply(bot, "This command can only be used in Groups.", nil)
		return ext.EndGroups
	}

	channelId, err := extractChannelId(msg)

	if channelId == -1 {
		_, _ = msg.Reply(bot, "Please reply to a message from a channel or pass the channel id to add a user to ignore list.", nil)
		return ext.EndGroups
	}

	if err != nil {
		_, _ = msg.Reply(bot, "Failed to extract channel id: "+err.Error(), nil)
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

	_, _ = msg.Reply(bot, "This channel is not in ignore list.", nil)

	return ext.EndGroups
}

func ignoreList(bot *gotgbot.Bot, ctx *ext.Context) error {

	msg := ctx.EffectiveMessage
	chat := ctx.EffectiveChat

	if chat.Type != "supergroup" {
		_, _ = msg.Reply(bot, "This command can only be used in Groups.", nil)
		return ext.EndGroups
	}

	var text string

	ignoreList := getIgnoreSettings(chat.Id).IgnoredChannels

	if len(ignoreList) > 0 {
		text = "Here is the list of channels currently ignored by me:"
		for _, i := range ignoreList {
			text += fmt.Sprintf("\n - <code>%d</code>", i)
		}
	} else {
		text = "There are no channels in ignore list."
	}

	_, _ = msg.Reply(bot, text, &gotgbot.SendMessageOpts{ParseMode: "HTML"})

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
