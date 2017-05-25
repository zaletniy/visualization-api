GOPATH=$(PWD)/.gopath
GO=GOPATH=$(GOPATH) go
GIT_SHA=$(shell git log -n 1 --pretty=format:%H)
VERSION=1.0

clean:
	rm -rf $(GOPATH)

init:
	#dependency management
	mkdir -p .gopath
	$(GO) get github.com/spf13/viper
	$(GO) get github.com/spf13/pflag
	$(GO) get -u gopkg.in/alecthomas/gometalinter.v1
	gometalinter.v1 --install

fmt:
	$(GO) fmt ./pkg/...

lint:
	gometalinter.v1 ./pkg/...

test:
	$(GO) test ./pkg/...

build: fmt lint
	$(GO) build ./pkg/cmd/...

build-all: fmt lint
	mkdir -p build/linux-amd64
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags "-X main.Version=$(VERSION) -X main.GitVersion=$(GIT_SHA)" -o $(PWD)/build/linux-amd64/visualizationapi ./pkg/cmd/visualizationapi

package:
	exit 1

all: init build-all package
