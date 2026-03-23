# Shell Completion

`unraid-tui` supports autocompletion for commands, flags, and dynamic values (server names, languages).

**Zsh** (add to `~/.zshrc`):

```bash
eval "$(unraid-tui completion zsh)"
```

**Oh-My-Zsh** plugin:

```bash
mkdir -p ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/unraid-tui
unraid-tui completion zsh > ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/unraid-tui/_unraid-tui
# Then add "unraid-tui" to plugins=(...) in ~/.zshrc
```

**Bash:**

```bash
# Linux
unraid-tui completion bash > /etc/bash_completion.d/unraid-tui

# macOS (requires bash-completion@2)
unraid-tui completion bash > $(brew --prefix)/etc/bash_completion.d/unraid-tui
```

**Fish:**

```fish
unraid-tui completion fish > ~/.config/fish/completions/unraid-tui.fish
```

**PowerShell:**

```powershell
unraid-tui completion powershell | Out-String | Invoke-Expression
```

> If you installed via Homebrew, Zsh completions are handled automatically.
