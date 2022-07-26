export GOFLAGS := -mod=vendor
export GO111MODULE := on
export GOPROXY := direct
export GOSUMDB := off

branch := $(shell git rev-parse --abbrev-ref HEAD)
tag := $(shell git describe --abbrev=0 --tags)
rev := $(shell git rev-parse --short HEAD)

golangci-lint := ./tools/golangci-lint-1.30.0-$(shell go env GOOS)-amd64

gosrc := $(shell find cmd internal service -name "*.go")

gobin := $(strip $(shell go env GOBIN))
ifeq ($(gobin),)
gobin := $(shell go env GOPATH)/bin
endif

ldflags := -ldflags "-X github.com/variadico/noti/internal/command.Version=$(branch)-$(rev)"
ldflags_rel := -ldflags "-s -w -X github.com/variadico/noti/internal/command.Version=$(tag)"

cmd/noti/noti: $(gosrc) vendor
	go build -race -o $@ $(ldflags) github.com/variadico/noti/cmd/noti

vendor: go.mod go.sum
	go mod tidy
	go mod vendor
	touch $@

release/noti.linuxrelease: $(gosrc) vendor
	mkdir --parents release
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -o $@ $(ldflags_rel) github.com/variadico/noti/cmd/noti
release/noti$(tag).linux-amd64.tar.gz: release/noti.linuxrelease
	tar czvf $@ --transform 's#$<#noti#g' $<

release/noti.darwinrelease: $(gosrc) vendor
	mkdir -p release
	GOOS=darwin GOARCH=amd64 \
		go build -o $@ $(ldflags_rel) github.com/variadico/noti/cmd/noti
release/noti$(tag).darwin-amd64.tar.gz: release/noti.darwinrelease
	tar czvf $@ --transform 's#$<#noti#g' $<

release/noti.windowsrelease: $(gosrc) vendor
	mkdir --parents release
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 \
		go build -o $@ $(ldflags_rel) github.com/variadico/noti/cmd/noti
release/noti$(tag).windows-amd64.tar.gz: release/noti.windowsrelease
	tar czvf $@ --transform 's#$<#noti.exe#g' $<

docs/man/noti.1: docs/man/noti.1.md
	pandoc -s -t man $< -o $@
docs/man/noti.yaml.5: docs/man/noti.yaml.5.md
	pandoc -s -t man $< -o $@

.PHONY: build
build: cmd/noti/noti

.PHONY: install
install: cmd/noti/noti
	mv $< $(gobin)

.PHONY: lint
lint:
	go vet ./...
	$(golangci-lint) run --no-config --exclude-use-default=false --max-same-issues=0 \
	--timeout 30s \
	--disable errcheck \
	--disable stylecheck \
	--enable bodyclose \
	--enable golint \
	--enable interfacer \
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
	--enable scopelint \
	--enable gocritic \
	--enable gochecknoinits \
	./...

.PHONY: test
test:
	go test -v -cover -race $$(go list ./... | grep -v "noti/tests")

.PHONY: test-integration
test-integration:
	go install \
		-ldflags "-X github.com/variadico/noti/internal/command.Version=$(branch)-$(rev)" \
		github.com/variadico/noti/cmd/noti
	go test -v -cover ./tests/...

.PHONY: clean
clean:
	go clean
	rm -f cmd/noti/noti
	rm -rf release/
	git clean -x -f -d

.PHONY: man
man: docs/man/noti.1 docs/man/noti.yaml.5

.PHONY: release
release: release/noti$(tag).linux-amd64.tar.gz release/noti$(tag).windows-amd64.tar.gz

.PHONY: release-darwin
release-darwin: release/noti$(tag).darwin-amd64.tar.gz

.PHONY: update-mod
update-mod:
	go get -u ./cmd/...
	go get -u ./internal/...
	go get -u ./service/...
	go mod tidy
