export GOFLAGS := -mod=vendor
export GO111MODULE := on
export GOPROXY := off
export GOSUMDB := off

branch := $(shell git rev-parse --abbrev-ref HEAD)
tag := $(shell git describe --abbrev=0 --tags)
rev := $(shell git rev-parse --short HEAD)

golangciLint := ./tools/golangci-lint-1.39.0-$(shell go env GOOS)-amd64

goSrc := $(shell find cmd internal service -name "*.go")

goBin := $(strip $(shell go env GOBIN))
ifeq ($(goBin),)
goBin := $(shell go env GOPATH)/bin
endif

vendor: go.mod go.sum
	GOPROXY=direct go mod vendor
	touch $@

noti: $(goSrc) vendor
	go build -race -o $@ \
		-ldflags "-X github.com/variadico/noti/internal/command.Version=$(branch)-$(rev)" \
		github.com/variadico/noti/cmd/$@

release/noti.linuxrelease: $(gosrc) vendor
	mkdir --parents $$(dirname $@)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $@ \
		-ldflags "-s -w -X github.com/variadico/noti/internal/command.Version=$(tag)" \
		github.com/variadico/noti/cmd/noti
release/noti.darwinrelease: $(gosrc) vendor
	mkdir --parents $$(dirname $@)
	GOOS=darwin GOARCH=amd64 go build -o $@ \
		-ldflags "-s -w -X github.com/variadico/noti/internal/command.Version=$(tag)" \
		github.com/variadico/noti/cmd/noti
release/noti.windowsrelease: $(gosrc) vendor
	mkdir --parents $$(dirname $@)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $@ \
		-ldflags "-s -w -X github.com/variadico/noti/internal/command.Version=$(tag)" \
		github.com/variadico/noti/cmd/noti
release/noti$(tag).%-amd64.tar.gz: release/noti.%release
	tar czvf $@ --transform 's#$<#noti#g' $<

docs/man/noti.1: docs/man/noti.1.md
	pandoc -s -t man $< -o $@
docs/man/noti.yaml.5: docs/man/noti.yaml.5.md
	pandoc -s -t man $< -o $@

.PHONY: build
build: noti

.PHONY: install
install: noti
	mv $< $(goBin)

.PHONY: test
test:
	go test -v -cover -race ./internal/... ./service/...

.PHONY: test-integration
test-integration: install
	go test -v ./tests/...

.PHONY: lint
lint:
	go vet ./...
	$(golangciLint) run --no-config --exclude-use-default=false --max-same-issues=0 \
	--timeout 15s \
	--disable errcheck \
	--disable stylecheck \
	--enable bodyclose \
	--enable golint \
	--enable unconvert \
	--enable dupl \
	--enable gocyclo \
	--enable gofmt \
	--enable goimports \
	--enable misspell \
	--enable lll \
	--enable unparam \
	--enable nakedret \
	--enable prealloc \
	--enable exportloopref \
	--enable gocritic \
	--enable gochecknoinits \
	./...

.PHONY: man
man: docs/man/noti.1 docs/man/noti.yaml.5

.PHONY: release
release: release/noti$(tag).linux-amd64.tar.gz release/noti$(tag).windows-amd64.tar.gz

# Special because uses CGO.
.PHONY: release-darwin
release-darwin: release/noti$(tag).darwin-amd64.tar.gz
