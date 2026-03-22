# unraid-tui

Interface terminal (TUI) pour monitorer et gerer un serveur [Unraid](https://unraid.net/) depuis la ligne de commande, via l'API GraphQL.

```
╭─────────────────── UNRAID CLI ────────────────────╮
│  [Dashboard]  Docker                               │
├────────────────────────────────────────────────────┤
│  ╭── CPU ──────────────╮  ╭── Memory ────────────╮ │
│  │ AMD Ryzen 9 5900X   │  │ Used: 24.3 / 64.0 GB│ │
│  │ Cores: 12 / 24T     │  │ ████████░░░░ 38.0%   │ │
│  │ Usage: ███████░ 72%  │  │                      │ │
│  ╰─────────────────────╯  ╰──────────────────────╯ │
│  ╭── System ───────────────────────────────────────╮│
│  │ Hostname: tower  │ OS: Unraid 6.12.6           ││
│  │ Kernel: 6.1.64   │ Board: ASRock X570 Taichi   ││
│  ╰─────────────────────────────────────────────────╯│
├────────────────────────────────────────────────────┤
│  tab changer page  │  1/2 aller a  │  q quitter   │
╰────────────────────────────────────────────────────╯
```

## Fonctionnalites

- **Dashboard** — CPU, memoire, infos systeme avec rafraichissement automatique (3s)
- **Docker** — Tableau interactif de tous les containers (nom, image, etat, ports)
- **Onboarding** — Assistant de configuration guide au premier lancement
- **Navigation clavier** — Changement de page par `Tab` ou touches numeriques

## Prerequis

- **Unraid 7.2+** (API GraphQL integree) ou Unraid 6.x avec le plugin [Unraid Connect](https://docs.unraid.net/API/)
- Une **cle API** Unraid (l'assistant de configuration vous guidera)

## Installation

### Homebrew (macOS)

```bash
brew install Greite/tap/unraid-tui
```

Mise a jour :

```bash
brew upgrade unraid-tui
```

### Binaires pre-compiles

Telecharger le binaire correspondant a votre OS/architecture depuis la page [Releases](https://github.com/Greite/unraid-tui/releases) :

| OS      | Architecture | Fichier                                |
|---------|-------------|----------------------------------------|
| macOS   | Apple Silicon| `unraid-tui_x.x.x_darwin_arm64.tar.gz`|
| macOS   | Intel       | `unraid-tui_x.x.x_darwin_amd64.tar.gz`|
| Linux   | x86_64      | `unraid-tui_x.x.x_linux_amd64.tar.gz` |
| Linux   | ARM64       | `unraid-tui_x.x.x_linux_arm64.tar.gz` |
| Windows | x86_64      | `unraid-tui_x.x.x_windows_amd64.zip`  |

```bash
# Exemple macOS Apple Silicon
tar xzf unraid-tui_*_darwin_arm64.tar.gz
sudo mv unraid-tui /usr/local/bin/
```

### Go install

Necessite Go 1.22+ :

```bash
go install github.com/Greite/unraid-tui@latest
```

### Depuis les sources

```bash
git clone https://github.com/Greite/unraid-tui.git
cd unraid-tui
make install
```

Cela compile le binaire et le copie dans `/usr/local/bin/`.

Pour desinstaller :

```bash
make uninstall
```

## Demarrage rapide

```bash
unraid-tui
```

Au premier lancement, un assistant interactif vous guide pour :

1. Saisir l'adresse de votre serveur Unraid
2. Tester la connexion
3. Creer et configurer votre cle API
4. Sauvegarder la configuration

La configuration est enregistree dans `~/.unraid-tui.yaml`.

### Configuration manuelle

Si vous preferez configurer manuellement, creez le fichier `~/.unraid-tui.yaml` :

```yaml
server_url: "http://192.168.1.100:3001"
api_key: "votre-cle-api"
```

Ou via variables d'environnement :

```bash
export UNRAID_SERVER_URL="http://192.168.1.100:3001"
export UNRAID_API_KEY="votre-cle-api"
```

### Obtenir une cle API

1. Ouvrir l'interface web Unraid
2. **Settings > Management Access > Developer Options**
3. Ouvrir Apollo GraphQL Studio
4. Executer :

```graphql
mutation {
  apiKey {
    create(input: {
      name: "unraid-tui"
      roles: [ADMIN]
    }) { key }
  }
}
```

5. Copier la cle retournee

## Utilisation

### Raccourcis clavier

| Touche       | Action                         |
|--------------|--------------------------------|
| `Tab`        | Page suivante                  |
| `Shift+Tab`  | Page precedente                |
| `1`          | Dashboard                      |
| `2`          | Docker                         |
| `r`          | Rafraichir (page Docker)       |
| `↑` / `↓`   | Naviguer dans les tableaux     |
| `q`          | Quitter                        |
| `Ctrl+C`     | Quitter                        |

### Pages

#### Dashboard

Affiche en temps reel :
- **CPU** — modele, nombre de coeurs, frequence, utilisation (%)
- **Memoire** — utilisation avec barre de progression
- **Systeme** — hostname, OS, kernel, carte mere

Les metriques se rafraichissent automatiquement toutes les 3 secondes.

#### Docker

Tableau interactif listant tous les containers :

| Colonne | Description |
|---------|-------------|
| NAME    | Nom du container |
| IMAGE   | Image Docker |
| STATE   | Etat avec indicateur (● running, ○ exited, ◑ paused) |
| STATUS  | Detail ("Up 14 days", "Exited (0) 2 days ago") |
| PORTS   | Mapping host:container |

## Developpement

### Commandes

```bash
make build         # Compiler le binaire
make test          # Lancer les tests
make test-verbose  # Tests avec detail
make test-cover    # Tests avec couverture HTML
make lint          # go vet
make run           # Build + executer
make clean         # Nettoyer
```

### Architecture

```
cmd/                     → Point d'entree Cobra
internal/api/            → Client GraphQL (interface + HTTP)
internal/config/         → Configuration Viper
internal/model/          → Types domaine
internal/tui/            → App Bubbletea (routeur)
internal/tui/common/     → Styles, messages, helpers
internal/tui/dashboard/  → Page dashboard
internal/tui/docker/     → Page Docker
internal/tui/onboarding/ → Assistant de configuration
```

### Tests

58 tests couvrent l'ensemble du projet :

| Package      | Tests | Ce qui est teste |
|--------------|-------|------------------|
| `api`        | 6     | Client HTTP, parsing, erreurs auth/connexion |
| `config`     | 5     | Chargement fichier, env vars, sauvegarde |
| `tui`        | 7     | Navigation pages, quit, resize |
| `dashboard`  | 7     | Panneaux CPU/memoire, loading, erreurs |
| `docker`     | 9     | Table containers, formatage, helpers |
| `onboarding` | 22    | Chaque transition d'etape, validation, normalisation |

```bash
make test
```

### Release

Les releases sont gerees par [GoReleaser](https://goreleaser.com/). A chaque tag Git, GoReleaser :

1. Compile les binaires pour macOS, Linux et Windows (amd64 + arm64)
2. Cree les archives `.tar.gz` / `.zip`
3. Publie la release GitHub
4. Met a jour la formule Homebrew dans [Greite/homebrew-tap](https://github.com/Greite/homebrew-tap)

```bash
# Tester la release localement (sans publier)
make release-dry

# Publier une release
git tag v0.1.0
git push origin v0.1.0
goreleaser release --clean
```

### Stack technique

| Composant | Librairie |
|-----------|-----------|
| TUI       | [Bubbletea v2](https://github.com/charmbracelet/bubbletea) |
| Styling   | [Lipgloss v2](https://github.com/charmbracelet/lipgloss) |
| Composants| [Bubbles v2](https://github.com/charmbracelet/bubbles) (table, spinner, textinput) |
| CLI       | [Cobra](https://github.com/spf13/cobra) |
| Config    | [Viper](https://github.com/spf13/viper) |
| API       | [Unraid GraphQL API](https://docs.unraid.net/API/) via `net/http` |

## Documentation

La documentation detaillee de chaque fonctionnalite est dans le dossier [`docs/`](docs/) :

- [Configuration](docs/configuration.md)
- [Onboarding](docs/onboarding.md)
- [Dashboard](docs/dashboard.md)
- [Docker](docs/docker.md)
- [Navigation](docs/navigation.md)
- [Client API](docs/api-client.md)

## Licence

MIT
