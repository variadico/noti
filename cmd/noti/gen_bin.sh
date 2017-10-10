#!/usr/bin/env bash
set -ex

VERSION="2.6.0"
rm -f noti
rm -f noti.exe

GOOS=darwin GOARCH=amd64 go build
tar -czf "noti$VERSION.darwin-amd64.tar.gz" noti
rm -f noti

GOOS=linux GOARCH=amd64 go build
tar -czf "noti$VERSION.linux-amd64.tar.gz" noti
rm -f noti

GOOS=windows GOARCH=amd64 go build
tar -czf "noti$VERSION.windows-amd64.tar.gz" noti.exe
rm -f noti.exe
