# Plan de migration TinyMonitor Python → Go

## Instructions pour Claude
Exécute ce plan étape par étape. Après chaque étape, compile avec `go build ./...` et corrige les erreurs avant de passer à la suivante. Demande confirmation avant de passer à une nouvelle phase.

## Phase 1 : Initialisation
- [x] Créer `tinymonitor-go/` avec `go mod init github.com/Gu1llaum-3/tinymonitor`
- [x] Créer l'arborescence de dossiers (cmd/, internal/, etc.)
- [x] Créer main.go minimal qui compile

## Phase 2 : Types de base
- [x] Créer internal/models/types.go (Severity, Alert, MetricResult)
- [x] Compiler et valider

## Phase 3 : Configuration
- [x] Analyser config.py
- [x] Créer internal/config/config.go avec structs typées
- [x] Implémenter Load() avec cascade de chemins
- [x] Écrire tests config_test.go
- [x] Compiler et valider

## Phase 4 : Metrics
- [x] Créer internal/metrics/collector.go (interface)
- [x] Migrer cpu.py → cpu.go
- [x] Migrer memory.py → memory.go
- [x] Migrer disk.py → disk.go
- [x] Migrer load.py → load.go
- [x] Migrer io.py → io.go
- [x] Migrer reboot.py → reboot.go
- [x] Tests pour chaque collector

## Phase 5 : Alertes
- [x] Créer internal/alerts/provider.go (interface)
- [x] Créer internal/alerts/manager.go (worker pool)
- [x] Migrer console provider
- [x] Migrer google_chat.go
- [x] Migrer ntfy.go
- [x] Migrer smtp.go
- [x] Migrer webhook.go
- [x] Migrer gotify.go

## Phase 6 : Core
- [x] Migrer monitor.py → internal/monitor/monitor.go
- [x] Compléter cmd/tinymonitor/main.go
- [x] Gestion des signaux (graceful shutdown)

## Phase 7 : Validation finale
- [x] go vet ./...
- [x] go test ./...
- [x] Build multi-arch
- [x] Test manuel avec config.json existante