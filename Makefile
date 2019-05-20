REPO_NAME=kubectl-bindrole
PKGROOT=github.com/Ladicle/kubectl-bindrole
VERSION ?= $(shell git rev-parse --short HEAD)
OUTDIR=_output

build:
	GO111MODULE=on CGO_ENABLED=0 go build -ldflags "-w -X $(PKGROOT)/cmd.version=$(VERSION) -X $(PKGROOT)/cmd.command=$(REPO_NAME)" -o $(OUTDIR)/kubectl-bindrole

build-linux:
	GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-w -X $(PKGROOT)/cmd.version=$(VERSION) -X $(PKGROOT)/cmd.command=$(REPO_NAME)" -o $(OUTDIR)/kubectl-bindrole_linux-amd64/kubectl-bindrole

build-darwin:
	GO111MODULE=on GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-w -X $(PKGROOT)/cmd.version=$(VERSION) -X $(PKGROOT)/cmd.command=$(REPO_NAME)" -o $(OUTDIR)/kubectl-bindrole_darwin-amd64/kubectl-bindrole

install:
	GO111MODULE=on CGO_ENABLED=0 go install -ldflags "-w -X $(PKGROOT)/cmd.version=$(VERSION) -X $(PKGROOT)/cmd.command=$(REPO_NAME)"
check:
	GO111MODULE=on go vet $(PKGROOT)/...
	./test

.PHONY: clean
clean:
	-rm -r $(OUTDIR)
