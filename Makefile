GOPATH=$(PWD)/.gopath
GO=GOPATH=$(GOPATH) go
GIT_SHA=$(shell git log -n 1 --pretty=format:%H)
VERSION=1.0
PROJECT_NAME=$(notdir $(basename $(PWD)))
CODE_SUBDIR=pkg

clean:
	rm -rf $(GOPATH)

init:
	#dependency management
	mkdir -p .gopath
	$(GO) get github.com/spf13/viper
	$(GO) get github.com/spf13/pflag
	$(GO) get github.com/op/go-logging
	$(GO) get -u gopkg.in/alecthomas/gometalinter.v1
	GOPATH=$(GOPATH) $(GOPATH)/bin/gometalinter.v1 --install
	# as soon as our application does not use relative imports - source code
	# has to be present in GOPATH to make lint work
	# as soon as we created isolated GOPATH - we have to create a symlink
	# from GOPATH to our source code
	mkdir $(GOPATH)/src/$(PROJECT_NAME)/
	ln -s $(PWD)/$(CODE_SUBDIR) $(GOPATH)/src/$(PROJECT_NAME)/$(CODE_SUBDIR)

fmt:
	$(GO) fmt ./pkg/...

lint:
	GOPATH=$(GOPATH) $(GOPATH)/bin/gometalinter.v1 --disable=gotype ./pkg/...

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
