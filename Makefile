
#
OUT_DIR=./bin
BINARY=covid19-timeseries
VERSION=1.1.0
BUILD=$(shell git rev-parse HEAD)
PLATFORMS=darwin linux windows
ARCHITECTURES=386 amd64

#
LDFLAGS=-ldflags "-X 'main.Version=${VERSION}' -X 'main.Build=${BUILD}'"

#
default: build

#
all: clean build_all

#
build:
	go build ${LDFLAGS} -o $(OUT_DIR)/${BINARY}

#
build_all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build -v -o $(OUT_DIR)/$(BINARY)-$(GOOS)-$(GOARCH))))

#
clean:
	find ${OUT_DIR} -name '${BINARY}[-?][a-zA-Z0-9]*[-?][a-zA-Z0-9]*' -delete

#
.PHONY: check clean build_all all