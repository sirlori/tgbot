# tgbot
[![GoDoc](https://godoc.org/github.com/Syfaro/tgbot?status.svg)](http://godoc.org/github.com/Syfaro/tgbot)

A Golang Telegram bot framework

Expect many breaking changes. Open an issue if you have any feature requests. 

If you want just API bindings, check out [telegram-bot-api](https://github.com/Syfaro/telegram-bot-api).

###Examples

```go
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetOutput(os.Stdout)

	// Create a new bot
	bot, err := tgbot.NewBot(TOKEN)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.API.Self.UserName)

	// Add commands to your plugin
	plugin := tgbot.NewPlugin("BOT")
	// Aommand setup in your plugin	
	plugin.AddSetup(InitHelp)
	plugin.AddSetup(InitEcho)
	
	// Add plugin to the bot:
	// In future this will be handled by our package
	// so it will be possible to add plugins from command-line
	// from our repo.
	bot.AddPlugin(plugin)
	
	// A function that runs fore every update without conditions
	bot.BeforeCommands = SomeStuff
	// Start to recieve commands
	bot.Start()
}

```

This is the init function for the echo command:

```go
func InitEcho(plugincommands *tgbot.MultipleParser) error {

	parser := tgbot.NewParser()
	parser.Cmd = "^/echo"
	parser.Profile = []tgbot.Argouments{
		tgbot.NewArgs("printer", []string{"string"}),
	}

	parser.Tasks["printer"] = func(cmd *tgbot.Command, up *tgbotapi.Update,
		bot *tgbotapi.BotAPI) {
		msg := tgbotapi.NewMessage(up.Message.Chat.ID, cmd.Args["string"])
		bot.SendMessage(msg)
	}
	plugincommands.Add(parser)
	return nil
}
```

The docs are still poor, but try to read them to understand more!
