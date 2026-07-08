package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

// ANSI Escape Color Codes
const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Ghost  = "\033[38;5;141m" // Soft purple/magenta glow
	Cyan   = "\033[36m"
	Green  = "\033[32m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
)

func CreateCommand(appName *string) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		app := *appName
		if app == "" {
			log.Fatalln(Red + "❌ Error: Application name flag (--app) must be provided" + Reset)
		}

		if len(args) == 0 {
			log.Fatalln(Red + "❌ Error: Provide a command to execute (e.g., 'ghostenv --app billing npm run dev')" + Reset)
		}

		// --- Fun, colorful intro ---
		fmt.Printf("%s👻 ghostenv%s %s⟡ Summoning environment for:%s %s%s%s...\n",
			Ghost, Reset, Cyan, Reset, Bold+Green, app, Reset)

		client := &http.Client{Timeout: 10 * time.Second}
		apiURL := fmt.Sprintf("%s/api/v1/pull/%s", os.Getenv("GHOST_ENV_BASE_URL"), app)

		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			log.Fatalf(Red+"❌ Failed to initialize request lifecycle: %v\n"+Reset, err)
		}

		req.SetBasicAuth(os.Getenv("GHOST_ENV_USERNAME"), os.Getenv("GHOST_ENV_PASSWORD"))

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf(Red+"❌ Failed to reach environment configuration server: %v\n"+Reset, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf(Red+"❌ Server rejected credentials or application scope with status: %s\n"+Reset, resp.Status)
		}

		var remoteEnv []map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&remoteEnv); err != nil {
			log.Fatalf(Red+"❌ Failed to decode environment payload schema: %v\n"+Reset, err)
		}

		executionTarget := args[0]
		var executionArgs []string
		if len(args) > 1 {
			executionArgs = args[1:]
		}

		command := exec.Command(executionTarget, executionArgs...)

		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		command.Stdin = os.Stdin

		localEnv := os.Environ()

		// --- Track injection stats ---
		varsInjected := 0
		for _, env := range remoteEnv {
			localEnv = append(localEnv, fmt.Sprintf("%s=%s", env["key"], env["value"]))
			varsInjected++
		}
		command.Env = localEnv

		// --- Success message before spawning ---
		fmt.Printf("%s✨ Success!%s %sInjected %d ghostly variables into context.%s\n",
			Green, Reset, Cyan, varsInjected, Reset)
		fmt.Printf("%s🚀 Spawning process:%s %s%s%s\n\n",
			Yellow, Reset, Bold, executionTarget, Reset)

		if err := command.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			log.Fatalf(Red+"❌ Process system thread dropped unexpectedly: %v\n"+Reset, err)
		}
	}
}
