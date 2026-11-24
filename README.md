# pwsh-updater

[![Go Version](https://img.shields.io/github/go-mod/go-version/mzeahmed/pwsh-updater)](go.mod)
[![License](https://img.shields.io/github/license/mzeahmed/pwsh-updater)](LICENSE)
[![Last commit](https://img.shields.io/github/last-commit/mzeahmed/pwsh-updater)](https://github.com/mzeahmed/pwsh-updater/commits/main)
[![Go Report Card](https://goreportcard.com/badge/github.com/mzeahmed/pwsh-updater)](https://goreportcard.com/report/github.com/mzeahmed/pwsh-updater)
![Platform](https://img.shields.io/badge/platform-linux-informational)

Installe ou met à jour PowerShell sur Linux depuis les releases GitHub officielles. Binaire Go unique, sans dépendance à l'exécution (ni `wget`, ni `tar`). `make` n'est utile que pour compiler.

## Installation

### Depuis les sources (recommandé)

```bash
git clone https://github.com/mzeahmed/pwsh-updater.git
cd pwsh-updater
make install
```

Le binaire est compilé puis copié dans `/usr/local/bin` (sudo requis).

### Avec `go install`

```bash
go install github.com/mzeahmed/pwsh-updater@latest
```

Le binaire est placé dans `$(go env GOPATH)/bin` (pensez à l'ajouter au `PATH`).

## Utilisation

```bash
pwsh-updater            # installe la dernière version stable
pwsh-updater 7.5.1      # installe une version précise
pwsh-updater --check    # vérifie si une mise à jour est disponible
pwsh-updater --force    # réinstalle même si la version est déjà présente
pwsh-updater --version  # version de l'outil
```

Si l'outil n'est pas lancé en root, il se relance automatiquement avec `sudo`.

## Ce que fait l'outil

- Détecte l'architecture (x64, arm64, arm32) et la dernière version si aucune n'est donnée
- Télécharge l'archive et vérifie son SHA-256 (`hashes.sha256` de la release)
- Remplace `/opt/microsoft/powershell/7` et met à jour le lien `/usr/bin/pwsh`
- Nettoie les fichiers temporaires

## Développement

Un `Makefile` regroupe les commandes utiles (`make help` pour la liste) :

| Commande | Action |
|---|---|
| `make build` | Compile le binaire `./pwsh-updater` |
| `make install` | Compile, installe dans `/usr/local/bin` (sudo, `BINDIR=` pour changer) puis retire le binaire local |
| `make uninstall` | Supprime le binaire de `/usr/local/bin` |
| `make run ARGS="--check"` | Compile puis lance avec des arguments |
| `make check` | Vérifie si une mise à jour de PowerShell existe |
| `make fmt` / `make vet` / `make test` | Format, analyse statique, tests |
| `make clean` | Supprime le binaire compilé |

La version embarquée vient de `git describe` (surchargeable : `make build VERSION=1.0.0`).
