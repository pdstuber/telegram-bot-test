/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gocarina/gocsv"
	"github.com/pdstuber/telegram-bot-test/internal/bot"
	"github.com/pdstuber/telegram-bot-test/internal/predict/prediction/tensorflow"
	"github.com/spf13/cobra"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const defaultColorChannels = 3

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run the bot",
	Long:  `run the telegram bot and act on messages`,
	Run: func(cmd *cobra.Command, args []string) {
		botAPI, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
		if err != nil {
			log.Panic(err)
		}

		botAPI.Debug = true

		log.Printf("Authorized on account %s", botAPI.Self.UserName)
		var labels []tensorflow.Label

		modelPath := os.Getenv("MODEL_PATH")
		model, err := os.ReadFile(fmt.Sprintf("%s/model.pb", modelPath))
		if err != nil {
			log.Panic(err)
		}
		labelBytes, err := os.ReadFile(fmt.Sprintf("%s/labels.csv", modelPath))
		if err != nil {
			log.Panic(err)
		}

		if err := gocsv.UnmarshalBytes(labelBytes, &labels); err != nil {
			log.Fatalf("could not unmarshal labels csv: %v\n", err)
		}

		bot := bot.New(botAPI, tensorflow.New(model, labels, defaultColorChannels))

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		bot.Start(ctx)

		log.Println("after start")

		<-ctx.Done()
		bot.Stop()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
