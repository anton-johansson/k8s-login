.DEFAULT_GOAL := build

BINARY = k8s-login
VERSION = v0.0.2

GO_VERSION = $(shell go version | awk -F\go '{print $$3}' | awk '{print $$1}')
COMMIT = $(shell git rev-parse HEAD)
PACKAGE_LIST = $$(go list ./...)
OUTPUT_DIRECTORY = ./build
LDFLAGS = -ldflags "\
	-X github.com/anton-johansson/k8s-login/version.version=${VERSION} \
	-X github.com/anton-johansson/k8s-login/version.goVersion=${GO_VERSION} \
	-X github.com/anton-johansson/k8s-login/version.commit=${COMMIT} \
	"

install:
	go get -v -d ./...

fmt:
	gofmt -s -d -e -w .

vet:
	go vet ${PACKAGE_LIST}

test: install
	go test ${PACKAGE_LIST}

linux: install
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${OUTPUT_DIRECTORY}/${BINARY}-linux-amd64 .

darwin: install
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${OUTPUT_DIRECTORY}/${BINARY}-darwin-amd64 .

windows: install
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${OUTPUT_DIRECTORY}/${BINARY}-windows-amd64.exe .

build: linux darwin windows

clean:
	rm -rf ${OUTPUT_DIRECTORY}
