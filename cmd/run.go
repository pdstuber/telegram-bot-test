/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const token = "6853139137:AAHt9ab7_9dmcC-zq0cF4CbCrAaSxQKTPu8"

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run the bot",
	Long:  `run the telegram bot and act on messages`,
	Run: func(cmd *cobra.Command, args []string) {
		bot, err := tgbotapi.NewBotAPI(token)
		if err != nil {
			log.Panic(err)
		}

		bot.Debug = true

		log.Printf("Authorized on account %s", bot.Self.UserName)

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates := bot.GetUpdatesChan(u)

		for update := range updates {
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

			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
		}
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
