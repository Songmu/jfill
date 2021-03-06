VERSION = $(shell gobump show -r)
CURRENT_REVISION = $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS = "-X github.com/Songmu/jfill.revision=$(CURRENT_REVISION)"
ifdef update
  u=-u
endif

deps:
	go get ${u} github.com/golang/dep/cmd/dep
	dep ensure

devel-deps: deps
	go get ${u} golang.org/x/lint/golint
	go get ${u} github.com/mattn/goveralls
	go get ${u} github.com/motemen/gobump/cmd/gobump
	go get ${u} github.com/Songmu/goxz/cmd/goxz
	go get ${u} github.com/Songmu/ghch/cmd/ghch
	go get ${u} github.com/tcnksm/ghr

test: deps
	go test

lint: devel-deps
	go vet
	golint -set_exit_status

cover: devel-deps
	goveralls

build: deps
	go build -ldflags=$(BUILD_LDFLAGS) ./cmd/jfill

crossbuild: devel-deps
	$(eval ver = $(shell gobump show -r))
	goxz -pv=v$(ver) -build-ldflags=$(BUILD_LDFLAGS) \
	  -d=./dist/v$(ver) ./cmd/jfill

bump: devel-deps
	_tools/releng

upload:
	ghr v$(VERSION) dist/v$(VERSION)

release: bump crossbuild upload

.PHONY: test deps devel-deps lint cover crossbuild release
