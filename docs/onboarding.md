# Onboarding

L'onboarding est un assistant de configuration interactif qui se lance automatiquement au premier demarrage du CLI, quand aucune configuration n'est detectee.

## Declenchement

L'onboarding se lance si :
- Le fichier `~/.unraid-tui/config.yaml` n'existe pas
- **Et** les variables d'environnement `UNRAID_SERVER_URL` / `UNRAID_API_KEY` ne sont pas definies

Si la config existe (fichier ou env vars), le dashboard se lance directement.

## Etapes

### 1. Ecran d'accueil

```
╭────────────────────────────────────────────────────────╮
│  Bienvenue !                                           │
│                                                        │
│  Cet assistant va vous aider a configurer la connexion │
│  a votre serveur Unraid en quelques etapes :           │
│                                                        │
│    1. Saisir l'adresse de votre serveur                │
│    2. Tester la connexion                              │
│    3. Configurer votre cle API                         │
│    4. Sauvegarder la configuration                     │
│                                                        │
│  Le fichier sera sauvegarde dans ~/.unraid-tui/config.yaml    │
╰────────────────────────────────────────────────────────╯

  enter commencer  esc retour
```

Presente le processus et ses etapes.

### 2. Adresse du serveur (Etape 1/3)

L'utilisateur saisit l'URL de son serveur Unraid. L'input accepte plusieurs formats :

| Saisie                        | Normalise en                  |
|-------------------------------|-------------------------------|
| `192.168.1.100:3001`          | `http://192.168.1.100:3001`   |
| `http://tower:3001`           | `http://tower:3001`           |
| `https://secure.local:3001/`  | `https://secure.local:3001`   |

Normalisation automatique :
- Ajout de `http://` si aucun schema n'est present
- Suppression du `/` final

Validation : l'URL ne peut pas etre vide.

### 3. Test de connexion

Envoie une requete GraphQL minimale (`{ __typename }`) au serveur. Le test verifie seulement que le serveur est joignable (meme une reponse 401 est consideree comme un succes a cette etape — cela signifie que l'API est bien la).

- **Succes** : passage a l'etape suivante
- **Echec** : retour a la saisie de l'URL avec le message d'erreur
- **Timeout** : 5 secondes

### 4. Instructions pour la cle API (Etape 2/3)

Affiche les instructions detaillees pour creer une cle API depuis l'interface web Unraid :

1. Ouvrir l'interface web Unraid
2. Settings > Management Access > Developer Options
3. Ouvrir Apollo GraphQL Studio
4. Executer la mutation `apiKey.create`
5. Copier la cle retournee

La mutation GraphQL exacte est affichee dans l'interface.

### 5. Saisie de la cle API (Etape 3/3)

Champ de saisie masque (mode password — les caracteres sont remplaces par `*`).

Validation : la cle ne peut pas etre vide.

### 6. Verification de la cle API

Envoie une requete authentifiee (`info { os { hostname } }`) pour verifier que la cle fonctionne.

- **200 OK** : la cle est valide, passage a la sauvegarde
- **401/403** : cle invalide ou permissions insuffisantes, retour a la saisie
- **Autre** : erreur affichee, retour a la saisie

### 7. Sauvegarde

Ecrit le fichier `~/.unraid-tui/config.yaml` avec les permissions `0600` (lecture/ecriture proprietaire uniquement).

Format du fichier :
```yaml
server_url: "http://192.168.1.100:3001"
api_key: "votre-cle-api"
```

### 8. Ecran de confirmation

```
╭────────────────────────────────────────────────────────╮
│  Configuration terminee !                              │
│                                                        │
│  Votre configuration a ete sauvegardee dans :          │
│    ~/.unraid-tui/config.yaml                                  │
│                                                        │
│  Serveur : http://192.168.1.100:3001                   │
│  Cle API : ********** (sauvegardee)                    │
│                                                        │
│  Le dashboard va maintenant se lancer.                 │
╰────────────────────────────────────────────────────────╯

  enter lancer le dashboard
```

Apres confirmation, le dashboard principal se lance automatiquement.

## Navigation

| Touche   | Action                                  |
|----------|-----------------------------------------|
| `enter`  | Valider / Passer a l'etape suivante     |
| `esc`    | Revenir a l'etape precedente            |
| `Ctrl+C` | Annuler et quitter                     |

## Barre de progression

Une barre de progression en haut de l'ecran indique l'avancement :

```
  ● Serveur  —  ◉ Connexion  —  ○ Cle API  —  ○ Termine
```

- `●` etape terminee (vert)
- `◉` etape en cours (violet)
- `○` etape a venir (gris)

## Gestion des erreurs

Les erreurs sont affichees en rouge sous le contenu de l'etape. L'utilisateur reste sur l'etape en cours et peut corriger sa saisie.

Erreurs possibles :
- URL vide
- Serveur injoignable (timeout, DNS, connexion refusee)
- URL invalide
- Cle API vide
- Cle API invalide (401/403)
- Erreur de sauvegarde du fichier

## Fichiers concernes

- `internal/tui/onboarding/onboarding.go` — Modele Bubbletea multi-etapes
- `internal/tui/onboarding/onboarding_test.go` — 22 tests unitaires
- `internal/config/config.go` — `Exists()`, `Save()`, `FilePath()`
- `cmd/root.go` — Detection et lancement de l'onboarding
