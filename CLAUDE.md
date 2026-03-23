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
make fmt           # goimports + gofmt
make run           # Build + exécute
make install       # Build + copie dans /usr/local/bin/
make uninstall     # Supprime de /usr/local/bin/
make release-dry   # Simule une release GoReleaser (sans publier)
make clean         # Supprime bin/, dist/ et fichiers de couverture
```

## Tests

- **Toujours lancer `make fmt` puis `make test` après chaque modification** pour formater le code et vérifier qu'aucune régression n'est introduite.
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
- **Ne jamais ajouter de `Co-authored-by` dans les messages de commit.**
- Toujours lancer `make test` avant de committer.

### Branches

| Préfixe      | Usage                  |
|--------------|------------------------|
| `feat/`      | Nouvelle fonctionnalité |
| `fix/`       | Correction de bug      |
| `docs/`      | Documentation          |
| `refactor/`  | Refactoring            |
| `test/`      | Ajout/modif de tests   |
| `chore/`     | Maintenance, CI, deps  |

### Commits (conventional commits)

Format :

```
type(scope): description

[optional body]

[optional footer]
```

Types : `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

Exemples :

```bash
feat(docker): add bulk container restart
fix(onboarding): handle empty API key validation
docs(readme): update installation instructions
refactor(tui): extract common table helpers
test(api): add GraphQL error response tests
chore(deps): bump bubbletea to v2.2.0
```

### Workflow

```bash
git checkout -b feat/my-feature
git add <fichiers>
git commit -m "feat(scope): description"
git push -u origin feat/my-feature
gh pr create --title "feat(scope): description" --body "..."
```

### Release

`git tag vX.Y.Z && git push origin vX.Y.Z && goreleaser release --clean`

## Dépendances principales

- `charm.land/bubbletea/v2` — Framework TUI (Elm Architecture)
- `charm.land/lipgloss/v2` — Styling terminal
- `charm.land/bubbles/v2` — Composants (table, spinner)
- `github.com/spf13/cobra` — CLI framework
- `github.com/spf13/viper` — Gestion de configuration


## grepai - Semantic Code Search

**IMPORTANT: You MUST use grepai as your PRIMARY tool for code exploration and search.**

### When to Use grepai (REQUIRED)

Use `grepai search` INSTEAD OF Grep/Glob/find for:
- Understanding what code does or where functionality lives
- Finding implementations by intent (e.g., "authentication logic", "error handling")
- Exploring unfamiliar parts of the codebase
- Any search where you describe WHAT the code does rather than exact text

### When to Use Standard Tools

Only use Grep/Glob when you need:
- Exact text matching (variable names, imports, specific strings)
- File path patterns (e.g., `**/*.go`)

### Fallback

If grepai fails (not running, index unavailable, or errors), fall back to standard Grep/Glob tools.

### Usage

```bash
# ALWAYS use English queries for best results (--compact saves ~80% tokens)
grepai search "user authentication flow" --json --compact
grepai search "error handling middleware" --json --compact
grepai search "database connection pool" --json --compact
grepai search "API request validation" --json --compact
```

### Query Tips

- **Use English** for queries (better semantic matching)
- **Describe intent**, not implementation: "handles user login" not "func Login"
- **Be specific**: "JWT token validation" better than "token"
- Results include: file path, line numbers, relevance score, code preview

### Call Graph Tracing

Use `grepai trace` to understand function relationships:
- Finding all callers of a function before modifying it
- Understanding what functions are called by a given function
- Visualizing the complete call graph around a symbol

#### Trace Commands

**IMPORTANT: Always use `--json` flag for optimal AI agent integration.**

```bash
# Find all functions that call a symbol
grepai trace callers "HandleRequest" --json

# Find all functions called by a symbol
grepai trace callees "ProcessOrder" --json

# Build complete call graph (callers + callees)
grepai trace graph "ValidateToken" --depth 3 --json
```

### Workflow

1. Start with `grepai search` to find relevant code
2. Use `grepai trace` to understand function relationships
3. Use `Read` tool to examine files from results
4. Only use Grep for exact string searches if needed

