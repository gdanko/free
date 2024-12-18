GOPATH := $(shell go env GOPATH)
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
FREE_VERSION := "0.3.4"

GOOS ?= $(shell uname | tr '[:upper:]' '[:lower:]')
GOARCH ?=$(shell arch)

.PHONY: all build install

all: build install

.PHONY: mod-tidy
mod-tidy:
	go mod tidy

.PHONY: build OS ARCH
build: guard-FREE_VERSION mod-tidy clean
	@echo "================================================="
	@echo "Building free"
	@echo "=================================================\n"

	@if [ ! -d "${GOOS}" ]; then \
		mkdir "${GOOS}"; \
	fi
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o "bin/free"
	sleep 2
	tar -czvf "free_${FREE_VERSION}_${GOOS}_${GOARCH}.tgz" ./bin ./share; \

.PHONY: clean
clean:
	@echo "================================================="
	@echo "Cleaning free"
	@echo "=================================================\n"
	@for OS in darwin; do \
		if [ -f $${OS}/free ]; then \
			rm -f $${OS}/free; \
		fi; \
	done

.PHONY: clean-all
clean-all: clean
	@echo "================================================="
	@echo "Cleaning tarballs"
	@echo "=================================================\n"
	@rm -f *.tgz 2>/dev/null

.PHONY: install
install:
	@echo "================================================="
	@echo "Installing free in ${GOPATH}/bin"
	@echo "=================================================\n"

	go install -race

#
# General targets
#
guard-%:
	@if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi