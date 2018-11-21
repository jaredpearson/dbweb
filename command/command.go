/*
Package command implements command line parser

This is a custom command line parser. This app is not using the
golang's flag because we need to support nested commands (styled
like Docker or Git).
*/
package command

import (
	"fmt"
	"os"
)

type Arg struct {
	name        string
	description string
	Value       string
}

type Command struct {
	name          string
	description   string
	selected      bool
	subCommandSet *CommandSet
	args          []*Arg
}

func (command *Command) AddArg(name, description string) *Arg {
	if command.args == nil {
		command.args = []*Arg{}
	}
	newArg := &Arg{
		name:        name,
		description: description,
	}
	command.args = append(command.args, newArg)
	return newArg
}
func (command *Command) GetArg(index int) (*Arg, bool) {
	if command.args == nil {
		return nil, false
	}
	if index < 0 || index >= len(command.args) {
		return nil, false
	}
	arg := command.args[index]
	return arg, true
}
func (command *Command) AddSubcommand(name, description string) *Command {
	if command.subCommandSet == nil {
		command.subCommandSet = NewCommandSet()
	}
	return command.subCommandSet.AddCommand(name, description)
}
func (command *Command) DisplayUsage() {
	if command.subCommandSet == nil {
		panic(fmt.Sprintf("Cannot displayUsage for command without subcommands: %s", command.name))
	}
	command.subCommandSet.DisplayUsage()
}
func (command *Command) IsSelected() bool {
	return command.selected
}

type CommandSet struct {
	commandsByName map[string]*Command
}

func NewCommandSet() *CommandSet {
	r := &CommandSet{}
	return r
}

func (commands *CommandSet) AddCommand(name, description string) *Command {
	cmd := &Command{
		name:        name,
		description: description,
	}
	if commands.commandsByName == nil {
		commands.commandsByName = make(map[string]*Command)
	}
	commands.commandsByName[name] = cmd
	return cmd
}

func (commands *CommandSet) VisitAll(fn func(*Command)) {
	for _, cmd := range commands.commandsByName {
		fn(cmd)
	}
}

func (commands *CommandSet) Parse() {
	args := os.Args[1:]

	var selected string
	if len(args) == 0 {
		selected = "help"
	} else {
		selected = args[0]
		args = args[1:] // remove the first and continue parsing
	}

	cmdInst, exists := commands.commandsByName[selected]
	if !exists {
		fmt.Fprintf(os.Stderr, "Unknown command specified: %s\n", selected)
		os.Exit(1)
	}
	cmdInst.selected = true

	// if there are subcommands, then the next value is a subcommand
	if len(args) > 0 && cmdInst.subCommandSet != nil {
		subCommandName := args[0]
		args = args[1:] // remove the first and continue parsing
		subInst, exists := cmdInst.subCommandSet.commandsByName[subCommandName]
		if !exists {
			fmt.Fprintf(os.Stderr, "Unknown subcommand specified for %s: %s", cmdInst.name, subCommandName)
		}
		subInst.selected = true
		cmdInst = subInst
	}

	// set any args. ignore if none have been setup
	if cmdInst.args != nil {
		for i, value := range args {
			if i >= len(cmdInst.args) {
				break
			}
			cmdInst.args[i].Value = value
		}
	}
}
func (commands *CommandSet) DisplayUsage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	commands.VisitAll(func(a *Command) {
		fmt.Fprintf(os.Stderr, "\t%s\t%s\n", a.name, a.description)
	})
}
