package tgbot

import (
	"github.com/syfaro/telegram-bot-api"
)

// Plugin is the interface for all basic methods a plugin must have.
type Plugin interface {
	// MUST BE UNIQUE name of plugin
	GetName() string
	// Slice of available commands, WITH forward slash (if needed)
	GetCommands() []string
	// Slice of help text messages, sent below the plugin name in help command
	GetHelpText() []string
	// Function to set up needed configuration, etc.
	Setup(*Bot) bool
	// Called when Message is received, params are Message, command, args
	GotCommand(tgbotapi.Message, string, []string)
}

func (bot *Bot) AddPlugin(plugin Plugin) {
	bot.Plugins = append(bot.Plugins, plugin)
}
