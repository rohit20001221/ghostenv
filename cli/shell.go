package main

import (
	"fmt"
	"log"
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
		if len(args) < 2 {
			log.Fatalln("insufficent args")
		}

		command := exec.Command(args[0], args[1:]...)

		command.Stdout = os.Stdout
		command.Stdin = os.Stdin
		command.Env = os.Environ()

		command.Run()
	}
}
