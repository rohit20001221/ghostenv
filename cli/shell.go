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

func CreateCommand(appName *string) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		app := *appName
		if app == "" {
			log.Fatalln("Error: Application name flag (--app) must be provided")
		}

		if len(args) == 0 {
			log.Fatalln("Error: Insufficient arguments. Provide a command to run (e.g., 'ghostenv --app my-app npm run dev')")
		}

		fmt.Printf("[x] Fetching environment for application: %s...\n", app)

		// 1. Build the HTTP Client and Request targeting your standardized API layout
		client := &http.Client{Timeout: 10 * time.Second}
		apiURL := fmt.Sprintf("%s/api/v1/pull/%s", os.Getenv("GHOST_ENV_BASE_URL"), app)

		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			log.Fatalf("Failed to initialize request lifecycle: %v", err)
		}

		// 2. Attach Basic Authentication parameters
		req.SetBasicAuth(os.Getenv("GHOST_ENV_USERNAME"), os.Getenv("GHOST_ENV_PASSWORD"))

		// 3. Dispatch request to backend server
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Failed to reach environment configuration server: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Server rejected credentials or application scope with status: %s", resp.Status)
		}

		// 4. Decode the key-value dictionary from the payload pipeline
		var remoteEnv []map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&remoteEnv); err != nil {
			log.Fatalf("Failed to decode environment payload schema: %v", err)
		}

		// 5. Structure the execution target context
		executionTarget := args[0]
		var executionArgs []string
		if len(args) > 1 {
			executionArgs = args[1:]
		}

		command := exec.Command(executionTarget, executionArgs...)

		command.Stdout = os.Stdout
		command.Stderr = os.Stderr // Added so sub-process runtime failures dump logs explicitly
		command.Stdin = os.Stdin

		// 6. Inherit local system environment strings cleanly
		localEnv := os.Environ()

		// 7. Inject remote system variables formatting them as "KEY=VALUE"
		for _, env := range remoteEnv {
			localEnv = append(localEnv, fmt.Sprintf("%s=%s", env["key"], env["value"]))
		}
		command.Env = localEnv

		// 8. Hand execution flow directly to the sub-process context loop
		fmt.Printf("[x] Spawning runner context: %s\n", executionTarget)
		if err := command.Run(); err != nil {
			// Catch non-zero exit states cleanly without generating a fatalln stack trace crash
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			log.Fatalf("Process system thread dropped unexpectedly: %v", err)
		}
	}
}
