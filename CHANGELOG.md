# Changelog

Toutes les modifications notables de ce projet sont documentées ici.
Le format suit [Keep a Changelog](https://keepachangelog.com/fr/1.1.0/) et le projet adhère au [versionnement sémantique](https://semver.org/lang/fr/).

## [Unreleased]

## [1.0.0] - 2026-11-24

### Ajouté
- Réécriture complète en Go : binaire unique, sans dépendance à l'exécution.
- Détection automatique de la dernière version stable de PowerShell et de l'architecture (x64, arm64, arm32).
- Vérification du SHA-256 de l'archive via `hashes.sha256`.
- Relance automatique avec `sudo` si l'outil n'est pas lancé en root.
- Options `--check`, `--force` et `--version`.
- `Makefile` (build, install, uninstall, run, check, fmt, vet, test, clean).
- Script `install.sh` pour une installation en une commande `curl`.
- Workflows GitHub Actions : CI (format, vet, tests, build) et Release (GoReleaser sur les tags `v*`).
- Badges, instructions d'installation et section Release dans le README.

### Modifié
- `make install` compile, copie le binaire dans `/usr/local/bin` puis retire le binaire local.
- `license` renommé en `LICENSE`.

### Supprimé
- Script `pwsh-i.sh` et ancien `makefile`.

## [0.1.0]

### Ajouté
- Script bash `pwsh-i.sh` d'installation de PowerShell sur Linux, lancé via `make i v=<version>`.

[Unreleased]: https://github.com/mzeahmed/pwsh-updater/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/mzeahmed/pwsh-updater/releases/tag/v1.0.0
