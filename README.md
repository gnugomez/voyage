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
voyage deploy -r <repo-url> -b <branch> -c <compose-path> -o <out-path> [-f] [-l debug|info|error|fatal]
```

Alternatively, you can use a JSON configuration file:

```sh
voyage deploy -config /path/to/config.json
```

| Flag      | Description                                                    |
| --------- | -------------------------------------------------------------- |
| `-r`      | Git repository URL                                             |
| `-b`      | Branch name                                                    |
| `-c`      | Path to `docker-compose.yml` (can be specified multiple times) |
| `-o`      | Output directory for the repo                                  |
| `-f`      | Force deployment (optional)                                    |
| `-l`      | Log level (default: info)                                      |
| `-config` | Path to a JSON configuration file (optional)                   |

### Configuration File

As an alternative to providing all arguments on the command line, you can use a JSON configuration file by specifying the `-config` flag.

> [!NOTE]
> When the `-config` flag is used, all other command-line arguments for the deploy command are ignored.

**Example `config.json`:**

```json
{
  "repo": "https://github.com/user/repo.git",
  "branch": "main",
  "outPath": "~/deployments/repo",
  "remoteComposePaths": [
    "docker/app1/compose.yml",
    "docker/app2/compose.yml"
  ],
  "force": false,
}
```

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
