package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type GithubContentResponse struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

func main() {
	var (
		repo    string
		profile string
	)

	flag.StringVar(&repo, "manifest", "manifest", "manifest repo path")
	flag.StringVar(&profile, "profile", "dev", "profile dev/staging/prod etc")
	flag.Parse()

	commandArgs := flag.Args()

	log.Println(commandArgs)
	if len(commandArgs) < 1 {
		log.Fatalln("command required")
	}

	owner := os.Getenv("GITHUB_USERNAME")
	token := os.Getenv("GITHUB_TOKEN")

	if token == "" || owner == "" {
		log.Fatalln("Invalid credentials check the GITHUB_USERNAME and GITHUB_TOKEN environment variables")
	}

	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/contents/%s",
		owner,
		repo,
		fmt.Sprintf("%s/env.json", profile),
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v\n", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "Ghostenv")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error performing request: %v\n", err)
	}

	defer resp.Body.Close()

	// Check for a successful response status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("API returned status %d: %s\n", resp.StatusCode, string(bodyBytes))
		return
	}

	// Read and parse the JSON payload
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	var contentResp GithubContentResponse
	if err := json.Unmarshal(body, &contentResp); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	// GitHub returns file contents as a Base64-encoded string if it's standard text
	if contentResp.Encoding == "base64" {
		decodedBytes, err := base64.StdEncoding.DecodeString(contentResp.Content)
		if err != nil {
			fmt.Printf("Error decoding base64 content: %v\n", err)
			return
		}

		var envVars map[string]string

		if err := json.Unmarshal(decodedBytes, &envVars); err != nil {
			log.Fatalln(err)
		}

		command := exec.Command(commandArgs[0], commandArgs[1:]...)
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		localEnv := os.Environ()
		// --- Track injection stats ---
		varsInjected := 0
		for _, env := range envVars {
			localEnv = append(localEnv, fmt.Sprintf("%s=%s", env, envVars[env]))
			varsInjected++
		}
		command.Env = localEnv

		if err := command.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			log.Fatalf("❌ Process system thread dropped unexpectedly: %v\n", err)
		}
	} else {
		fmt.Printf("Unknown or missing encoding type: %s\n", contentResp.Encoding)
	}
}
