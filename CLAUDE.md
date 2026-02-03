# TinyMonitor - Go

## Description
TinyMonitor est un agent de monitoring système léger écrit en Go avec alertes multi-canaux.

**Fonction** : Surveillance système (CPU, RAM, disque, load, I/O, reboot) avec alertes configurables
**Architecture** : Boucle de polling + collectors metrics + worker pool d'alertes + providers multi-canaux
**Cible** : Binaire standalone pour Linux/macOS (AMD64 + ARM64), ~9MB
**Framework CLI** : Cobra avec 10 commandes (run, validate, info, test-alert, service, update, uninstall, version, completion)

## Architecture Complète

### Structure des Répertoires
```
tinymonitor/
├── main.go                      # Entry point (délègue à cmd.Execute())
├── cmd/                         # Commandes CLI (Cobra)
│   ├── root.go                 # Commande principale + logging setup
│   ├── version.go              # Affichage version/commit/date
│   ├── validate.go             # Validation config TOML
│   ├── info.go                 # Résumé config (metrics, alerts)
│   ├── test_alert.go           # Test des providers d'alertes
│   ├── service.go              # Gestion service systemd
│   ├── uninstall.go            # Suppression service/binaire/config
│   ├── update.go               # Self-update depuis GitHub releases
│   └── config.example.toml     # Template config embarqué
│
├── internal/
│   ├── config/
│   │   ├── config.go           # Parsing TOML + validation (446 lignes)
│   │   └── config_test.go      # Tests de validation
│   │
│   ├── monitor/
│   │   └── monitor.go          # Boucle principale (ticker, state machine, 225 lignes)
│   │
│   ├── models/
│   │   └── types.go            # Types partagés (Alert, MetricResult, Severity, AlertState)
│   │
│   ├── metrics/                # Collectors de métriques
│   │   ├── collector.go        # Interface Collector
│   │   ├── cpu.go              # CPU usage (gopsutil)
│   │   ├── memory.go           # Mémoire (RAM/Swap)
│   │   ├── disk.go             # Filesystem usage
│   │   ├── load.go             # Load average (1/5/15 min)
│   │   ├── io.go               # I/O disque
│   │   ├── reboot.go           # Détection reboot système
│   │   └── metrics_test.go
│   │
│   ├── alerts/                 # Providers d'alertes
│   │   ├── provider.go         # Interface Provider + BaseProvider
│   │   ├── manager.go          # Worker pool (5 goroutines, channel-based)
│   │   ├── ntfy.go             # Push notifications (ntfy.sh)
│   │   ├── google_chat.go      # Google Chat webhook
│   │   ├── smtp.go             # Email SMTP
│   │   ├── webhook.go          # Generic webhook POST
│   │   └── gotify.go           # Gotify self-hosted
│   │
│   ├── system/
│   │   ├── service.go          # Installation systemd (template embarqué, 8663 bytes)
│   │   ├── uninstall.go        # Suppression service/config
│   │   └── update.go           # Self-update (GitHub API, atomic replacement)
│   │
│   └── utils/
│       └── system_info.go      # Helpers système (hostname, uptime)
│
├── configs/
│   └── config.example.toml     # Config exemple complète (3796 bytes)
│
├── docs/                        # Documentation Zensical (mkdocs-like)
│   ├── index.md                # Homepage
│   ├── configuration.md        # Guide config TOML
│   ├── development.md          # Guide développeur
│   ├── commands/               # Doc des 10 commandes CLI
│   ├── metrics/                # Doc des 6 types de métriques
│   ├── alerts/                 # Doc des 5 providers d'alertes
│   └── getting-started/        # Quick start
│
├── install/
│   └── install.sh              # Script d'installation cross-platform (7848 bytes)
│
├── .github/workflows/
│   ├── ci-cd.yml               # Tests + vet sur push/PR
│   ├── release.yml             # GoReleaser sur tag v*
│   └── docs.yml                # Déploiement docs sur GitHub Pages
│
├── Makefile                     # Automation build/test/lint
├── .goreleaser.yaml            # Config releases multi-arch
├── zensical.toml               # Config documentation builder
├── go.mod                       # Go 1.21+
└── go.sum
```

### Flux de Données
1. **Configuration** : Fichier TOML → Parsing → Validation → Structs typées
2. **Monitoring** : Ticker (interval configurable) → Collectors → MetricResults
3. **Alertes** : MetricResult → State machine (warning/critical/recovery) → Alert channel → Worker pool → Providers
4. **Rate Limiting** : Cooldown period + duration-based triggers (évite faux positifs)

