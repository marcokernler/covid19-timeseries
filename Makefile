
#
OUT_DIR=./bin
BINARY=covid19-timeseries
VERSION=1.1.0
BUILD=$(shell git rev-parse HEAD)
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
build_all: clean linux darwin windows

linux:
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=linux; export GOARCH=$(GOARCH); go build -v -o $(OUT_DIR)/$(BINARY)-linux-$(GOARCH)))

darwin:
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=darwin; export GOARCH=$(GOARCH); go build -v -o $(OUT_DIR)/$(BINARY)-darwin-$(GOARCH)))

windows:
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=windows; export GOARCH=$(GOARCH); go build -v -o $(OUT_DIR)/$(BINARY)-windows-$(GOARCH).exe))

#
clean:
	find ${OUT_DIR} -name '${BINARY}[-?][a-zA-Z0-9]*[-?][a-zA-Z0-9]*' -delete

#
.PHONY: check clean build_all all