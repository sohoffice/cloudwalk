include ../.env
export

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=heroku-gcp-migrate-tools
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

all: test build

build:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ../dist/$(VERSION)/heroku-gcp-migrate-tools/darwin_amd64/heroku-gcp-migrate-tools -v
	chmod a+x ../dist/$(VERSION)/heroku-gcp-migrate-tools/darwin_amd64/heroku-gcp-migrate-tools
	@echo "        Built darwin-amd64"
test:
	$(GOTEST) -v ./... -args -logtostderr
