branch = $(shell git rev-parse --abbrev-ref HEAD)
tag = $(shell git describe --abbrev=0 --tags)
rev = $(shell git rev-parse --short HEAD)

export GOFLAGS = -mod=vendor

staticcheck = ./tools/staticcheck_2019.1.1_$(shell go env GOOS)_$(shell go env GOARCH)

.PHONY: build
build:
	go build \
		-race -o cmd/noti/noti \
		-ldflags "-X github.com/variadico/noti/internal/command.Version=$(branch)-$(rev)" \
		github.com/variadico/noti/cmd/noti

.PHONY: install
install:
	go install \
		-ldflags "-X github.com/variadico/noti/internal/command.Version=$(branch)-$(rev)" \
		github.com/variadico/noti/cmd/noti

.PHONY: test
test:
	go test -v -cover -race $(shell go list ./... | grep -v "noti/tests")
	$(staticcheck) ./...

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
	git remote prune origin

.PHONY: release-linux
release-linux:
	mkdir -p release
	GOOS=linux GOARCH=amd64 \
		go build \
		-ldflags "-s -w -X github.com/variadico/noti/internal/command.Version=$(tag)" \
		github.com/variadico/noti/cmd/noti
	tar -czf release/noti$(tag).linux-amd64.tar.gz noti
	rm -f noti

.PHONY: release-darwin
release-darwin:
	mkdir -p release
	GOOS=darwin GOARCH=amd64 \
		go build \
		-ldflags "-s -w -X github.com/variadico/noti/internal/command.Version=$(tag)" \
		github.com/variadico/noti/cmd/noti
	tar -czf release/noti$(tag).darwin-amd64.tar.gz noti
	rm -f noti

.PHONY: release-windows
release-windows:
	mkdir -p release
	GOOS=windows GOARCH=amd64 \
		go build \
		-ldflags "-s -w -X github.com/variadico/noti/internal/command.Version=$(tag)" \
		github.com/variadico/noti/cmd/noti
	tar -czf release/noti$(tag).windows-amd64.tar.gz noti.exe
	rm -f noti.exe

.PHONY: man
man:
	pandoc -s -t man docs/man/noti.1.md -o docs/man/noti.1
	pandoc -s -t man docs/man/noti.yaml.5.md -o docs/man/noti.yaml.5
