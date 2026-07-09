
# ghostenv 👻

A zero-dependency, secure CLI tool built in Go to hot-inject environment variable maps fetched directly from your private GitHub configurations into downstream child processes.

No more hardcoded `.env` files floating around local machines or leaking out into your git history.

---

## 🛠️ How It Works

`ghostenv` acts as a lightweight proxy wrapper. When executed, it pulls down a structured JSON config file from a repository via the GitHub REST API, securely base64-decodes it in memory, merges those variables seamlessly with your existing system shell variables, and boots up your downstream application.

---

## 🔒 Setup & Authentication

The utility relies on system environment variables to authenticate against private repositories securely. Set these up in your active terminal session or deployment profile:

```bash
# Define your personal or organization GitHub username
export GITHUB_USERNAME="your_github_username_or_org"

# Provide a Classic or Fine-Grained Personal Access Token (PAT) with repository content read scopes
export GITHUB_TOKEN="ghp_yourSecureGitHubTokenGoesHere"

```

---

## 📂 Expected Repository Structure

`ghostenv` dynamically maps paths according to the `-profile` value requested. By default, it looks for an `env.json` configuration nested within folders corresponding to your target environment stage:

```text
your-manifest-repository/
├── dev/
│   └── env.json
├── staging/
│   └── env.json
└── prod/
    └── env.json

```

### `env.json` Layout Format

Ensure your configuration uses standard, flat key-value strings:

```json
{
  "DATABASE_URL": "postgresql://localhost:5432/mydb",
  "API_SECRET": "super-secret-cryptographic-string",
  "DEBUG_MODE": "true"
}

```

---

## 🚀 Usage

Execute `ghostenv` by supplying your destination settings as flags, immediately followed by the specific runtime or binary command you wish to wrap.

### Configuration Flags

* **`-manifest`** *(string, default: "manifest")*: The name of your target configuration repository.
* **`-profile`** *(string, default: "dev")*: The directory scope or deployment phase to fetch.

### Practical Code Execution

```bash
# Example 1: Launch a Go microservice inside your local development sandbox
ghostenv -manifest="app-configs" -profile="dev" go run main.go

# Example 2: Wrap an npm deployment engine running in a staging cluster
ghostenv -manifest="infra-secrets" -profile="staging" npm run start

```

---

## 🎨 Console Output & Diagnostics

`ghostenv` provides streamlined, timestamp-free terminal diagnostics wrapped in ANSI fallback escape colors so you can inspect active environment injection pipelines at a glance:

```text
[ghostenv] Fetching remote manifest (staging) from repo app-configs...
[ghostenv] 🚀 Executing command: [go run main.go] (14 variables injected)
