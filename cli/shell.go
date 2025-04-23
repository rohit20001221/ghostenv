package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func CreateCommand(appName *string) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		app := *appName
		fmt.Println("[x] application:", app)

		/*
			fetch the environment variables based on the appName via some api populate them into the command
		*/
		command := exec.Command(args[0], args[1:]...)

		command.Stdout = os.Stdout
		command.Stdin = os.Stdin
		command.Env = []string{}

		command.Run()
	}
}
