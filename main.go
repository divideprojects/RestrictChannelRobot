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
			APIURL:      apiUrl,
			Client:      http.Client{},
			GetTimeout:  gotgbot.DefaultGetTimeout,
			PostTimeout: gotgbot.DefaultPostTimeout,
		},
	)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	// Create updater and dispatcher.
	updater := ext.NewUpdater(nil)
	dispatcher := updater.Dispatcher

	// Handlers for runnning commands.
	dispatcher.AddHandler(handlers.NewCommand("start", start))
	dispatcher.AddHandler(handlers.NewCommand("help", help))
	dispatcher.AddHandler(handlers.NewCommand("source", source))
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
		fmt.Println("[Webhook] Starting webhook...")

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

		fmt.Printf("[Webhook] Set Webhook to: %s\n", webhookUrl)

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

		fmt.Println("[Webhook] Webhook started Successfully!")
	} else {
		err = updater.StartPolling(b, &ext.PollingOpts{DropPendingUpdates: false})
		if err != nil {
			log.Fatalf("[Polling] Failed to start polling: %s\n", err.Error())
		}
		log.Println("[Polling] Started Polling...!")
	}

	// log msg telling that bot has started
	fmt.Printf("%s has been started...!\nMade with ❤️ by @DivideProjects\n", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}

func start(bot *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	user := ctx.EffectiveSender.User
	chat := ctx.EffectiveChat

	var text string

	// stay silent in group chats
	if chat.Type != "private" {
		return ext.EndGroups
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

	_, err := msg.Reply(
		bot,
		text,
		&gotgbot.SendMessageOpts{
			ParseMode:             "HTML",
			DisableWebPagePreview: true,
		},
	)
	if err != nil {
		fmt.Println("[Start] Failed to reply:", err.Error())
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
		"Just add me to a group with these basic permissions and I'll do the rest!\n",
		" - Ban Permissions: To ban the channels\n",
		" - Delete Message Permissions: To delete the messages sent by channel",
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
		fmt.Println("[Start] Failed to reply:", err.Error())
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
		fmt.Println("[Start] Failed to reply:", err.Error())
		return err
	}

	return ext.EndGroups
}

func restrictChannels(bot *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	chat := ctx.EffectiveChat
	sender := ctx.EffectiveSender

	_, err := msg.Delete(bot)
	if err != nil {
		_, err := msg.Reply(bot, "Failed to delete message: "+err.Error(), nil)
		if err != nil {
			fmt.Println("[RestrictChannels] Failed to reply:", err.Error())
		}
		fmt.Println("[RestrictChannels] Failed to delete message:", err.Error())
		return err
	}

	_, err = chat.BanSenderChat(bot, sender.Id())
	if err != nil {
		_, err := msg.Reply(bot, "Failed to ban sender: "+err.Error(), nil)
		if err != nil {
			fmt.Println("[RestrictChannels] Failed to reply:", err.Error())
		}
		fmt.Println("[RestrictChannels] Failed to ban sender:", err.Error())
		return err
	}

	fmt.Printf("[RestrictChannels] Banning %s (%d) in %s (%d)\n", sender.Name(), sender.Id(), chat.Title, chat.Id)

	return ext.ContinueGroups
}
