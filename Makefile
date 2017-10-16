branch = $(shell git rev-parse --abbrev-ref HEAD)
tag = $(shell git describe --abbrev=0 --tags)
rev = $(shell git rev-parse --short HEAD)
pkgs = $(shell go list ./... | grep -v /vendor/)

.PHONY: build install tools test update-deps clean

build:
	go build -race -o cmd/noti/noti \
		-ldflags "-X github.com/variadico/noti/internal/command.Version=$(branch)-$(rev)" \
		github.com/variadico/noti/cmd/noti
install:
	go install \
		-ldflags "-X github.com/variadico/noti/internal/command.Version=$(branch)-$(rev)" \
		github.com/variadico/noti/cmd/noti
tools:
	go install ./vendor/github.com/golang/dep/cmd/dep
	go install ./vendor/honnef.co/go/tools/cmd/megacheck
	go install ./vendor/github.com/golang/lint/golint
test:
	golint -set_exit_status $(pkgs)
	megacheck $(pkgs)
	go vet $(pkgs)
	go test -v -cover -race $(pkgs)
update-deps:
	dep ensure
	dep ensure -update
	dep prune
clean:
	go clean
	rm -f cmd/noti/noti
	git clean -x -f -d
	git remote prune origin
