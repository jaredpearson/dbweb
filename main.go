package main

import (
	"fmt"
	"os"
	"strings"

	"jaredpearson.com/dbweb/command"
	"jaredpearson.com/dbweb/data"
	"jaredpearson.com/dbweb/web"
)

func executeUserCommand(userCmd *command.Command, usersAddCmd *command.Command) {
	if usersAddCmd.IsSelected() {
		usernameArg, _ := usersAddCmd.GetArg(0)
		username := strings.Trim(usernameArg.Value, " ")
		if username == "" {
			fmt.Fprint(os.Stderr, "Username is required when adding a new user\n")
			os.Exit(1)
		}

		_, err := data.GetUserByUsername(username)
		// we expect a ErrUserNotFound
		if err == nil {
			fmt.Fprintf(os.Stderr, "User already exists with username %s\n", username)
			os.Exit(1)
		} else if ae, ok := err.(*data.ErrUserNotFound); !ok {
			fmt.Fprintf(os.Stderr, "Unable to add user: %s\n%v\n", username, ae)
			os.Exit(1)
		}

		err = data.AddUser(username)
		if err != nil {
			fmt.Printf("Error adding user: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "Added new user %s\n", username)
		os.Exit(0)
	} else {
		userCmd.DisplayUsage()
		os.Exit(1)
	}
}

func main() {
	var commands = command.NewCommandSet()
	helpCmd := commands.AddCommand("help", "Displays the help information")
	startCmd := commands.AddCommand("start", "Starts the web server")
	userCmd := commands.AddCommand("users", "Manage users")
	usersAddCmd := userCmd.AddSubcommand("add", "Adds a new user")
	usersAddCmd.AddArg("username", "The username of the new user")

	commands.Parse()

	if helpCmd.IsSelected() {
		commands.DisplayUsage()
		os.Exit(0)
	} else if startCmd.IsSelected() {
		web.ServerStart()
	} else if userCmd.IsSelected() {
		executeUserCommand(userCmd, usersAddCmd)
	} else {
		fmt.Fprint(os.Stderr, "Invalid or unknown command specified\n")
		commands.DisplayUsage()
		os.Exit(1)
	}
}
