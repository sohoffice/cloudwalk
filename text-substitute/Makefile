include ../.env
export

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=text-substitute
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

all: test build

build:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ../dist/$(VERSION)/text-substitute/darwin_amd64/text-substitute -v
	chmod a+x ../dist/$(VERSION)/text-substitute/darwin_amd64/text-substitute
	@echo "        Built darwin-amd64"
test:
	$(GOTEST) -v ./... -args -logtostderr
