PKGROOT=github.com/Ladicle/kubectl-rolesum

CMD=kubectl-rolesum
VERSION ?= $(shell git describe --abbrev=0 --tags 2>/dev/null || echo no-version)
GIT_COMMIT := $(shell git rev-parse --short HEAD)

OUTDIR=_output

GOLDFLAGS=-w -X $(PKGROOT)/cmd.version=$(VERSION) -X $(PKGROOT)/cmd.command=$(CMD) -X $(PKGROOT)/cmd.commit=$(GIT_COMMIT)

export GO111MODULE=on

.PHONY: build build-linux build-darwin build-windows install check clean

build:
	CGO_ENABLED=0 \
	go build -ldflags "$(GOLDFLAGS)" -o $(OUTDIR)/$(CMD)

build-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build -ldflags "$(GOLDFLAGS)" -o $(OUTDIR)/$(CMD)_linux-amd64/$(CMD)

build-darwin:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 \
	go build -ldflags "$(GOLDFLAGS)" -o $(OUTDIR)/$(CMD)_darwin-amd64/$(CMD)

build-windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
	go build -ldflags "$(GOLDFLAGS)" -o $(OUTDIR)/$(CMD)_windows-amd64/$(CMD)

install:
	CGO_ENABLED=0 go install -ldflags "$(GOLDFLAGS)"

check: build
	go vet $(PKGROOT)/...
	./test.sh

clean:
	-rm -r $(OUTDIR)
