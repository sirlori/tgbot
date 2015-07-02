package tgbot

import (
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/Syfaro/telegram-bot-api"
)

// This structs use the information you give it of a command
// to parse a command
type Argouments struct {
	Name string
	Args []string
}

// Struct to parse command and automatically execute
// Tasks with them once they are recognized
// by the parser
type Command struct {
	Cmd  string
	Args map[string]string
	Task func(*Command, *tgbotapi.Update,
		*tgbotapi.BotAPI)
	Match    string
	FullText string
}

// This is parser that you want to declare
// to add your command definition to the bot
type CommandParser struct {
	Cmd     string
	Profile []Argouments
	Tasks   map[string]func(*Command, *tgbotapi.Update,
		*tgbotapi.BotAPI)
	Help string
	Sep  string
}

// This groups parsers and choose which
// command to use for every update
type MultipleParser struct {
	Botname string
	parsers []CommandParser
}

// Call it to create a MultiParser, that store and
// uses your parsers to call their tasks when needed
func NewMultipleParser(name string) MultipleParser {
	return MultipleParser{
		name,
		make([]CommandParser, 0, 5),
	}
}

// Add a parser to your multiparser
func (p *MultipleParser) Add(cp CommandParser) {
	p.parsers = append(p.parsers, cp)
}

// Parses with all your added parsers and returns
// the proper command.
// It is limited to only one command for performance
// reasons
func (p *MultipleParser) Parse(text string) (Command, error) {
	for _, v := range p.parsers {
		cmd, err := v.Parse(text)
		if err == nil {
			return cmd, err
		}
	}
	return Command{}, errors.New("No matches for these parsers")
}

// Parses only for a single command.
func (parser *CommandParser) Parse(text string) (Command, error) {
	text = strings.TrimSpace(text)
	re, err := regexp.Compile(parser.Cmd)
	if err != nil {
		log.Printf(err.Error())
		return Command{}, err
	}

	matches := re.FindAllString(text, 1)
	if len(matches) != 1 || len(matches) == 1 && matches[0] == "" {
		return Command{}, errors.New("No matches in this parser")
	}

	cmdlen := len(matches[0])
	argsval := make([]string, 1)

	if parser.Sep != "" {
		argsval = ClearSlice(strings.Split(text[cmdlen:], parser.Sep))
	} else if !strings.ContainsRune(text, ' ') {
		argsval = []string{}
	} else {
		argsval[0] = text[cmdlen:]
	}
	args := make(map[string]string)
	argname := ""

	for _, argouments := range parser.Profile {
		if len(argouments.Args) == len(argsval) {
			argname = argouments.Name
			for i := 0; i < len(argsval); i++ {
				args[argouments.Args[i]] = argsval[i]
			}
		}
	}

	task := parser.Tasks[argname]
	command := Command{}
	command.Cmd = parser.Cmd
	command.Args = args
	if task == nil {
		return command, errors.New("No task found for this command: " +
			argname)
	}
	command.Task = task
	return command, nil
}

// Add argouments to your command. If zero argouments , meaning: []string{}
// the command should be only his name, otherwise other text is needed
// taskname is for binding a command syntax to a task
func NewArgs(taskname string, args []string) Argouments {
	return Argouments{
		taskname,
		args,
	}
}

// Create a parser for a single command
func NewParser() CommandParser {
	parser := CommandParser{}
	parser.Tasks = make(map[string]func(*Command, *tgbotapi.Update,
		*tgbotapi.BotAPI))
	return parser
}

// Run the task of the resulting command
// that has probably been parsed by a MultipleParser
func (cmd *Command) RunTask(update *tgbotapi.Update,
	bot *tgbotapi.BotAPI) {
	cmd.Task(cmd, update, bot)
}
