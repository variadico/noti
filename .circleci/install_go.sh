#!/usr/bin/env bash

set -ex

# Save Go env in bash env.
echo 'export GOROOT=$HOME/go$GOVERSION' >> "$BASH_ENV"
echo 'export PATH=$GOROOT/bin:$HOME/go/bin:$PATH' >> "$BASH_ENV"
source "$BASH_ENV"

godl="${TMPDIR}go${GOVERSION}.tar.gz"
url="https://redirector.gvt1.com/edgedl/go/go$GOVERSION.darwin-amd64.tar.gz"

curl -L -o "$godl" "$url"

mkdir -p "$GOROOT"
# --strip 1 to remove the extra "go" directory.
tar -C "$GOROOT" -xzf "$godl" --strip 1

go version
go env
