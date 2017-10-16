.PHONY: test tools

test:
	golint -set_exit_status $(pkgs)
	megacheck -unused.enabled=false $(pkgs)
	go vet $(pkgs)
tools:
	go get -u honnef.co/go/tools/cmd/megacheck
	go get -u github.com/golang/lint/golint
