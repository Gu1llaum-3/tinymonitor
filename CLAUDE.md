# TinyMonitor - Go

## Description
TinyMonitor est un agent de monitoring système léger avec alertes multi-canaux.
- **Fonction** : Surveillance CPU, RAM, disque, load, I/O avec alertes
- **Architecture** : Boucle de polling + collectors metrics + providers d'alertes
- **Cible** : Binaire standalone pour Linux/macOS (AMD64 + ARM64)

## Architecture
```
tinymonitor/
├── cmd/tinymonitor/main.go     # Point d'entrée
├── internal/
│   ├── config/config.go        # Structs typées + chargement
│   ├── monitor/monitor.go      # Boucle + state machine
│   ├── alerts/
│   │   ├── manager.go          # Worker pool avec goroutines
│   │   ├── provider.go         # Interface Alert
│   │   ├── google_chat.go, ntfy.go, smtp.go, webhook.go, gotify.go
│   ├── metrics/
│   │   ├── collector.go        # Interface Metric
│   │   ├── cpu.go, memory.go, disk.go, load.go, io.go, reboot.go
│   ├── models/                 # Types partagés (Alert, MetricResult, etc.)
│   └── utils/                  # Helpers (system_info)
├── configs/config.example.json
├── go.mod
└── go.sum
```

## Conventions Go
- Gestion d'erreurs explicite (jamais de panic en production)
- Context pour timeouts et cancellation
- Interfaces définies par le consommateur
- Logging avec `log/slog` (stdlib Go 1.21+)
- Tests dans fichiers `*_test.go` avec table-driven tests

## Dépendances
- `github.com/shirou/gopsutil/v3` - Métriques système
- Stdlib uniquement pour le reste (net/http, encoding/json, context, sync)

## Commandes utiles
```bash
# Build
go build -o tinymonitor ./cmd/tinymonitor

# Tests
go test ./...

# Vet
go vet ./...

# Build multi-arch
GOOS=linux GOARCH=amd64 go build -o tinymonitor-linux-amd64 ./cmd/tinymonitor
GOOS=linux GOARCH=arm64 go build -o tinymonitor-linux-arm64 ./cmd/tinymonitor
GOOS=darwin GOARCH=arm64 go build -o tinymonitor-darwin-arm64 ./cmd/tinymonitor
```

## Usage
```bash
# Avec config par défaut (recherche config.json, ~/.config/tinymonitor/config.json, /etc/tinymonitor/config.json)
./tinymonitor

# Avec config spécifique
./tinymonitor -c /path/to/config.json

# Version
./tinymonitor -v
```
