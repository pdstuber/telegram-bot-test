/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	bot "github.com/pdstuber/telegram-bot-test/internal"
	"github.com/spf13/cobra"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates := botAPI.GetUpdatesChan(u)

		bot := bot.New(botAPI)

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		go bot.HandleUpdates(ctx, updates)
		<-ctx.Done()
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
