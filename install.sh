#!/bin/sh
# Installe la dernière release de pwsh-updater dans /usr/local/bin.
#   curl -fsSL https://raw.githubusercontent.com/mzeahmed/pwsh-updater/main/install.sh | sh
# Variables : BINDIR (défaut /usr/local/bin), VERSION (défaut : dernière release, ex. v1.0.0)
set -eu

REPO="mzeahmed/pwsh-updater"
BINDIR="${BINDIR:-/usr/local/bin}"

[ "$(uname -s)" = "Linux" ] || { echo "Erreur : seul Linux est supporté." >&2; exit 1; }

case "$(uname -m)" in
  x86_64 | amd64) ARCH="x86_64" ;;
  aarch64 | arm64) ARCH="arm64" ;;
  armv7l | armv7) ARCH="armv7" ;;
  *) echo "Erreur : architecture non supportée : $(uname -m)" >&2; exit 1 ;;
esac

ARCHIVE="pwsh-updater_Linux_${ARCH}.tar.gz"
if [ -n "${VERSION:-}" ]; then
  BASE="https://github.com/${REPO}/releases/download/${VERSION}"
else
  BASE="https://github.com/${REPO}/releases/latest/download"
fi

TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

echo "Téléchargement de ${ARCHIVE}..."
curl -fsSL "${BASE}/${ARCHIVE}" -o "${TMP}/${ARCHIVE}"
curl -fsSL "${BASE}/checksums.txt" -o "${TMP}/checksums.txt"

echo "Vérification du checksum..."
(cd "$TMP" && grep " ${ARCHIVE}\$" checksums.txt | sha256sum -c -)

tar -xzf "${TMP}/${ARCHIVE}" -C "$TMP" pwsh-updater

SUDO=""
[ -w "$BINDIR" ] || SUDO="sudo"
$SUDO install -m 755 "${TMP}/pwsh-updater" "${BINDIR}/pwsh-updater"

echo "pwsh-updater installé dans ${BINDIR}/pwsh-updater"
"${BINDIR}/pwsh-updater" --version
