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
	$(GO) get github.com/go-sql-driver/mysql
	$(GO) get github.com/go-xorm/xorm
	$(GO) get github.com/pressly/chi
	$(GO) get github.com/auth0/go-jwt-middleware
	$(GO) get github.com/dgrijalva/jwt-go
	$(GO) get -u gopkg.in/alecthomas/gometalinter.v1
	GOPATH=$(GOPATH) $(GOPATH)/bin/gometalinter.v1 --install
	# as soon as our application does not use relative imports - source code
	# has to be present in GOPATH to make lint work
	# as soon as we created isolated GOPATH - we have to create a symlink
	# from GOPATH to our source code
	mkdir -p $(GOPATH)/src/$(PROJECT_NAME)/
	ln -s $(PWD)/$(CODE_SUBDIR) $(GOPATH)/src/$(PROJECT_NAME)/$(CODE_SUBDIR)

fmt:
	$(GO) fmt ./pkg/...

lint:
	GOPATH=$(GOPATH) $(GOPATH)/bin/gometalinter.v1 --disable=gotype ./pkg/...

test:
	$(GO) test ./pkg/...

build: fmt lint
	$(GO) build ./pkg/cmd/...

build-all: fmt
	mkdir -p build/linux-amd64
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags "-X main.version=$(VERSION) -X main.gitVersion=$(GIT_SHA)" -o $(PWD)/build/linux-amd64/visualizationapi ./pkg/cmd/visualizationapi

package-init:
	docker build -t com.mirantis.pv/build ./tools/build

package-clean:
	docker image rm -f com.mirantis.pv/build
	rm -rf build/deb/*

package:
	docker run -e VERSION=$(VERSION) -v $(PWD):/app com.mirantis.pv/build /app/tools/build/build_deb.sh

package-debug:
	docker run -it -v $(PWD):/app com.mirantis.pv/build /bin/bash

docker:
	docker build -t $(DOCKER_USERNAME)/visualization-api -f tools/docker/visualization-api/Dockerfile .

docker-push:
	docker push $(DOCKER_USERNAME)/visualization-api

all: init fmt lint  build-all package-init package docker
