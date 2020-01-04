PKGNAME    = beluga
DESTDIR   ?=
PREFIX    ?= /usr
BINDIR     = $(PREFIX)/bin
SYSCONFDIR = /etc

GOBIN       = _build/bin
GOPROJROOT  = $(GOSRC)/$(PROJREPO)

GOLDFLAGS   = -ldflags "-s -w"
GOTAGS      = --tags "libsqlite3 linux"
GOCC        = go
GOFMT       = $(GOCC) fmt -x
GOGET       = $(GOCC) get $(GOLDFLAGS)
GOBUILD     = $(GOCC) build -v $(GOLDFLAGS) $(GOTAGS)
GOTEST      = $(GOCC) test
GOVET       = $(GOCC) vet
GOINSTALL   = $(GOCC) install $(GOLDFLAGS)

include Makefile.waterlog

GOLINT = golint -set_exit_status

all: build

build:
	@$(call stage,BUILD)
	@$(GOBUILD)
	@$(call pass,BUILD)

test: build
	@$(call stage,TEST)
	@$(GOTEST) ./...
	@$(call pass,TEST)

validate:
	@$(call stage,FORMAT)
	@$(GOFMT) ./...
	@$(call pass,FORMAT)
	@$(call stage,VET)
	@$(call task,Running 'go vet'...)
	@$(GOVET) ./...
	@$(call pass,VET)
	@$(call stage,LINT)
	@$(call task,Running 'golint'...)
	@$(GOLINT) ./...
	@$(call pass,LINT)


install:
	@$(call stage,INSTALL)
	install -Dm 00755 $(PKGNAME) $(DESTDIR)$(BINDIR)/$(PKGNAME)
	install -Dm 00644 data/beluga.conf $(SYSCONFDIR)/beluga/beluga.conf
	@$(call pass,INSTALL)

uninstall:
	@$(call stage,UNINSTALL)
	rm -f $(DESTDIR)$(BINDIR)/$(PKGNAME)
	@$(call pass,UNINSTALL)

clean:
	@$(call stage,CLEAN)
	@$(call task,Removing _build directory...)
	@rm -rf _build
	@$(call task,Removing executable...)
	@rm $(PKGNAME)
	@$(call pass,CLEAN)