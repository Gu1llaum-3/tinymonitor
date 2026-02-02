# tinymonitor completion

Generate shell autocompletion scripts for TinyMonitor commands.

## Usage

```bash
tinymonitor completion <shell>
```

## Supported Shells

| Shell | Command |
|-------|---------|
| Bash | `tinymonitor completion bash` |
| Zsh | `tinymonitor completion zsh` |
| Fish | `tinymonitor completion fish` |
| PowerShell | `tinymonitor completion powershell` |

## Installation

### Bash

=== "Linux"
    ```bash
    tinymonitor completion bash | sudo tee /etc/bash_completion.d/tinymonitor > /dev/null
    ```

=== "macOS (Homebrew)"
    ```bash
    tinymonitor completion bash > $(brew --prefix)/etc/bash_completion.d/tinymonitor
    ```

### Zsh

First, ensure completion is enabled in your `~/.zshrc`:

```bash
autoload -U compinit; compinit
```

Then add the completion:

```bash
tinymonitor completion zsh > "${fpath[1]}/_tinymonitor"
```

Or for Oh My Zsh:

```bash
tinymonitor completion zsh > ~/.oh-my-zsh/completions/_tinymonitor
```

### Fish

```bash
tinymonitor completion fish > ~/.config/fish/completions/tinymonitor.fish
```

### PowerShell

Add to your PowerShell profile:

```powershell
tinymonitor completion powershell | Out-String | Invoke-Expression
```

To make it permanent, add the above line to your `$PROFILE` file.

## Usage After Installation

Restart your shell or source the completion file. Then use Tab to autocomplete:

```bash
tinymonitor ser<TAB>
# Completes to: tinymonitor service

tinymonitor service <TAB>
# Shows: install  status  uninstall

tinymonitor update --<TAB>
# Shows: --check  --yes
```

## Updating Completions

After updating TinyMonitor, regenerate the completion scripts to include any new commands or flags:

```bash
# Example for Bash on Linux
tinymonitor completion bash | sudo tee /etc/bash_completion.d/tinymonitor > /dev/null
```
