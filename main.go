package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

// ANSI Escape Sequences for terminal coloring
const (
	ColorReset  = "\033[0m"
	ColorBold   = "\033[1m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorRed    = "\033[31m"
)

// Config holds the application settings parsed from flags and env vars
type Config struct {
	Repo        string
	Profile     string
	Owner       string
	Token       string
	CommandArgs []string
}

// GithubContentResponse mirrors the GitHub API payload structure
type GithubContentResponse struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

func main() {
	cfg, err := parseConfig()
	if err != nil {
		fmt.Printf("%s[ghostenv] %sConfiguration error:%s %v\n", ColorCyan, ColorRed, ColorReset, err)
		os.Exit(1)
	}

	fmt.Printf("%s[ghostenv]%s Fetching remote manifest (%s%s%s) from repo %s%s%s...\n",
		ColorCyan, ColorReset,
		ColorYellow, cfg.Profile, ColorReset,
		ColorBlue, cfg.Repo, ColorReset,
	)

	envVars, err := fetchGitHubEnv(cfg)
	if err != nil {
		fmt.Printf("%s[ghostenv] %sFailed to fetch environment:%s %v\n", ColorCyan, ColorRed, ColorReset, err)
		os.Exit(1)
	}

	if err := runSubprocess(cfg.CommandArgs, envVars); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Printf("%s[ghostenv] %sSubprocess failed:%s %v\n", ColorCyan, ColorRed, ColorReset, err)
		os.Exit(1)
	}
}

// parseConfig handles flags and validates required environment variables
func parseConfig() (*Config, error) {
	var cfg Config

	flag.StringVar(&cfg.Repo, "manifest", "manifest", "manifest repo path")
	flag.StringVar(&cfg.Profile, "profile", "dev", "profile dev/staging/prod etc")
	flag.Parse()

	cfg.CommandArgs = flag.Args()
	if len(cfg.CommandArgs) < 1 {
		return nil, fmt.Errorf("target command is required (e.g., %sghostenv npm run dev%s)", ColorBold, ColorReset)
	}

	cfg.Owner = os.Getenv("GITHUB_USERNAME")
	cfg.Token = os.Getenv("GITHUB_TOKEN")
	if cfg.Owner == "" || cfg.Token == "" {
		return nil, fmt.Errorf("missing credentials; check %sGITHUB_USERNAME%s and %sGITHUB_TOKEN%s variables", ColorYellow, ColorReset, ColorYellow, ColorReset)
	}

	return &cfg, nil
}

// fetchGitHubEnv reaches out to GitHub, decodes the base64 payload, and parses the JSON mapping
func fetchGitHubEnv(cfg *Config) (map[string]string, error) {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/contents/%s/env.json",
		cfg.Owner, cfg.Repo, cfg.Profile,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfg.Token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "Ghostenv")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var contentResp GithubContentResponse
	if err := json.NewDecoder(resp.Body).Decode(&contentResp); err != nil {
		return nil, fmt.Errorf("parsing GitHub JSON response: %w", err)
	}

	if contentResp.Encoding != "base64" {
		return nil, fmt.Errorf("unsupported encoding from GitHub: %s", contentResp.Encoding)
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(contentResp.Content)
	if err != nil {
		return nil, fmt.Errorf("decoding base64 payload: %w", err)
	}

	var envVars map[string]string
	if err := json.Unmarshal(decodedBytes, &envVars); err != nil {
		return nil, fmt.Errorf("parsing decrypted environment JSON: %w", err)
	}

	return envVars, nil
}

// runSubprocess spawns the target command passing down local + github fetched variables
func runSubprocess(args []string, githubEnv map[string]string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Merge system environment variables with fetched ones
	env := os.Environ()
	for key, value := range githubEnv {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}
	cmd.Env = env

	fmt.Printf("%s[ghostenv]%s %s🚀 Executing command:%s %v (%s%d variables injected%s)\n",
		ColorCyan, ColorReset, ColorGreen, ColorReset, args, ColorGreen, len(githubEnv), ColorReset,
	)
	return cmd.Run()
}
