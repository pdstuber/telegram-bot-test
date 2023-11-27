package bot

import (
	"bytes"
	"context"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pdstuber/telegram-bot-test/internal/predict/prediction/model"
	"golang.org/x/image/draw"
)

const (
	defaultColorChannels = 3
	targetImageSize      = 256
	targetImageMime      = "image/jpeg"
	targetImageQuality   = 99
)

// A ImagePredictor predicts the class of an image
type ImagePredictor interface {
	PredictImage(imageBytes []byte, colorChannels int64) (*model.Result, error)
}

type Bot struct {
	botAPI          *tgbotapi.BotAPI
	wg              *sync.WaitGroup
	workers         int
	fetchBuffer     int
	shutdownChannel chan interface{}
	imagePredictor  ImagePredictor
	httpClient      *http.Client
}

func New(botAPI *tgbotapi.BotAPI, imagePredictor ImagePredictor) *Bot {
	return &Bot{
		botAPI:          botAPI,
		wg:              &sync.WaitGroup{},
		workers:         16,
		fetchBuffer:     100,
		shutdownChannel: make(chan interface{}),
		httpClient:      http.DefaultClient,
		imagePredictor:  imagePredictor,
	}
}

func (b *Bot) Start(ctx context.Context) {
	log.Println("starting bot")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10
	ch := make(chan *tgbotapi.Update, b.fetchBuffer)
	go b.FetchAsync(ctx, u, ch)
	b.StartWorkers(ctx, ch)
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

func (b *Bot) StartWorkers(ctx context.Context, ch chan *tgbotapi.Update) {
	for i := 0; i < b.workers; i++ {
		b.wg.Add(1)

		go func() {
			for update := range ch {
				if update.Message == nil { // ignore any non-Message updates
					continue
				}

				if len(update.Message.Photo) == 0 { // ignore any messages not containing a photo
					continue
				}

				msg := b.handlePhoto(ctx, update.Message)

				if _, err := b.botAPI.Send(msg); err != nil {
					log.Panic(err)
				}
			}
			b.wg.Done()
		}()
	}
}

func (b *Bot) handlePhoto(ctx context.Context, message *tgbotapi.Message) tgbotapi.MessageConfig {
	fileConfig := tgbotapi.FileConfig{
		FileID: message.Photo[2].FileID,
	}

	file, err := b.botAPI.GetFile(fileConfig)
	if err != nil {
		log.Println(err)
		return tgbotapi.NewMessage(message.Chat.ID, "There was an error, I'm sorry :(")
	}

	link := file.Link(b.botAPI.Token)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		log.Println(err)
		return tgbotapi.NewMessage(message.Chat.ID, "There was an error, I'm sorry :(")
	}

	response, err := b.httpClient.Do(request)
	if err != nil {
		log.Println(err)
		return tgbotapi.NewMessage(message.Chat.ID, "There was an error, I'm sorry :(")
	}

	defer response.Body.Close()

	photoBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return tgbotapi.NewMessage(message.Chat.ID, "There was an error, I'm sorry :(")
	}

	log.Printf("photo size: %d from url: %s\n", len(photoBytes), link)

	resizedBytes, err := resizeImage(photoBytes)
	if err != nil {
		log.Println(err)
		return tgbotapi.NewMessage(message.Chat.ID, "There was an error, I'm sorry :(")
	}
	result, err := b.imagePredictor.PredictImage(resizedBytes, defaultColorChannels)

	if err != nil {
		log.Println(err)
		return tgbotapi.NewMessage(message.Chat.ID, "There was an error, I'm sorry :(")
	}

	return tgbotapi.NewMessage(message.Chat.ID, result.String())
}

func resizeImage(imageBytes []byte) ([]byte, error) {

	src, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}

	// Set the expected size that you want:
	dst := image.NewRGBA(image.Rect(0, 0, targetImageSize, targetImageSize))

	// Resize:
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	var buf bytes.Buffer
	jpeg.Encode(&buf, dst, &jpeg.Options{Quality: targetImageQuality})

	return buf.Bytes(), nil
}
