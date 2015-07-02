// Package tgbot is a framework for building Telegram bots with the Bot API.
package tgbot

import (
	"fmt"
	"log"
	"strings"

	"github.com/Syfaro/telegram-bot-api"
)

// Config holds all plugin settings
type Config map[string]string

//Bot is all the methods you need to operate the bot.
type Bot struct {
	API            *tgbotapi.BotAPI
	Plugins        []Plugin
	Config         Config
	BeforeCommands func(*tgbotapi.BotAPI, *tgbotapi.Update)
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
	var bot Bot
	ok := false
	for !ok {
		fmt.Println("What is your bot token?")
		var token string
		fmt.Scanln(&token)

		bot, err := NewBot(token)
		if err != nil {
			log.Println("Your token is invalid!\nRetry.")
			continue
		}

		log.Printf("Are you sure you wish to use the account %s? [Y/n] ", bot.API.Self.UserName)
		var answer string
		fmt.Scanln(&answer)
		if strings.ToLower(answer) == "n" {
			continue
		}
		ok = true
	}
	return bot
}

// This starts your plugin and run their tasks in a goroutine
// So keep attention in what you are doing
func (bot *Bot) Start() {
	for _, plugin := range bot.Plugins {
		plugin.Setup(bot.API.Debug)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.API.UpdatesChan(u)

	for update := range updates {
		for _, plugin := range bot.Plugins {
			go bot.BeforeCommands(bot.API, &update)
			cmd, err := plugin.Commands.Parse(update.Message.Text)
			if err != nil {
				log.Printf(err.Error())
				continue
			}
			go cmd.RunTask(&update, bot.API)
		}
	}
}
