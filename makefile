GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install
BUILDDIR=$(PWD)/.bin

all: prep build test s3fileupload
.PHONY: prep
prep:
	mkdir -p $(BUILDDIR)
.PHONY: build
build:
	$(GOBUILD) -v ./...
.PHONY: test
test: 
	$(GOTEST) -v ./...
.PHONY: s3fileupload
s3fileupload:
	$(GOBUILD) -o $(BUILDDIR)/s3fileupload ./cmd/fileupload
