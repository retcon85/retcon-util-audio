TARGET := retcon-audio
TARGETDIR := build
BINDIR := $(TARGETDIR)/binaries
VERSION ?= 0.0.0-DEV
BINROOT := $(BINDIR)/$(TARGET)-$(VERSION)
SOURCES := $(shell find . -name '*.go')

# default make [build] outputs an executable specific to this environment 
build: $(TARGETDIR)/$(TARGET)

$(TARGETDIR)/$(TARGET): $(SOURCES)
	go build -o $(TARGETDIR)/$(TARGET)

# make build/binaries outputs a complete set of cross-compiled binaries
$(BINDIR): $(SOURCES) $(BINROOT)-linux-386.bz2 $(BINROOT)-linux-amd64.bz2 $(BINROOT)-linux-arm.bz2 $(BINROOT)-linux-arm64.bz2 $(BINROOT)-darwin-amd64.bz2 $(BINROOT)-darwin-arm64.bz2 $(BINROOT)-windows-386.bz2 $(BINROOT)-windows-amd64.bz2 $(BINROOT)-windows-arm.bz2 $(BINROOT)-windows-arm64.bz2

$(BINROOT)-windows-%/$(TARGET):
	GOOS=windows GOARCH=$* go build -o $@.exe
	cp LICENSE $(BINROOT)-windows-$*

$(BINROOT)-%/$(TARGET):
	GOOS=$(firstword $(subst -, ,$*)) GOARCH=$(word 2,$(subst -, ,$*)) go build -o $@
	cp LICENSE $(BINROOT)-$*

$(BINROOT)-%.bz2: $(BINROOT)-%/$(TARGET)
	tar -jcC $(BINROOT)-$* . > $@
	rm -rf $(BINROOT)-$*

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -rf $(TARGETDIR)
	