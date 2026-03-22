# Dashboard

Le dashboard est la page d'accueil de l'application. Il affiche les informations système du serveur Unraid en temps réel.

## Panneaux

### CPU

Affiche les informations du processeur et son utilisation en temps réel.

```
╭── CPU ──────────────────────╮
│  AMD Ryzen 9 5900X          │
│  Cores: 16  @ 3.7 GHz      │
│  ████████████░░░░░░░░ 62.5% │
╰─────────────────────────────╯
```

- **Brand** : Modèle du processeur
- **Cores** : Nombre de cœurs physiques
- **Speed** : Fréquence en GHz
- **Barre de progression** : Utilisation CPU en temps réel (%)

### Memory

Affiche l'utilisation mémoire du serveur.

```
╭── Memory ───────────────────╮
│  Used: 32.0 GB / 64.0 GB   │
│  ██████████░░░░░░░░░░ 50.0% │
╰─────────────────────────────╯
```

- **Used / Total** : Mémoire utilisée sur la mémoire totale (format humain)
- **Barre de progression** : Pourcentage d'utilisation

### System

Affiche les informations générales du système.

```
╭── System ───────────────────────────────────────────╮
│  Hostname: tower  │  OS: Unraid 6.12.6              │
│  Kernel: 6.1.64   │  Platform: linux                │
│  Board: ASRock X570 Taichi                          │
╰─────────────────────────────────────────────────────╯
```

- **Hostname** : Nom du serveur
- **OS** : Distribution et version
- **Kernel** : Version du noyau Linux
- **Platform** : Architecture
- **Board** : Carte mère (fabricant + modèle)

## Rafraîchissement automatique

Les métriques (CPU et mémoire) sont rafraîchies automatiquement toutes les **3 secondes** via un mécanisme de polling. Les informations système statiques (CPU model, OS, baseboard) ne sont chargées qu'une seule fois au premier affichage.

## États

| État         | Affichage                              |
|--------------|----------------------------------------|
| Chargement   | Spinner animé + "Chargement du dashboard..." |
| Données OK   | Les 3 panneaux CPU, Memory, System     |
| Erreur       | Bandeau rouge avec le message d'erreur + les données en cache restent visibles |

## Requêtes GraphQL utilisées

- `info` — Informations système (CPU, mémoire, OS, baseboard)
- `metrics` — Métriques temps réel (usage CPU, usage mémoire)

## Fichiers concernés

- `internal/tui/dashboard/dashboard.go` — Modèle Bubbletea, rendering des panneaux
- `internal/tui/dashboard/dashboard_test.go` — Tests unitaires
- `internal/api/queries.go` — Requêtes GraphQL (`querySystemInfo`, `querySystemMetrics`)
