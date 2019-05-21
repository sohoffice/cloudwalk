include .env
export

# Go parameters
PROJECT_NAME=cloudwalk
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

all: test build

build: build-heroku-gcp-migrate-tools build-text-substitute
test:
	$(GOTEST) -v ./... -args -logtostderr
tests:
	$(GOTEST) -v ./... -count=10 -args -logtostderr
clean:
	$(GOCLEAN)
	rm -rf dist

# Build on all supported platforms
build-all: build
	@echo "\n$(PROJECT_NAME) version: $(VERSION)\n"

# Cross compilation on linux
build-linux:
	cd main && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ../dist/$(VERSION)/linux_amd64/$(BINARY_NAME) -v && cd ..
	chmod a+x dist/$(VERSION)/linux_amd64/$(BINARY_NAME)
	@echo "        Built linux-amd64"

build-heroku-gcp-migrate-tools:
	cd heroku-gcp-migrate-tools && make

build-text-substitute:
	cd text-substitute && make
