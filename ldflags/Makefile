GO=go
VERSION=beta

LDFLAGS += -X "main.BuildTS=$(shell date -u '+%Y-%m-%d %I:%M:%S')"
LDFLAGS += -X "main.GitHash=$(shell git rev-parse --short HEAD)"
LDFLAGS += -X "main.Version=$(VERSION)"
LDFLAGS += -X "main.GitBranch=$(shell git rev-parse --abbrev-ref HEAD)"
SRC = main.go

all: build

build: $(SRC)
	$(GO) build -ldflags '$(LDFLAGS)' main.go

clean:
	rm -f main