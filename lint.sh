#!/usr/bin/env bash
set -euo pipefail

CMD="go vet github.com/variadico/noti/..."
if [[ ! -z "$($CMD 2>&1)" ]]; then
	echo "Fix go vet errors"
	echo "For more info, run: $CMD"
	exit 1
fi

CMD="golint github.com/variadico/noti/..."
if [[ ! -z "$($CMD 2>&1)" ]]; then
	echo "Fix golint warnings"
	echo "For more info, run: $CMD"
	exit 1
fi

CMD="gofmt -d ."
if [[ ! -z "$($CMD)" ]]; then
	echo "Format your code"
	echo "For more info, run: $CMD"
	exit 1
fi

CMD="gofmt -s -d ."
if [[ ! -z "$($CMD)" ]]; then
	echo "Simplify your code"
	echo "For more info, run: $CMD"
	exit 1
fi
