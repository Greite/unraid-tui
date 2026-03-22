# Client API

Le client API encapsule toute la communication avec l'API GraphQL d'Unraid.

## Architecture

```
UnraidClient (interface)
    │
    ├── httpClient (implémentation réelle)
    │   └── envoie des POST HTTP à /graphql
    │
    └── MockClient (tests)
        └── fonctions configurables par test
```

### Interface

```go
type UnraidClient interface {
    GetSystemInfo(ctx context.Context) (*model.SystemInfo, error)
    GetSystemMetrics(ctx context.Context) (*model.SystemMetrics, error)
    GetContainers(ctx context.Context) ([]model.Container, error)
}
```

Tout le code TUI dépend de cette interface, jamais de l'implémentation concrète.

## Implémentation HTTP

### Requêtes

Chaque méthode envoie une requête POST à `<server_url>/graphql` avec :

- **Header** `Content-Type: application/json`
- **Header** `Authorization: Bearer <api_key>`
- **Body** : `{"query": "<graphql_query>"}`

### Réponses

Les réponses GraphQL sont décodées en structs Go spécifiques (package `api/types.go`), puis converties en types domaine (`model/model.go`) via des méthodes `toDomain()`.

### Gestion d'erreurs

| Situation                | Erreur retournée            |
|--------------------------|-----------------------------|
| Serveur injoignable      | `ErrConnectionFailed`       |
| Réponse 401 ou 403      | `ErrUnauthorized`           |
| Code HTTP inattendu      | `unexpected status <code>`  |
| Erreur dans le JSON GraphQL | `graphql: <message>`     |
| JSON invalide            | `decoding response: <err>`  |

Les erreurs sentinelles `ErrConnectionFailed` et `ErrUnauthorized` permettent une gestion spécifique côté TUI (affichage adapté, pas de retry sur erreur d'auth).

## Mock pour les tests

`api/mock.go` fournit un `MockClient` avec des champs de fonction remplaçables :

```go
mock := &api.MockClient{
    GetContainersFn: func(ctx context.Context) ([]model.Container, error) {
        return []model.Container{{Name: "test"}}, nil
    },
}
```

Si une fonction n'est pas définie, le mock retourne `(nil, nil)`.

## Requêtes GraphQL

Les requêtes sont des constantes string dans `api/queries.go`. Elles correspondent à l'API GraphQL Unraid :

| Constante            | Données récupérées                          |
|----------------------|---------------------------------------------|
| `querySystemInfo`    | CPU, mémoire, OS, baseboard                 |
| `querySystemMetrics` | Usage CPU %, usage mémoire %                |
| `queryContainers`    | Liste des containers avec ports et réseaux  |

## Tests

6 tests couvrent le client HTTP :

- `TestGetSystemInfo_Success` — Vérifie headers, parsing, conversion domaine
- `TestGetSystemMetrics_Success` — Parsing des métriques
- `TestGetContainers_Success` — Parsing multi-containers avec ports
- `TestGetSystemInfo_Unauthorized` — Détection du 401
- `TestGetSystemInfo_GraphQLError` — Erreur dans la réponse GraphQL
- `TestGetSystemInfo_ConnectionError` — Serveur injoignable

## Fichiers concernés

- `internal/api/client.go` — Interface + implémentation HTTP
- `internal/api/client_test.go` — Tests avec httptest
- `internal/api/queries.go` — Requêtes GraphQL
- `internal/api/types.go` — Structs réponses JSON + conversion domaine
- `internal/api/mock.go` — Mock client pour tests TUI
