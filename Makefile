REPO_NAME=kubectl-bindrole
PKGROOT=github.com/Ladicle/kubectl-bindrole
VERSION ?= $(shell git describe --abbrev=0 --tags 2>/dev/null || echo no-version)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
OUTDIR=_output

build:
	GO111MODULE=on CGO_ENABLED=0 go build -ldflags "-w -X $(PKGROOT)/cmd.version=$(VERSION) -X $(PKGROOT)/cmd.command=$(REPO_NAME) -X $(PKGROOT)/cmd.commit=$(GIT_COMMIT)" -o $(OUTDIR)/kubectl-bindrole

build-linux:
	GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-w -X $(PKGROOT)/cmd.version=$(VERSION) -X $(PKGROOT)/cmd.command=$(REPO_NAME) -X $(PKGROOT)/cmd.commit=$(GIT_COMMIT)" -o $(OUTDIR)/kubectl-bindrole_linux-amd64/kubectl-bindrole

build-darwin:
	GO111MODULE=on GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-w -X $(PKGROOT)/cmd.version=$(VERSION) -X $(PKGROOT)/cmd.command=$(REPO_NAME) -X $(PKGROOT)/cmd.commit=$(GIT_COMMIT)" -o $(OUTDIR)/kubectl-bindrole_darwin-amd64/kubectl-bindrole

build-windows:
	GO111MODULE=on GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-w -X $(PKGROOT)/cmd.version=$(VERSION) -X $(PKGROOT)/cmd.command=$(REPO_NAME) -X $(PKGROOT)/cmd.commit=$(GIT_COMMIT)" -o $(OUTDIR)/kubectl-bindrole_windows-amd64/kubectl-bindrole


install:
	GO111MODULE=on CGO_ENABLED=0 go install -ldflags "-w -X $(PKGROOT)/cmd.version=$(VERSION) -X $(PKGROOT)/cmd.command=$(REPO_NAME) -X $(PKGROOT)/cmd.commit=$(GIT_COMMIT)"

check:
	GO111MODULE=on go vet $(PKGROOT)/...
	./test.sh

.PHONY: clean
clean:
	-rm -r $(OUTDIR)
