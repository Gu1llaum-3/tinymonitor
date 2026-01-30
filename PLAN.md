# Plan de migration vers Cobra CLI

## Analyse du code actuel

### Gestion actuelle (stdlib `flag`)
```go
configPath := flag.String("c", "", "Path to configuration file")
flag.StringVar(configPath, "config", "", "Path to configuration file")
showVersion := flag.Bool("v", false, "Show version")
flag.BoolVar(showVersion, "version", false, "Show version")
```

### Fonctionnalités existantes
- `-c` / `--config` : chemin vers le fichier de configuration
- `-v` / `--version` : affiche la version
- Chargement config avec cascade de chemins
- Logging multi-handler (stdout + fichier)
- Graceful shutdown (SIGINT/SIGTERM)

## Architecture Cobra cible

```
cmd/tinymonitor/
├── main.go           # Point d'entrée (appelle cmd.Execute())
├── cmd/
│   ├── root.go       # Commande racine (run monitor)
│   └── version.go    # Sous-commande version
```

### Commandes prévues
- `tinymonitor` : Lance le monitoring (comportement par défaut)
- `tinymonitor version` : Affiche la version détaillée
- `tinymonitor --help` : Aide auto-générée

### Flags globaux (persistent)
- `-c, --config` : Chemin vers le fichier de configuration

---

## Plan d'exécution

### Phase 1 : Setup Cobra
- [x] Ajouter la dépendance `github.com/spf13/cobra`
- [x] Créer le dossier `cmd/tinymonitor/cmd/`
- [x] Compiler et valider

### Phase 2 : Commande racine (root.go)
- [x] Créer `cmd/tinymonitor/cmd/root.go` avec la commande racine
- [x] Migrer le flag `--config` comme PersistentFlag
- [x] Migrer la logique de run (config, logging, signals, monitor)
- [x] Compiler et tester

### Phase 3 : Commande version
- [x] Créer `cmd/tinymonitor/cmd/version.go`
- [x] Ajouter variables pour Version, Commit, BuildDate (ldflags)
- [x] Compiler et tester `tinymonitor version`

### Phase 4 : Refactoring main.go
- [x] Simplifier main.go (appel unique à cmd.Execute())
- [x] Déplacer multiHandler dans un package séparé ou inline
- [x] Compiler et tester

### Phase 5 : Mise à jour ldflags
- [x] Mettre à jour Makefile pour injecter Version/Commit/Date
- [x] Mettre à jour .goreleaser.yaml
- [x] Tester le build avec version injectée

### Phase 6 : Validation finale
- [x] `go vet ./...`
- [x] `go test ./...`
- [x] Test manuel : `tinymonitor --help`
- [x] Test manuel : `tinymonitor version`
- [x] Test manuel : `tinymonitor -c config.json`

---

## Résultat attendu

```bash
$ tinymonitor --help
TinyMonitor - Lightweight system monitoring agent

Usage:
  tinymonitor [flags]
  tinymonitor [command]

Available Commands:
  version     Print version information
  help        Help about any command

Flags:
  -c, --config string   Path to configuration file
  -h, --help            Help for tinymonitor

$ tinymonitor version
TinyMonitor v1.0.0
Commit: abc1234
Built: 2024-01-30T12:00:00Z
```
