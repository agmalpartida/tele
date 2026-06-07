#!/usr/bin/env bash
set -euo pipefail

# Updates Formula/tele.rb in the homebrew-tele tap repo with SHA256 checksums
# for a new release.
#
# Usage:
#   ./scripts/update-homebrew-formula.sh <version>
#
# Example:
#   ./scripts/update-homebrew-formula.sh 1.0.9
#
# The script downloads tarballs from the GitHub release, computes their
# SHA256, and rewrites the formula file in-place.

HOMEBREW_REPO="$HOME/Git/personal/homebrew-tele"
FORMULA="$HOMEBREW_REPO/Formula/tele.rb"

if [ ! -f "$FORMULA" ]; then
  echo "error: formula not found at $FORMULA" >&2
  exit 1
fi

VERSION="${1:?usage: $0 <version>}"
# strip leading "v" if present
VERSION="${VERSION#v}"

REPO="agmalpartida/tele"
BASE_URL="https://github.com/$REPO/releases/download/v$VERSION"

# platform -> tarball name pattern (goreleaser name_template: "tele_{{ .Os }}_{{ .Arch }}")
declare -A PLATFORMS=(
  [darwin_amd64]="tele_darwin_amd64.tar.gz"
  [darwin_arm64]="tele_darwin_arm64.tar.gz"
  [linux_amd64]="tele_linux_amd64.tar.gz"
  [linux_arm64]="tele_linux_arm64.tar.gz"
)

# Choose sha256 command
if command -v sha256sum &>/dev/null; then
  SHA_CMD="sha256sum"
elif command -v shasum &>/dev/null; then
  SHA_CMD="shasum -a 256"
else
  echo "error: no sha256 command found" >&2
  exit 1
fi

TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

echo ":: downloading tarballs for v$VERSION ..."
declare -A CHECKS
for plat in "${!PLATFORMS[@]}"; do
  tarball="${PLATFORMS[$plat]}"
  url="$BASE_URL/$tarball"
  echo "   $tarball"
  curl -sL -o "$TMPDIR/$tarball" "$url"
  sum="$($SHA_CMD "$TMPDIR/$tarball" | awk '{print $1}')"
  CHECKS[$plat]="$sum"
done

echo ""
echo ":: updating $FORMULA ..."

# Update version line
sed -i '' 's/version ".*"/version "'"$VERSION"'"/' "$FORMULA"

# Update sha256 lines by matching the URL that precedes each one.
# Each block has the form:
#   url "...tele_{os}_{arch}.tar.gz"
#   sha256 "hex..."
for plat in "${!PLATFORMS[@]}"; do
  tarball="${PLATFORMS[$plat]}"
  checksum="${CHECKS[$plat]}"
  # Match the url line for this tarball and replace sha256 on the next line
  sed -i '' -E '/url ".*'"$tarball"'"/{n;s/sha256 "[a-f0-9]*"/sha256 "'"$checksum"'"/;}' "$FORMULA"
done

echo "done. updated to v$VERSION"
echo ""
echo "review changes:"
echo "  git -C $HOMEBREW_REPO diff"
