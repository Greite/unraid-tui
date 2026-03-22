# CLAUDE.md — unraid-tui

## Projet

CLI TUI en Go pour l'API Unraid (GraphQL). Monitoring serveur avec dashboard et gestion Docker.

## Commandes

```bash
make build         # Compile le binaire dans bin/ (avec ldflags version)
make test          # Lance tous les tests
make test-verbose  # Tests avec détail
make test-cover    # Tests avec couverture HTML
make lint          # go vet
make run           # Build + exécute
make install       # Build + copie dans /usr/local/bin/
make uninstall     # Supprime de /usr/local/bin/
make release-dry   # Simule une release GoReleaser (sans publier)
make clean         # Supprime bin/, dist/ et fichiers de couverture
```

## Tests

- **Toujours lancer `make test` après chaque modification** pour vérifier qu'aucune régression n'est introduite.
- Les tests utilisent `httptest.NewServer` pour mocker l'API GraphQL (pas de serveur réel nécessaire).
- Le mock client est dans `internal/api/mock.go` — l'utiliser pour tous les tests TUI.
- Les tests TUI testent le modèle (Update/View), pas le rendu terminal réel.
- Pas de dépendance externe pour les assertions — utiliser `testing` standard avec `t.Errorf` / `t.Fatal`.

## Architecture

```
cmd/                     → Point d'entrée Cobra
internal/api/            → Client GraphQL (interface UnraidClient + HTTP)
internal/config/         → Chargement config Viper + détection/sauvegarde
internal/model/          → Types domaine partagés
internal/tui/            → App Bubbletea (routeur de pages)
internal/tui/common/     → Styles, messages, helpers (package séparé pour éviter les cycles d'import)
internal/tui/dashboard/  → Page dashboard
internal/tui/docker/     → Page Docker
internal/tui/onboarding/ → Assistant de configuration au premier lancement
```

## Conventions

- **Import paths Charmbracelet v2** : utiliser `charm.land/bubbletea/v2`, `charm.land/lipgloss/v2`, `charm.land/bubbles/v2` (pas `github.com/charmbracelet/...`).
- **Bubbletea v2** : `View()` retourne `tea.View` (pas `string`), pas de `tea.WithAltScreen()` (utiliser `v.AltScreen = true` dans la View), spinner n'a pas de `Init()` (utiliser `m.spinner.Tick` comme Cmd initial).
- **Pas de cycles d'import** : les sous-packages TUI (dashboard, docker) importent `tui/common`, jamais `tui` directement.
- **Interface `UnraidClient`** : tout accès API passe par cette interface pour la testabilité.
- **Pas de `fmt.Sprintf` pour les assertions de test** — utiliser les helpers de `testing` directement.

## Config

Fichier `~/.unraid-tui/config.yaml` ou variables d'environnement `UNRAID_SERVER_URL` et `UNRAID_API_KEY`.

## Onboarding

- Au premier lancement (pas de config détectée), l'onboarding TUI se lance automatiquement.
- `config.Exists()` vérifie si la config est complète (fichier + env vars).
- `config.Save()` écrit le fichier avec permissions `0600`.
- L'onboarding teste la connexion au serveur **et** la validité de la clé API avant de sauvegarder.
- Les sous-modèles d'onboarding (textinput, spinner) suivent les mêmes patterns que les pages TUI.
- Les tests d'onboarding vérifient chaque transition d'étape indépendamment, sans serveur réel.

## Release

- GoReleaser gère la compilation cross-platform et la publication.
- Config dans `.goreleaser.yaml`.
- Les variables `version`, `commit`, `date` dans `cmd/root.go` sont injectées par ldflags au build.
- La formule Homebrew est publiée automatiquement dans `Greite/homebrew-tap`.
- Pour tagger : `git tag vX.Y.Z && git push origin vX.Y.Z` puis `goreleaser release --clean`.

## Git & GitHub

- Le repo distant est sur GitHub sous l'organisation/utilisateur **Greite**.
- Utiliser `gh` (GitHub CLI) pour les opérations GitHub (PRs, issues, releases).
- Committer avec des messages concis en anglais, au présent impératif (ex: `add onboarding wizard`, `fix config loading`).
- **Ne jamais ajouter de `Co-authored-by` dans les messages de commit.**
- Toujours lancer `make test` avant de committer.
- Workflow de commit et push :
  ```bash
  git add <fichiers>
  git commit -m "message"
  git push
  ```
- Créer une PR : `gh pr create --title "..." --body "..."`.
- Créer une release : `git tag vX.Y.Z && git push origin vX.Y.Z && goreleaser release --clean`.

## Dépendances principales

- `charm.land/bubbletea/v2` — Framework TUI (Elm Architecture)
- `charm.land/lipgloss/v2` — Styling terminal
- `charm.land/bubbles/v2` — Composants (table, spinner)
- `github.com/spf13/cobra` — CLI framework
- `github.com/spf13/viper` — Gestion de configuration