### Patterns Utilisés
- **Worker Pool** : 5 goroutines pour distribution des alertes (alerts/manager.go)
- **State Machine** : Suivi des niveaux d'alerte (OK → Warning → Critical → Recovery)
- **Interface-based Design** : `Collector` (metrics), `Provider` (alerts)
- **Graceful Shutdown** : Context cancellation + sync.WaitGroup
- **Embedded Resources** : Config template et systemd unit embarqués dans le binaire

## Conventions Go

### Standards de Code
- **Erreurs** : Gestion explicite, jamais de `panic` en production
- **Context** : Utilisé pour timeouts, cancellation, graceful shutdown
- **Interfaces** : Définies par le consommateur, suffixe `er` (Collector, Provider)
- **Logging** : `log/slog` (stdlib Go 1.21+), multi-handler (stdout + fichier optionnel)
- **Tests** : Table-driven tests dans fichiers `*_test.go`
- **Naming** : Constructeurs `New<Type>()`, configs `<Name>Config`

### Patterns de Concurrence
- Channels pour communication entre goroutines
- sync.WaitGroup pour coordination
- Context pour propagation de cancellation
- Mutex uniquement si nécessaire (évité au maximum)

### Organisation du Code
- Single responsibility par package
- Pas de state global (sauf logger)
- Minimal dependencies (stdlib préférée)
- Clear separation : cmd/ (CLI), internal/ (logic)

## Directives de Développement

### Workflow & Méthodologie
- **Clarification Proactive** : Ne jamais hésiter à demander des précisions si une demande est ambiguë, incomplète, ou si plusieurs approches sont possibles. Utiliser `AskUserQuestion` pour clarifier les exigences avant de commencer l'implémentation.
- **Mode Plan** : Toujours fonctionner en "mode plan" (Planning Mode). Avant toute modification de code, proposer une analyse détaillée de la tâche avec les étapes prévues. Utiliser `EnterPlanMode` pour les tâches non-triviales.
- **Documentation Systématique** : Mettre à jour la documentation (`docs/`) pour toute nouvelle fonctionnalité ou changement de logique. Ajouter des commentaires dans le code uniquement quand la logique n'est pas évidente.
- **Maintien du CLAUDE.md** : Mettre à jour ce fichier d'instructions (`CLAUDE.md`) lors de changements architecturaux significatifs : ajout/suppression de dépendances, nouvelles commandes CLI, changements de workflow, nouvelles conventions, ou modifications de structure de répertoires.
- **Documentation Officielle** : Utiliser la documentation officielle via Context7 (MCP tool `mcp__context7__query-docs`) pour valider les implémentations liées à Go, Cobra, ou toute autre dépendance externe.
- **Tests Obligatoires** : Ne jamais considérer une tâche comme terminée sans avoir écrit ou mis à jour les tests unitaires (`*_test.go`). Lancer `make test` et `make test-race` pour validation complète.
- **Validation Continue** : Exécuter `make vet` et `make lint` avant tout commit pour garantir la qualité du code.
- **Commits en Anglais** : Tous les messages de commit doivent être rédigés en anglais en suivant la convention Conventional Commits (feat:, fix:, docs:, refactor:, test:, etc.).

### Style de Code & Langue
- **Langue** : Le code, les commentaires, les messages de commit, et la documentation technique doivent être rédigés exclusivement en **anglais**.
- **Qualité** : Suivre scrupuleusement les standards `golangci-lint` définis dans le projet. Aucun warning ne doit subsister.
- **Cohérence** : Respecter les patterns existants (naming conventions, error handling, logging format) pour maintenir l'homogénéité du codebase.

### Checklist Avant Commit
1. ✅ Code et commentaires en anglais
2. ✅ Message de commit en anglais (format: `type: description`)
3. ✅ `make test` passe sans erreur
4. ✅ `make test-race` ne détecte pas de race conditions
5. ✅ `make vet` ne remonte aucun problème
6. ✅ `make lint` ne produit aucun warning
7. ✅ Documentation mise à jour (`docs/` si nécessaire)
8. ✅ Tests unitaires ajoutés/modifiés si applicable

## Dépendances

### Go Modules (go.mod)
```go
require (
    github.com/shirou/gopsutil/v3  // Métriques système cross-platform
    github.com/BurntSushi/toml     // Parsing TOML
    github.com/spf13/cobra         // Framework CLI
    golang.org/x/sys               // Syscalls bas niveau
)
```

### Outils de Développement
- **golangci-lint** : Linting (configuré dans Makefile)
- **GoReleaser** : Automation des releases (v2)
- **Zensical** : Builder de documentation (Python, config dans zensical.toml)
- **GitHub Actions** : CI/CD (3 workflows : ci-cd, release, docs)

### Versions Minimales
- Go 1.21+ (requis pour `log/slog`)
- Linux : Systemd (pour service management)
- macOS : LaunchD (documentation fournie)

