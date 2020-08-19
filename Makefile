GOARCH = amd64

CURRENT_DIR=$(shell pwd)
BUILD_DIR=${CURRENT_DIR}/cmd/sugokud/
DOCKERFILE_DIR=${CURRENT_DIR}/build/

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

docker: linux
	cp ${BUILD_DIR}/sugokud-linux-${GOARCH} ${DOCKERFILE_DIR}/sugokud-linux-${GOARCH} ; \
	cd ${DOCKERFILE_DIR} ; \
	docker build -t sugokud:latest . ; \
	rm sugokud-linux-${GOARCH}
