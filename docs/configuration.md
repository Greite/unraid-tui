# Configuration

## Fichier de configuration

L'application cherche un fichier `~/.unraid-tui/config.yaml` dans le répertoire home de l'utilisateur.

### Format

```yaml
server_url: "http://192.168.1.100:3001"
api_key: "votre-clé-api-unraid"
```

### Paramètres

| Paramètre    | Requis | Description                                      |
|--------------|--------|--------------------------------------------------|
| `server_url` | Oui    | URL du serveur Unraid (incluant le port)         |
| `api_key`    | Oui    | Clé API pour l'authentification Bearer           |

## Variables d'environnement

Chaque paramètre peut être surchargé par une variable d'environnement préfixée `UNRAID_` :

| Variable              | Surcharge      |
|-----------------------|----------------|
| `UNRAID_SERVER_URL`   | `server_url`   |
| `UNRAID_API_KEY`      | `api_key`      |

Les variables d'environnement ont priorité sur le fichier de configuration.

### Exemple

```bash
export UNRAID_SERVER_URL="http://10.0.0.5:3001"
export UNRAID_API_KEY="ma-clé-secrète"
unraid-tui
```

## Obtenir une clé API

1. Ouvrir l'interface web Unraid
2. Aller dans **Settings → Management Access → Developer Options**
3. Ouvrir Apollo GraphQL Studio
4. Créer une clé API via la mutation GraphQL :

```graphql
mutation {
  apiKey {
    create(input: {
      name: "unraid-tui"
      description: "CLI monitoring tool"
      roles: [VIEWER]
      permissions: [
        { resource: INFO, actions: [READ_ANY] }
        { resource: DOCKER, actions: [READ_ANY] }
      ]
    }) {
      key
    }
  }
}
```

5. Copier la clé retournée dans `~/.unraid-tui/config.yaml`

## Validation

Au lancement, l'application vérifie que `server_url` et `api_key` sont définis. Si l'un des deux manque, un message d'erreur indique la source attendue (fichier ou variable d'environnement).

## Fichiers concernés

- `internal/config/config.go` — Chargement et validation
- `internal/config/config_test.go` — Tests (fichier, env vars, override)
