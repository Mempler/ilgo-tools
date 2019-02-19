GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

BUILD_DIRECTORY=$(shell pwd)/build

.PHONY: all test clean gpacker

all: test build

test:
	echo "Run Tests"
	$(GOTEST) -v ./...

build: gen_build_dir gpacker

clean:
	echo "Cleanup"
	$(GOCLEAN)
	rm -rf $(BUILD_DIRECTORY)

gen_build_dir:
	mkdir -p $(BUILD_DIRECTORY)

gpacker:
	echo "Build GPacker"
	$(GOBUILD) -o $(BUILD_DIRECTORY)/gpacker -v gpacker/executable/main.go

