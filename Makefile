branch = $(shell git rev-parse --abbrev-ref HEAD)
tag = $(shell git describe --abbrev=0 --tags)
rev = $(shell git rev-parse --short HEAD)
pkgs = $(shell go list ./... | grep -v /vendor/)

.PHONY: build install install-tools lint-only test-only test update-deps clean man release-macos release-linux release-windows release

build:
	go build -race -o cmd/noti/noti \
		-ldflags "-X github.com/variadico/noti/internal/command.Version=$(branch)-$(rev)" \
		github.com/variadico/noti/cmd/noti
install:
	go install \
		-ldflags "-X github.com/variadico/noti/internal/command.Version=$(branch)-$(rev)" \
		github.com/variadico/noti/cmd/noti
install-tools:
	go install ./vendor/github.com/golang/dep/cmd/dep
	go install ./vendor/honnef.co/go/tools/cmd/megacheck
	go install ./vendor/github.com/golang/lint/golint
lint-only:
	golint -set_exit_status $(pkgs)
	megacheck $(pkgs)
	go vet $(pkgs)
test-only:
	go test -v -cover -race $(pkgs)
test: lint-only test-only
update-deps:
	dep ensure
	dep ensure -update
	dep prune
clean:
	go clean
	rm -f cmd/noti/noti
	git clean -x -f -d
	git remote prune origin
man:
	pandoc -s -t man docs/man/noti.1.md -o docs/man/noti.1
	pandoc -s -t man docs/man/noti.yaml.5.md -o docs/man/noti.yaml.5
release-macos:
	GOOS=darwin GOARCH=amd64 \
		go build \
		-ldflags "-s -w -X github.com/variadico/noti/internal/command.Version=$(tag)" \
		github.com/variadico/noti/cmd/noti
	tar -czf noti$(tag).darwin-amd64.tar.gz noti
	rm -f noti
release-linux:
	GOOS=linux GOARCH=amd64 \
		go build \
		-ldflags "-s -w -X github.com/variadico/noti/internal/command.Version=$(tag)" \
		github.com/variadico/noti/cmd/noti
	tar -czf noti$(tag).linux-amd64.tar.gz noti
	rm -f noti
release-windows:
	GOOS=windows GOARCH=amd64 \
		go build \
		-ldflags "-s -w -X github.com/variadico/noti/internal/command.Version=$(tag)" \
		github.com/variadico/noti/cmd/noti
	tar -czf noti$(tag).windows-amd64.tar.gz noti.exe
	rm -f noti.exe
release: release-macos release-linux release-windows
