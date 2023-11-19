package bot

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	botAPI *tgbotapi.BotAPI
}

func New(botAPI *tgbotapi.BotAPI) *Bot {
	return &Bot{botAPI: botAPI}
}

func (b *Bot) HandleUpdates(ctx context.Context, updatesChannel tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			return
		case update, ok := <-updatesChannel:
			if !ok {
				return
			}
			if update.Message == nil { // ignore any non-Message updates
				continue
			}

			if !update.Message.IsCommand() { // ignore any non-command Messages
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			// Extract the command from the Message.
			switch update.Message.Command() {
			case "help":
				msg.Text = "I understand /sayhi and /status."
			case "sayhi":
				msg.Text = "Hi :)"
			case "status":
				msg.Text = "I'm ok."
			default:
				msg.Text = "I don't know that command"
			}

			if _, err := b.botAPI.Send(msg); err != nil {
				log.Panic(err)
			}
		}
	}
}
