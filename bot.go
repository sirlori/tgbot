// Package tgbot is a framework for building Telegram bots with the Bot API.
package tgbot

import (
	"fmt"
	"github.com/syfaro/telegram-bot-api"
	"log"
	"strings"
)

// Config holds all plugin settings
type Config map[string]string

// Bot is all the methods you need to operate the bot.
type Bot struct {
	API     *tgbotapi.BotAPI
	Plugins []Plugin
	Config  Config
}

// NewBot creates a new Bot instance with a token
func NewBot(token string) (Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return Bot{}, err
	}

	return Bot{
		API: bot,
	}, nil
}

func InitBot() Bot {
	fmt.Println("What is your bot token?")
	var token string
	fmt.Scanln(&token)

	bot, err := NewBot(token)
	if err != nil {
		log.Println("Your token is invalid!")
		return InitBot()
	}

	log.Printf("Are you sure you wish to use the account %s? [Y/n] ", bot.API.Self.UserName)
	var answer string
	fmt.Scanln(&answer)

	if strings.ToLower(answer) == "n" {
		return InitBot()
	}

	return bot
}

func (bot *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.API.UpdatesChan(u)

	for update := range updates {
		args := strings.Split(update.Message.Text, " ")

		for _, plugin := range bot.Plugins {
			for _, command := range plugin.GetCommands() {
				if args[0] == command {
					plugin.GotCommand(update.Message, args[0], args[1:])
				}
			}
		}
	}
}