## Commandes de Développement

### Makefile Targets
```bash
make build          # Build pour plateforme actuelle
make build-all      # Build multi-arch (Linux AMD64/ARM64, macOS x86_64/ARM64)
make test           # Run tests unitaires
make test-cover     # Coverage report
make test-race      # Race detector
make vet            # Go vet static analysis
make lint           # golangci-lint checks
make clean          # Supprime dist/ et binaires
make install        # Copie vers /usr/local/bin
make uninstall      # Supprime de /usr/local/bin
make completions    # Génère shell completions (bash/zsh/fish)
make dev            # Build + run avec config par défaut
make help           # Affiche aide
```

### Build Flags Importants
```bash
# LDFLAGS injectés au build
-s -w                           # Strip debug info (réduit taille)
-X cmd.Version=$(VERSION)       # Version depuis git tag
-X cmd.Commit=$(COMMIT)         # Commit SHA
-X cmd.BuildDate=$(DATE)        # Date de build ISO8601
```

### Tests
```bash
go test ./...                   # Tous les tests
go test -v ./internal/config    # Tests config avec verbose
go test -race ./...             # Race detector
go test -cover ./...            # Coverage
```

### Build Manuel Cross-Platform
```bash
GOOS=linux GOARCH=amd64 go build -o tinymonitor-linux-amd64 .
GOOS=linux GOARCH=arm64 go build -o tinymonitor-linux-arm64 .
GOOS=darwin GOARCH=amd64 go build -o tinymonitor-darwin-amd64 .
GOOS=darwin GOARCH=arm64 go build -o tinymonitor-darwin-arm64 .
```

## CI/CD Pipeline

### GitHub Actions Workflows

**ci-cd.yml** (Tests Continus)
- Trigger : Push sur main (*.go, go.mod, go.sum) + Pull Requests
- Runner : ubuntu-latest
- Steps : Checkout → Setup Go 1.21 → Dependencies → `go vet ./...` → `go test -v ./...`

**release.yml** (Releases Automatiques)
- Trigger : Push de tag `v*` (ex: v1.2.3)
- Runner : ubuntu-latest
- Steps : Checkout → Go setup → `go mod tidy` → `go test ./...` → GoReleaser
- Outputs : Binaires multi-arch, checksums SHA256, release notes, install.sh

**docs.yml** (Documentation Déployée)
- Trigger : Changes dans `docs/**` ou `zensical.toml`
- Runner : ubuntu-latest + Python
- Steps : Checkout → Setup Python → Install Zensical → Build docs → Deploy vers GitHub Pages
- URL : https://gu1llaum-3.github.io/tinymonitor/

## Workflow de Release

### Process Complet
1. **Commit les changements** : `git commit -m "feat: nouvelle fonctionnalité"`
2. **Tag version** : `git tag -a v1.2.3 -m "Release v1.2.3"`
3. **Push tag** : `git push origin v1.2.3`
4. **GitHub Actions** : Déclenche automatiquement release.yml
5. **GoReleaser** : Build multi-arch, génère checksums, crée GitHub release
6. **Assets** : Binaires + install.sh uploadés sur release

### Format des Tags
- Semantic versioning : `v<major>.<minor>.<patch>`
- Exemple : `v1.2.3`, `v2.0.0-beta.1`

## Documentation

### Structure Zensical
- **Builder** : Zensical (config dans `zensical.toml`)
- **Format** : Markdown dans `docs/`
- **Navigation** : Auto-générée depuis structure de fichiers
- **Déploiement** : GitHub Pages via workflow docs.yml
- **Commandes Documentées** : 10 CLI commands
- **Métriques Documentées** : 6 types (CPU, Memory, Disk, Load, I/O, Reboot)
- **Providers Documentés** : 5 types (Ntfy, Google Chat, SMTP, Webhook, Gotify)

### Sections Principales
1. **Getting Started** : Installation, first config
2. **Configuration** : Format TOML, search paths, validation
3. **Commands** : Usage de chaque commande Cobra
4. **Metrics** : Détails de chaque collector
5. **Alerts** : Configuration de chaque provider
6. **Guides** : Troubleshooting, LaunchD setup
7. **Development** : Architecture, contribution, ajout de features

## Usage CLI

