package main

import (
	"bytes"
	"fmt"
	"strings"
)

func initCommandsHelp() {
	listBuffer := bytes.Buffer{}
	helpBuffer := bytes.Buffer{}

	commands := map[string]commandReference{
		"list": {
			aliases: commandAliases{"list", "ls"},
			help: `	List files in the current directory`,
		},
		"chdir": {
			aliases: commandAliases{"chdir", "cd"},
			help: `	DIRECTORY		Change current directory`,
		},
		"get": {
			aliases: commandAliases{"get", "download", "fetch"},
			help: `FILENAME [DIRECTORY]	Download file to the local directory (default is for -D option at the startup)`,
		},
		"help": {
			aliases: commandAliases{"help", "h"},
			help: `	Print commands usage help`,
		},
		"exit": {
			aliases: commandAliases{"exit", "quit", "q"},
			help:    `Quit the program`,
		},
	}

	for _, reference := range commands {
		aliasesText := strings.Join(reference.aliases, ", ")
		listBuffer.WriteString(fmt.Sprintf(" * %s\n", aliasesText))
		helpBuffer.WriteString(fmt.Sprintf(" * %s\t%s\n", aliasesText, reference.help))
	}

	commandsList = listBuffer.String()
	commandsHelp = helpBuffer.String()
}

func printCommandsHelp() {
	fmt.Printf(`Usage:

%s
`, commandsHelp)
}

func printCommandsList() {
	fmt.Printf(`Available commands:

%s
`, commandsList)
}
