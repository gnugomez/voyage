# ⛵ Voyage

A simple CLI tool to automate deployment of Docker Compose projects from a Git repository.

## ✨ Features
- 🐳 Pulls a Git repository and checks for changes in a subdirectory
- 🔄 Runs `docker compose up` only if changes are detected (or with `--force`)
- 🛠️ Supports custom compose file paths, branches, and output directories

## ⚡ Usage
```sh
voyage -r <repo-url> -b <branch> -c <compose-path> -o <out-path> [--force] [--daemon] [--log-level debug|info|error|fatal]
```

| Flag         | Description                          |
|--------------|--------------------------------------|
| `-r`         | Git repository URL                   |
| `-b`         | Branch name                          |
| `-c`         | Path to `docker-compose.yml`         |
| `-o`         | Output directory for the repo        |
| `-f`         | Force deployment (optional)          |
| `-d`         | Run in daemon mode (default: true)   |
| `-l`         | Log level (default: info)            |

### Example
```sh
voyage -r https://github.com/user/repo.git -b main -c compose.yml -o ~/local-deployments/repo-folder -l debug
```

## 📦 Requirements
- Docker & Docker Compose
- Go 1.24+

---
Made with ❤️ by gnugomez
