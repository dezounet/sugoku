GOARCH = amd64

CURRENT_DIR=$(shell pwd)
BUILD_DIR=${CURRENT_DIR}/cmd/sugokud/

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS =

linux:
	cd ${BUILD_DIR}; \
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o sugokud-linux-${GOARCH} . ; \
	cd - >/dev/null

darwin:
	cd ${BUILD_DIR}; \
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o sugokud-darwin-${GOARCH} . ; \
	cd - >/dev/null

all: linux darwin
