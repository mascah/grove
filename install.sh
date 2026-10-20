#!/bin/sh
# Installs grove from a GitHub release.
#   curl -fsSL https://raw.githubusercontent.com/mascah/grove/main/install.sh | sh
#   curl -fsSL https://raw.githubusercontent.com/mascah/grove/main/install.sh | sh -s -- v0.1.0
set -eu

releases=https://github.com/mascah/grove/releases
dir=${GROVE_INSTALL_DIR:-$HOME/.local/bin}

fail() {
	echo "grove: $1" >&2
	exit 1
}

case $(uname -s) in
Darwin) os=darwin ;;
Linux) os=linux ;;
*) fail "no build for $(uname -s)." ;;
esac
case $(uname -m) in
x86_64 | amd64) arch=amd64 ;;
arm64 | aarch64) arch=arm64 ;;
*) fail "no build for $(uname -m)." ;;
esac

if [ $# -gt 0 ]; then
	url=$releases/download/$1
else
	url=$releases/latest/download
fi
archive=grove_${os}_${arch}.tar.gz

if command -v curl >/dev/null 2>&1; then
	fetch() { curl -fsSL -o "$2" "$1"; }
elif command -v wget >/dev/null 2>&1; then
	fetch() { wget -q -O "$2" "$1"; }
else
	fail "needs curl or wget."
fi
if command -v sha256sum >/dev/null 2>&1; then
	sha256() { sha256sum "$1"; }
else
	sha256() { shasum -a 256 "$1"; }
fi

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

fetch "$url/$archive" "$tmp/$archive" || fail "could not download $url/$archive."
fetch "$url/checksums.txt" "$tmp/checksums.txt" || fail "could not download $url/checksums.txt."

want=$(awk -v f="$archive" '$2 == f { print $1 }' "$tmp/checksums.txt")
got=$(sha256 "$tmp/$archive" | awk '{ print $1 }')
if [ -z "$want" ] || [ "$want" != "$got" ]; then
	fail "$archive does not match checksums.txt; nothing was installed."
fi

tar -xzf "$tmp/$archive" -C "$tmp" grove
mkdir -p "$dir"
mv "$tmp/grove" "$dir/grove"
echo "Installed grove in $dir."

case ":$PATH:" in
*":$dir:"*) ;;
*) echo "$dir is not on your PATH; add it to run grove." ;;
esac
