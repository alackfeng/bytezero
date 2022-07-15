##### Basic go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod


# Binary names
BINARY_EXE=$(shell go env GOEXE)
BINARY_NAME=bytezero$(BINARY_EXE)
BINARY_PACKPB=bytezeroPackPb
BINARY_ANALYSE=analyse


all: build

.PHONY: bytezeroPackPb
bytezeroPackPb:
	$(GOBUILD) -o bin/$(BINARY_PACKPB) -v main.go

.PHONY: bytezero
bytezero:
	$(GOBUILD) -o bin/$(BINARY_NAME) -v -gcflags '-N -l' main.go
	##$(GOBUILD) -o bin/$(BINARY_NAME) -v main.go

.PHONY: analyse
analyse:
	$(GOBUILD) -o bin/$(BINARY_ANALYSE) -v -gcflags '-m -m -N -l' example/analyse/analyse.go

.PHONY: build
build: bytezero

.PHONY: run
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	##./$(BINARY_NAME)

.PHONY: test
test:
	$(GOTEST) -v ./...

.PHONY: deps
deps:
	$(GOMOD) tidy -v

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f bin/$(BINARY_NAME)
	rm -f bin/$(BINARY_PACKPB)
	rm -f bin/$(BINARY_ANALYSE)
