# TinyMonitor - Migration Python → Go

## Contexte du projet
TinyMonitor est un agent de monitoring système léger (~1,053 LOC Python) avec alertes multi-canaux.
- **Fonction** : Surveillance CPU, RAM, disque, load, I/O avec alertes
- **Architecture** : Boucle de polling + plugins metrics + providers d'alertes
- **Cible** : Binaire standalone pour Linux/macOS (AMD64 + ARM64)

## Architecture actuelle Python
```
tinymonitor/
├── monitor.py          # Boucle principale, machine d'état, rate limiting
├── config.py           # Chargement config JSON avec cascade
├── metrics/
│   ├── cpu.py, memory.py, disk.py, load.py, io.py, reboot.py
├── alerts/
│   ├── manager.py      # Dispatch async via ThreadPoolExecutor
│   ├── google_chat.py, ntfy.py, smtp.py, webhook.py, gotify.py
```

## Architecture Go cible
```
tinymonitor-go/
├── cmd/tinymonitor/main.go
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
│   └── models/                 # Types partagés (Alert, MetricResult, etc.)
├── go.mod
├── go.sum
└── configs/config.example.json
```

## Points critiques de migration

### 1. ThreadPoolExecutor → Worker Pool (CRITIQUE)
**Python :**
```python
self.executor = ThreadPoolExecutor(max_workers=5)
self.executor.submit(provider.send, alert)
```
**Go attendu :**
- Channel buffered pour les alertes
- Pool de goroutines workers
- Graceful shutdown avec context.Context
- sync.WaitGroup pour attendre la fin

### 2. Configuration dynamique → Structs typées (CRITIQUE)
**Python :** dict merge récursif
**Go attendu :**
- Structs avec tags JSON
- Valeurs par défaut via constructeur
- Validation explicite
- Pas de `interface{}` sauf nécessité absolue

### 3. psutil → gopsutil
Mapping direct des fonctions, API similaire.

## Conventions Go à respecter
- Gestion d'erreurs explicite (jamais de panic en production)
- Context pour timeouts et cancellation
- Interfaces définies par le consommateur (pas le provider)
- Logging avec `log/slog` (stdlib Go 1.21+)
- Tests dans fichiers `*_test.go` avec table-driven tests

## Dépendances Go autorisées
- `github.com/shirou/gopsutil/v3` - Métriques système
- `gopkg.in/gomail.v2` - Emails SMTP (optionnel, net/smtp suffit)
- Stdlib uniquement pour le reste (net/http, encoding/json, context, sync)

## Ordre de migration recommandé
1. `internal/models/` - Types de base (Alert, Severity, MetricResult)
2. `internal/config/` - Chargement configuration
3. `internal/metrics/` - Collectors (commencer par cpu.go)
4. `internal/alerts/provider.go` - Interface + console provider
5. `internal/alerts/manager.go` - Worker pool
6. `internal/monitor/` - Boucle principale
7. `cmd/tinymonitor/main.go` - Point d'entrée
8. Autres providers d'alertes
9. Tests

## Format de configuration JSON (à conserver)
La config JSON doit rester compatible avec le format Python existant pour faciliter la migration des utilisateurs.

## Commandes utiles
```bash
# Build
go build -o tinymonitor ./cmd/tinymonitor

# Tests
go test ./...

# Build multi-arch
GOOS=linux GOARCH=amd64 go build -o tinymonitor-linux-amd64 ./cmd/tinymonitor
GOOS=linux GOARCH=arm64 go build -o tinymonitor-linux-arm64 ./cmd/tinymonitor
GOOS=darwin GOARCH=arm64 go build -o tinymonitor-darwin-arm64 ./cmd/tinymonitor
```