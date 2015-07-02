package tgbot

import (
	"errors"
	"log"
)

type Plugin struct {
	Name     string
	Commands MultipleParser
	setups   []func(*MultipleParser) error
}

//Setup your plugin
func (plugin *Plugin) Setup(debug bool) error {
	var err error
	for _, f := range plugin.setups {
		err = f(&plugin.Commands)
		if debug {
			log.Printf(err.Error())
		}
	}
	if err != nil {
		return errors.New("The commands setups returned errors")
	}
	return nil
}

//Add setup that let's you create a command for your plugin
func (plugin *Plugin) AddSetup(setup func(*MultipleParser) error) {
	plugin.setups = append(plugin.setups, setup)
}

//Add a plugin to your bot
func (bot *Bot) AddPlugin(plugin Plugin) {
	bot.Plugins = append(bot.Plugins, plugin)
}

func NewPlugin(name string) Plugin {
	plugin := Plugin{}
	plugin.Name = name
	plugin.Commands = NewMultipleParser()
	return plugin
}
