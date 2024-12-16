GOPATH := $(shell go env GOPATH)
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
FREE_VERSION := "0.3.3"

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
	GOOS=${GOOS} GOARCH=${GOARCH} go build -o "${GOOS}/free"
	# sleep 2
	# tar -C "${GOOS}" -czvf "free_${FREE_VERSION}_${GOOS}_${GOARCH}.tgz" free free.1; \

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