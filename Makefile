branch = $(shell git rev-parse --abbrev-ref HEAD)
rev = $(shell git rev-parse --short HEAD)
pkgs = $(shell go list ./... | grep -v /vendor/)

.PHONY: build install tools test update clean

build:
	go build -race github.com/variadico/noti/cmd/noti
install:
	go install github.com/variadico/noti/cmd/noti
tools:
	go install ./vendor/github.com/golang/dep/cmd/dep
	go install ./vendor/honnef.co/go/tools/cmd/megacheck
	go install ./vendor/github.com/golang/lint/golint
test:
	golint -set_exit_status $(pkgs)
	megacheck -unused.enabled=false $(pkgs)
	go vet $(pkgs)
	go test -cover -race $(pkgs)
update:
	dep ensure
	dep ensure -update
	dep prune
clean:
	rm -f cmd/noti/noti
	git clean -x -f -d
	git remote prune origin