### Commandes Principales
```bash
# Run monitoring (cherche config dans : CLI flag > ./config.toml > ~/.config/tinymonitor/config.toml > /etc/tinymonitor/config.toml)
tinymonitor
tinymonitor -c /path/to/config.toml

# Configuration
tinymonitor info                          # Affiche résumé config
tinymonitor validate                      # Valide config sans lancer monitoring
tinymonitor validate -c custom.toml       # Valide config spécifique

# Tests
tinymonitor test-alert                    # Test tous les providers
tinymonitor test-alert --provider ntfy    # Test provider spécifique

# Service Management (Linux systemd)
sudo tinymonitor service install          # Installe service avec config par défaut
sudo tinymonitor service install -c /etc/tinymonitor/config.toml
tinymonitor service status                # Check status
sudo tinymonitor service uninstall        # Supprime service (garde config)

# Update
tinymonitor update --check                # Vérifie nouvelle version
tinymonitor update                        # Met à jour depuis GitHub
tinymonitor update --yes                  # Update sans confirmation

# Uninstall
sudo tinymonitor uninstall                # Supprime service + binaire (garde config)
sudo tinymonitor uninstall --purge        # Supprime tout (service + binaire + config)

# Utilities
tinymonitor version                       # Version, commit, build date
tinymonitor completion bash               # Shell completion script
```

### Chemins de Configuration (ordre de priorité)
1. Flag CLI : `-c /custom/path/config.toml`
2. Répertoire actuel : `./config.toml`
3. User config : `~/.config/tinymonitor/config.toml`
4. System config : `/etc/tinymonitor/config.toml`

## Configuration TOML

### Exemple Minimal
```toml
refresh = 5       # Intervalle de polling (secondes)
cooldown = 60     # Cooldown entre alertes (secondes)

[cpu]
enabled = true
warning = 70
critical = 90
duration = 30

[memory]
enabled = true
warning = 80
critical = 95

[alerts.ntfy]
enabled = true
topic_url = "https://ntfy.sh/your_topic"
```

### Sections Principales
- **Global** : `refresh`, `cooldown`, `log_file`
- **Metrics** : `cpu`, `memory`, `filesystem`, `load`, `io`, `reboot`
- **Alerts** : `alerts.ntfy`, `alerts.google_chat`, `alerts.smtp`, `alerts.webhook`, `alerts.gotify`

Voir `configs/config.example.toml` pour config complète avec tous les paramètres.

## Ajout de Nouvelles Features

### Ajouter un Nouveau Collector
1. Créer `internal/metrics/new_metric.go`
2. Implémenter interface `Collector` :
   ```go
   type Collector interface {
       Name() string
       Collect() (*models.MetricResult, error)
   }
   ```
3. Ajouter struct config dans `internal/config/config.go`
4. Enregistrer dans `monitor.go` : `collectors = append(collectors, metrics.NewMyMetric(cfg))`
5. Ajouter tests dans `metrics_test.go`
6. Documenter dans `docs/metrics/new_metric.md`

### Ajouter un Nouveau Provider
1. Créer `internal/alerts/new_provider.go`
2. Implémenter interface `Provider` :
   ```go
   type Provider interface {
       Name() string
       Send(alert *models.Alert) error
   }
   ```
3. Ajouter struct config dans `internal/config/config.go`
4. Enregistrer dans `monitor.go` : `providers = append(providers, alerts.NewMyProvider(cfg))`
5. Ajouter support dans `cmd/test_alert.go`
6. Documenter dans `docs/alerts/new_provider.md`

### Ajouter une Nouvelle Commande Cobra
1. Créer `cmd/my_command.go`
2. Définir commande :
   ```go
   var myCmd = &cobra.Command{
       Use:   "mycommand",
       Short: "Description courte",
       Long:  "Description longue",
       Run: func(cmd *cobra.Command, args []string) {
           // Logic
       },
   }
   ```
3. Enregistrer dans `cmd/root.go` : `rootCmd.AddCommand(myCmd)`
4. Documenter dans `docs/commands/my_command.md`

## Notes Importantes

### Sécurité
- Service systemd avec hardening (PrivateTmp, NoNewPrivileges, etc.)
- Pas de credentials en clair dans les logs
- Validation stricte des configs (types, ranges)

### Performance
- Binaire ~9MB (stripped)
- Low CPU usage (polling loop avec ticker)
- Low RAM usage (pas de buffers massifs)
- Pas de dépendances runtime

### Compatibilité
- Linux : Systemd requis pour service management
- macOS : LaunchD (doc fournie dans docs/guides/)
- Windows : Non supporté (gopsutil compatible mais pas testé)

### Limites Connues
- Un seul fichier de config supporté
- Pas de hot-reload de config (restart requis)
- Metrics collectors synchrones (pas de parallélisation)
- Worker pool fixe (5 workers)

## Ressources

- **Repository** : https://github.com/Gu1llaum-3/tinymonitor
- **Documentation** : https://gu1llaum-3.github.io/tinymonitor/
- **Releases** : https://github.com/Gu1llaum-3/tinymonitor/releases
- **Issues** : https://github.com/Gu1llaum-3/tinymonitor/issues
- **License** : MIT
