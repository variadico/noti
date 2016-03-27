#!/bin/bash
set -euo pipefail

if [[ ! -z "$(go vet 2>&1)" ]]; then
	echo "Fix go vet errors"
	echo "For more info, run: go vet"
fi

if [[ ! -z "$(golint 2>&1)" ]]; then
	echo "Fix golint warnings"
	echo "For more info, run: golint"
fi

if [[ ! -z "$(gofmt -d ./*.go)" ]]; then
	echo "Format your code"
	echo "For more info, run: gofmt -d ./*.go"
	exit 1
fi

if [[ ! -z "$(gofmt -s -d ./*.go)" ]]; then
	echo "Simplify your code"
	echo "For more info, run: gofmt -s -d ./*.go"
fi
