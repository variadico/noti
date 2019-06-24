branch = $(shell git rev-parse --abbrev-ref HEAD)
tag = $(shell git describe --abbrev=0 --tags)
rev = $(shell git rev-parse --short HEAD)

export GOFLAGS = -mod=vendor

pkgs = $(shell go list ./... | grep -v /vendor/)
project_src = $(shell find . -name "*.go" | grep -v /vendor/ | grep -v _test.go)

bin_prefix = "vendor/.bin"

# Build
cmd/noti/noti: $(project_src)
	go build -race -o $@ \
		-ldflags "-X github.com/variadico/noti/internal/command.Version=$(branch)-$(rev)" \
		github.com/variadico/noti/cmd/noti
build: cmd/noti/noti

# Manual
docs/man/noti.1: docs/man/noti.1.md
	pandoc -s -t man $< -o $@
docs/man/noti.yaml.5: docs/man/noti.yaml.5.md
	pandoc -s -t man $< -o $@
man: docs/man/noti.1 docs/man/noti.yaml.5

# Release
release/noti$(tag).linux-amd64.tar.gz: $(shell find . -name "*.go")
	mkdir -p $(@D)
	GOOS=linux GOARCH=amd64 \
		go build \
		-ldflags "-s -w -X github.com/variadico/noti/internal/command.Version=$(tag)" \
		github.com/variadico/noti/cmd/noti
	tar -czf release/noti$(tag).linux-amd64.tar.gz noti
	rm -f noti
release/noti$(tag).darwin-amd64.tar.gz: $(shell find . -name "*.go")
	mkdir -p $(@D)
	GOOS=darwin GOARCH=amd64 \
		go build \
		-ldflags "-s -w -X github.com/variadico/noti/internal/command.Version=$(tag)" \
		github.com/variadico/noti/cmd/noti
	tar -czf release/noti$(tag).darwin-amd64.tar.gz noti
	rm -f noti
release/noti$(tag).windows-amd64.tar.gz: $(shell find . -name "*.go")
	mkdir -p $(@D)
	GOOS=windows GOARCH=amd64 \
		go build \
		-ldflags "-s -w -X github.com/variadico/noti/internal/command.Version=$(tag)" \
		github.com/variadico/noti/cmd/noti
	tar -czf release/noti$(tag).windows-amd64.tar.gz noti.exe
	rm -f noti.exe
release: release/noti$(tag).linux-amd64.tar.gz release/noti$(tag).darwin-amd64.tar.gz \
	release/noti$(tag).windows-amd64.tar.gz

.PHONY: install
install:
	go install \
		-ldflags "-X github.com/variadico/noti/internal/command.Version=$(branch)-$(rev)" \
		github.com/variadico/noti/cmd/noti

.PHONY: test
test:
	go test -v -cover -race ./...

.PHONY: vendor
vendor:
	go mod vendor
	go mod tidy

PHONY: clean
clean:
	go clean
	rm -f cmd/noti/noti
	rm -rf vendor/.bin
	rm -rf release/
	git clean -x -f -d
	git remote prune origin

