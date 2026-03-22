# Navigation

L'application utilise une interface TUI multi-pages avec navigation par clavier.

## Pages disponibles

| Page       | Touche directe | Description                       |
|------------|----------------|-----------------------------------|
| Dashboard  | `1`            | Monitoring CPU, mémoire, système  |
| Docker     | `2`            | Liste des containers Docker       |

## Raccourcis globaux

Ces raccourcis fonctionnent depuis n'importe quelle page :

| Touche       | Action                                        |
|--------------|-----------------------------------------------|
| `Tab`        | Page suivante (cyclique)                       |
| `Shift+Tab`  | Page précédente (cyclique)                     |
| `1`          | Aller au Dashboard                             |
| `2`          | Aller à la page Docker                         |
| `q`          | Quitter l'application                          |
| `Ctrl+C`     | Quitter l'application                          |

## Raccourcis par page

### Docker

| Touche | Action                   |
|--------|--------------------------|
| `↑`/`↓`| Naviguer dans le tableau |
| `r`    | Rafraîchir les containers |

## Interface

```
┌──────────────────────────────────────────────────────┐
│  UNRAID CLI    Dashboard   Docker                    │  ← Header avec tabs
├──────────────────────────────────────────────────────┤
│                                                      │
│  Contenu de la page active                           │  ← Zone de contenu
│                                                      │
├──────────────────────────────────────────────────────┤
│  tab changer page  │  1/2 aller à  │  q quitter     │  ← Footer avec aide
└──────────────────────────────────────────────────────┘
```

### Header

- Titre "UNRAID CLI" en violet
- Onglets des pages avec l'onglet actif en surbrillance

### Footer

- Rappel des raccourcis clavier principaux

## Comportement au changement de page

- Quand on navigue vers une page, sa commande `Init()` est ré-exécutée pour rafraîchir les données.
- Les données de la page précédente restent en cache (pas de rechargement inutile au retour).

## Fichiers concernés

- `internal/tui/app.go` — Routeur de pages, gestion des touches globales
- `internal/tui/app_test.go` — Tests de navigation (7 tests)
- `internal/tui/header.go` — Rendu de l'en-tête avec onglets
- `internal/tui/footer.go` — Rendu du pied de page avec raccourcis
