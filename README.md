# Ghostenv 👻

A minimalistic, secure environment variable manager written in Golang. 

**Ghostenv** eliminates the friction and security risks of dragging raw `.env` files around local workspaces. Instead, it lets you store your application parameters centrally through a clean self-hosted web dashboard and securely injects them into any shell command at runtime via a lightweight compiled CLI wrapper.

---

## 🛠️ Local Setup & Installation

Follow these steps to clone, compile, and boot the entire **Ghostenv** ecosystem (both the backend dashboard server and the companion CLI tool) locally from source.

### Prerequisites
* **Go** (`1.22+` required to leverage native `r.PathValue` multiplexer routing patterns)

---

### Step 1: Clone and Run the Server Engine

1. **Clone the repository:**
   ```bash
   git clone https://github.com/rohit20001221/ghostenv-server.git
   cd ghostenv-server
   ```

2. **Download Go modules:**
   ```bash
   go mod download
   ```

3. **Boot the backend server:**
   ```bash
   go run main.go
   ```
   *The application will boot up. Open your web browser and navigate to `http://localhost:8080` to access the registration, login, and application workspace views.*

---

### Step 2: Build and Install the CLI Client

Open a new terminal window inside the project directory to compile your standalone client binary:

1. **Compile the CLI executable:**
   ```bash
   go build -o ghostenv cmd/cli/main.go
   ```
   *(Note: Adjust the source path `cmd/cli/main.go` to match where your CLI entry point file resides in your project structure).*

2. **Make it executable:**
   ```bash
   chmod +x ghostenv
   ```

3. **Install it globally to your system PATH:**
   ```bash
   sudo mv ghostenv /usr/local/bin/
   ```

4. **Verify the installation:**
   ```bash
   ghostenv --help
   ```

---

## 🚀 How to Use It

Once your server is running and the CLI is globally compiled on your system, managing your configurations takes 3 steps:

### 1. Create your Application Space
Head over to the Web UI at `http://localhost:8080`, register an account, and initialize a new application workspace target (e.g., `payment-service`).

### 2. Populate Keys & Values
Inside your application view, use the form input fields to add the environment variables your code relies on (e.g., `DATABASE_URL`, `API_SECRET`).

### 3. Hot-Inject into Code Runtimes
Instead of running your code raw or utilizing a local `.env` flat file, wrap your typical shell script using `ghostenv`. It will reach out to your local running dashboard server API, stream the values safely straight into process memory, and kick off your script:

```bash
# Syntax Architecture Blueprint
ghostenv --app <app_name> "<shell command>"

# Real-world Running Examples
ghostenv --app payment-service "go run main.go"
ghostenv --app frontend-dashboard "npm run start"
```

---

## 📦 Project Directory Layout

```text
├── controllers/          # Server routing logic handlers (Auth, Apps, Variables)
├── middlewares/          # Access Control & Session Verification pipelines
├── templates/            # Web Canvas Blocks (HTML + Tailwind CSS components)
│   ├── base.html.tmpl    # Global UI HTML shell frame layout 
│   ├── home.html.tmpl    # Main application dashboard and drop utilities
│   └── ...               # Login, Register, Product Intro views
├── types/                # Strict typed internal structs & context keys
├── cmd/
│   └── cli/              # Source code directory for the ghostenv CLI tool
└── main.go               # Hub web server bootstrap initialization entry point
```

---

## 🔒 Security Posture
* **No Disk Pollution:** Variables stay completely ephemeral. They live inside process memory strings only while the wrapped command executes—leaving zero `.env` residual artifacts on hard drives or accidental repository pushes.
* **Session Integrity:** API endpoints and dashboard UI states are protected behind active middleware filters that verify cookie token state handshakes prior to exposing parameter values.

