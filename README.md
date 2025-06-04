# ⛵ Voyage

A simple CLI tool to automate deployment of Docker Compose projects from a Git repository.

I started this project to automate the deployment of my own homelab infrastructure. I don't have a very powerful server, and all the tools I found capable of what I wanted were either too complex or too heavy for my needs. So, I decided to build my own tool.

This project came out as a way for me to learn Go, so please don't expect it to be perfect.

## ✨ Features

- 🐳 Pulls a Git repository and checks for changes in a subdirectory
- 🔄 Runs `docker compose up` only if changes are detected (or with `-f`)
- 🛠️ Supports custom compose file paths, branches, and output directories

## ⚡ Usage


```sh
voyage deploy -r <repo-url> -b <branch> -c <compose-path> -o <out-path> [-f] [-d] [-l debug|info|error|fatal]
```

| Flag | Description                                                    |
| ---- | -------------------------------------------------------------- |
| `-r` | Git repository URL                                             |
| `-b` | Branch name                                                    |
| `-c` | Path to `docker-compose.yml` (can be specified multiple times) |
| `-o` | Output directory for the repo                                  |
| `-f` | Force deployment (optional)                                    |
| `-l` | Log level (default: info)                                      |

> [!IMPORTANT]  
> Since this tool detects what needs to be deployed by checking the remote repository for changes, you may want to run it as 
> a single instance if you're watching multiple compose files. Otherwise, it wouldn't be able to detect changes properly.

### Example

```sh
voyage deploy -r https://github.com/user/repo.git -b main -o ~/deployments/repo \
  -c docker/app1/compose.yml \
  -c docker/app1/compose.override.yml \
  -c docker/app2/compose.yml \
  -c frontend/compose.yml
```

## 📦 Requirements

- Docker & Docker Compose
- Go 1.24+

---

Made with ❤️ by gnugomez
