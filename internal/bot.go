package bot

import (
	"context"
	"log"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	botAPI          *tgbotapi.BotAPI
	wg              *sync.WaitGroup
	workers         int
	fetchBuffer     int
	shutdownChannel chan interface{}
}

func New(botAPI *tgbotapi.BotAPI) *Bot {
	return &Bot{
		botAPI:          botAPI,
		wg:              &sync.WaitGroup{},
		workers:         16,
		fetchBuffer:     100,
		shutdownChannel: make(chan interface{}),
	}
}

func (b *Bot) Start(ctx context.Context) {
	log.Println("starting bot")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10
	ch := make(chan *tgbotapi.Update, b.fetchBuffer)
	go b.FetchAsync(ctx, u, ch)
	b.StartWorkers(ch)
	log.Println("starting bot done")
}

func (b *Bot) Stop() {
	log.Println("stopping bot")
	close(b.shutdownChannel)
	b.wg.Wait()
	log.Println("stopping bot done")
}

func (bot *Bot) FetchAsync(ctx context.Context, config tgbotapi.UpdateConfig, ch chan *tgbotapi.Update) {
	for {
		select {
		case <-bot.shutdownChannel:
			close(ch)
			return
		case <-ctx.Done():
			close(ch)
			return
		default:
			updates, err := bot.botAPI.GetUpdates(config)
			if err != nil {
				log.Println(err)
				log.Println("Failed to get updates, retrying in 3 seconds...")
				time.Sleep(time.Second * 3)

				continue
			}

			for _, update := range updates {
				if update.UpdateID >= config.Offset {
					config.Offset = update.UpdateID + 1
					ch <- &update
				}
			}
		}
	}
}

func (b *Bot) StartWorkers(ch chan *tgbotapi.Update) {
	for i := 0; i < b.workers; i++ {
		b.wg.Add(1)

		go func() {
			for update := range ch {
				if update.Message == nil { // ignore any non-Message updates
					continue
				}

				if !update.Message.IsCommand() { // ignore any non-command Messages
					continue
				}

				msg := b.handleCommand(update.Message)

				if _, err := b.botAPI.Send(msg); err != nil {
					log.Panic(err)
				}
			}
			b.wg.Done()
		}()
	}
}

func (b *Bot) handleCommand(message *tgbotapi.Message) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")

	// Extract the command from the Message.
	switch message.Command() {
	case "help":
		msg.Text = "I understand /sayhi and /status."
	case "sayhi":
		msg.Text = "Hi :)"
	case "status":
		msg.Text = "I'm ok."
	default:
		msg.Text = "I don't know that command"
	}

	return msg
}
