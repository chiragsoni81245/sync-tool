# 🔄 sync-tool

A production-ready Go CLI that automatically syncs local directories to GitHub repos on a scheduled cron pattern.

---

## ✨ Features

- Background cron scheduler for auto-sync
- CLI to add and manage directories
- Git initialization and remote setup
- Modular, production-grade architecture
- SQLite database to track sync metadata
- Structured logging (zap)
- Configurable via YAML

---

## 🛠️ Setup

### 1. Install

```bash
git clone https://github.com/your-org/sync-tool.git
cd sync-tool
go build -o sync-tool
```

### 2. Configure

Create `config.yaml`:

```yaml
cron_schedule: "*/15 * * * *"
database_filepath: "data/data.db"
github_username: "username"
github_token: "ghp_..."
github_email: "email"
```

### 3. Use CLI

```bash
./sync-tool add --provider "github/gdrive" --mode "push/pull" --folder <local-directory> --remote <remote-url>  # Add directory to watch so that it can synced with given provider
./sync-tool list                                                                                                # List current sync targets
./sync-tool start                                                                                               # Start the background sync service
```

### 4. (Optional) Run as Service

```ini
# /etc/systemd/system/sync-tool.service
[Unit]
Description=Sync Tool GitHub Directory Sync Service
After=network.target

[Service]
ExecStart=/usr/local/bin/sync-tool start --config <path-to-config>
Restart=on-failure
WorkingDirectory=/home/ubuntu/sync-tool

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable sync-tool
sudo systemctl start sync-tool
```

---

## 📁 Architecture

```
.
├── cmd             # CLI commands
├── internal
│   ├── config      # Config loader (yaml + env)
│   ├── db          # SQLite models + GORM
│   ├── github      # Github operations
│   ├── logger      # Structured zap logger
│   └── scheduler   # Cron-based sync engine
├── config.yaml     # App config
└── main.go         # Entrypoint (optional)
```

---

## 👨‍💻 Author

Chirag Soni · Built with Go + ❤️
