# Docker

La page Docker affiche l'ensemble des containers Docker du serveur Unraid dans un tableau interactif.

## Tableau des containers

```
  Containers (5)  3 running

  ┌──────────────┬─────────────────────┬────────────┬──────────────────┬─────────────┐
  │ NAME         │ IMAGE               │ STATE      │ STATUS           │ PORTS       │
  ├──────────────┼─────────────────────┼────────────┼──────────────────┼─────────────┤
  │ plex         │ plexinc/pms:latest  │ ● running  │ Up 14 days       │ 32400:32400 │
  │ nextcloud    │ nextcloud:28        │ ● running  │ Up 14 days       │ 443:443     │
  │ homeassistant│ ghcr.io/ha:latest   │ ● running  │ Up 14 days       │ 8123:8123   │
  │ pihole       │ pihole/pihole:lat.. │ ○ exited   │ Exited (0) 2d    │ -           │
  │ test-db      │ postgres:15         │ ○ exited   │ Exited (1) 5d    │ -           │
  └──────────────┴─────────────────────┴────────────┴──────────────────┴─────────────┘

  ↑/↓: naviguer  │  r: rafraîchir
```

### Colonnes

| Colonne | Description                                           |
|---------|-------------------------------------------------------|
| NAME    | Nom du container                                      |
| IMAGE   | Image Docker (tronquée à 25 caractères si nécessaire) |
| STATE   | État avec indicateur visuel (● running, ○ exited, ◑ paused) |
| STATUS  | Détail du statut (ex: "Up 14 days", "Exited (0) 2 days ago") |
| PORTS   | Mapping de ports host:container (ou `-` si aucun)     |

### En-tête

- **Nombre total** de containers entre parenthèses
- **Nombre de containers running** affiché à côté

## Navigation

| Touche | Action                        |
|--------|-------------------------------|
| `↑`/`↓`| Naviguer dans le tableau      |
| `r`    | Rafraîchir la liste           |

Le tableau supporte le scroll si la liste dépasse la hauteur disponible.

## Colonnes responsives

Les largeurs des colonnes s'adaptent automatiquement à la taille du terminal. Chaque colonne utilise un pourcentage de la largeur disponible :

- NAME : 20%
- IMAGE : 25%
- STATE : 10%
- STATUS : 25%
- PORTS : 15%

## États

| État         | Affichage                              |
|--------------|----------------------------------------|
| Chargement   | Spinner animé + "Chargement des containers..." |
| Données OK   | Tableau complet avec compteur          |
| Erreur       | Bandeau rouge avec le message d'erreur |

## Requête GraphQL utilisée

```graphql
query {
  docker {
    containers {
      id name image state status
      ports { privatePort publicPort type }
      networks { networkId }
    }
  }
}
```

## Fichiers concernés

- `internal/tui/docker/docker.go` — Modèle Bubbletea, table Bubbles, formatage
- `internal/tui/docker/docker_test.go` — Tests unitaires (9 tests)
- `internal/api/queries.go` — Requête GraphQL (`queryContainers`)
