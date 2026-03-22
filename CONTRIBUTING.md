# Contribuer a unraid-tui

Merci de votre interet pour ce projet ! Ce guide vous explique comment contribuer.

## Prerequis

- [Go 1.22+](https://go.dev/dl/)
- [GNU Make](https://www.gnu.org/software/make/)
- (Optionnel) [GoReleaser](https://goreleaser.com/) pour tester les releases
- (Optionnel) [GitHub CLI](https://cli.github.com/) pour les PRs

## Mise en place

```bash
git clone https://github.com/Greite/unraid-tui.git
cd unraid-tui
make build
make test
```

Si le build et les tests passent, vous etes pret.

## Workflow de developpement

### 1. Creer une branche

```bash
git checkout -b feat/ma-feature
```

Conventions de nommage :

| Prefixe    | Usage                     |
|------------|---------------------------|
| `feat/`    | Nouvelle fonctionnalite   |
| `fix/`     | Correction de bug         |
| `refactor/`| Refactoring               |
| `docs/`    | Documentation uniquement  |
| `test/`    | Ajout/modification de tests|

### 2. Developper

```bash
# Lancer les tests en continu pendant le dev
make test

# Verifier que le code compile
make build

# Lancer le linter
make lint
```

### 3. Committer

Messages de commit en anglais, au present imperatif :

```bash
# Bon
git commit -m "add VM monitoring page"
git commit -m "fix container port display when host port is 0"

# Mauvais
git commit -m "Added VM page"
git commit -m "WIP"
git commit -m "fix stuff"
```

### 4. Ouvrir une PR

```bash
git push -u origin feat/ma-feature
gh pr create --title "Add VM monitoring page" --body "Description..."
```

## Structure du projet

```
cmd/                     → Point d'entree Cobra
internal/api/            → Client GraphQL (interface + HTTP)
internal/config/         → Configuration Viper
internal/model/          → Types domaine
internal/tui/            → App Bubbletea (routeur)
internal/tui/common/     → Styles, messages, helpers partages
internal/tui/dashboard/  → Page dashboard
internal/tui/docker/     → Page Docker
internal/tui/onboarding/ → Assistant de configuration
```

## Ajouter une nouvelle page TUI

1. Creer un package dans `internal/tui/<nom>/`
2. Implementer un `Model` avec `Init()`, `Update()`, `View()` (pattern Bubbletea)
3. Importer les styles et messages depuis `internal/tui/common/` (jamais depuis `internal/tui/` directement pour eviter les cycles d'import)
4. Ajouter la page dans le routeur `internal/tui/app.go`
5. Ajouter les requetes GraphQL dans `internal/api/queries.go` si necessaire
6. Ecrire les tests

## Ajouter une requete API

1. Ajouter la requete GraphQL dans `internal/api/queries.go`
2. Ajouter les types de reponse dans `internal/api/types.go` avec une methode `toDomain()`
3. Ajouter les types domaine dans `internal/model/model.go`
4. Ajouter la methode dans l'interface `UnraidClient` (`internal/api/client.go`)
5. Implementer la methode dans `httpClient`
6. Ajouter la methode dans `MockClient` (`internal/api/mock.go`)
7. Tester avec `httptest` dans `internal/api/client_test.go`

## Tests

Les tests sont obligatoires pour toute contribution.

```bash
make test          # Lancer les tests
make test-verbose  # Avec detail
make test-cover    # Avec couverture
```

### Regles

- Utiliser le package `testing` standard (pas de framework externe)
- Mocker l'API avec `httptest.NewServer` pour les tests client
- Utiliser `api.MockClient` pour les tests TUI
- Tester le modele Bubbletea (Update/View), pas le rendu terminal
- Chaque nouvelle feature doit etre accompagnee de tests

### Lancer un test specifique

```bash
go test -v -run TestNomDuTest ./internal/tui/docker/
```

## Conventions de code

- **Import paths Charmbracelet v2** : `charm.land/bubbletea/v2`, `charm.land/lipgloss/v2`, `charm.land/bubbles/v2`
- **Bubbletea v2** :
  - `View()` retourne `tea.View`, pas `string`
  - Pas de `tea.WithAltScreen()` — utiliser `v.AltScreen = true` dans la View
  - Spinner : utiliser `m.spinner.Tick` comme Cmd initial (pas de `Init()`)
- **Pas de cycles d'import** : les sous-packages TUI importent `tui/common`, jamais `tui`
- **Interface `UnraidClient`** : tout acces API passe par cette interface
- Formater le code avec `gofmt` (applique automatiquement par Go)

## Documentation

- Documenter les nouvelles features dans `docs/`
- Mettre a jour le `README.md` si la feature est visible par l'utilisateur
- Mettre a jour `CLAUDE.md` si la contribution change l'architecture ou les conventions

## Signaler un bug

Ouvrir une issue sur [GitHub](https://github.com/Greite/unraid-tui/issues) avec :

- La version (`unraid-tui version`)
- L'OS et l'architecture
- Les etapes pour reproduire
- Le comportement attendu vs observe
