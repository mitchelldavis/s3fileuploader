GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install
BUILDDIR=$(PWD)/.bin

all: prep build test s3fileupload darwin linux windows
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
.PHONY: darwin
darwin:
	export GOOS=darwin
	export GOARCH=arm64
	$(GOBUILD) -o $(BUILDDIR)/s3fileupload_darwin_amd64 ./cmd/fileupload
.PHONY: linux
linux:
	export GOOS=linux 
	export GOARCH=arm64 
	$(GOBUILD) -o $(BUILDDIR)/s3fileupload_linux_amd64 ./cmd/fileupload
.PHONY: windows
windows:
	export GOOS=windows 
	export GOARCH=arm64 
	$(GOBUILD) -o $(BUILDDIR)/s3fileupload_windows_amd64.exe ./cmd/fileupload
