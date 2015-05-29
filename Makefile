GO		:= $(shell which go)
GO_PKG		:= github.com/nesv/pm
GO_FILES	:= $(shell find . -type f -iname "*.go")

TEST_DIR	:= ./test
TEST_CACHE_DIR	:= ${TEST_DIR}/cache
TEST_BASE_DIR	:= ${TEST_DIR}/pkgs
TEST_BIN_DIR	:= ${TEST_DIR}/bin
PM_TEST_FLAGS	:= -C ${TEST_CACHE_DIR} -d ${TEST_BASE_DIR} -B ${TEST_BIN_DIR}

ifndef arch
	arch = ${GOARCH}
endif
ifndef platform
	platform = ${GOOS}
endif

all: install-deps bin/pm

install-deps:
	gpm install

bin/%: ${GO_FILES}
	@echo "=== Building $@"
	GOOS=${platform} GOARCH=${arch} \
	     go build -o $@ ${GO_PKG}/cmd/$(shell basename $@)

METADATA_FILE := metadata-${platform}-${arch}.json
package: install-deps \
	pm-${version}-${platform}-${arch}.tar.gz

pm-${version}-${platform}-${arch}.tar.gz: \
	pm-bootstrap \
	${METADATA_FILE} \
	bin/pm
	./pm-bootstrap build --metadata=metadata-${platform}-${arch}.json

pm-bootstrap: ${GO_FILES}
	go build -o $@ ${GO_PKG}/cmd/pm

${METADATA_FILE}: metadata.json
ifndef version
	$(error "version is not set")
endif
	cat $^ | sed \
		-e 's/{{.Version}}/${version}/' \
		-e 's/{{.Architecture}}/${arch}/' \
		-e 's/{{.Platform}}/${platform}/' \
		> $@

live-test: ${TEST_BASE_DIR} ${TEST_CACHE_DIR} ${TEST_BIN_DIR} \
	package \
	live-test-unpack \
	live-test-install

${TEST_DIR}/%:
	mkdir -p $@

live-test-fetch: bin/pm clean-test
	@echo "=== Testing fetch"
	bin/pm ${PM_TEST_FLAGS} fetch pm-${version}-${platform}-${arch}.tar.gz

live-test-unpack: bin/pm live-test-fetch
	@echo "=== Testing unpack"
	bin/pm ${PM_TEST_FLAGS} unpack pm 0.1.0

live-test-install: bin/pm
	@echo "=== Testing install"
	bin/pm ${PM_TEST_FLAGS} install pm-${version}-${platform}-${arch}.tar.gz

clean-test:
	rm -rvf ${TEST_BASE_DIR}/*
	rm -rvf ${TEST_CACHE_DIR}/*
	rm -rvf ${TEST_BIN_DIR}/*

clean:
	@echo "=== Cleaning up"
	@rm -rvf bin
	@rm -vf pm-bootstrap
	@rm -vf pm-*-${platform}-${arch}.tar.gz
	@rm -vf ${METADATA_FILE}

clean-all:
	@echo "=== Cleaning everything up"
	@rm -rvf bin
	@rm -vf pm-*-*-*.tar.gz
	@rm -vf metadata-*.json

.PHONY: \
	install-deps \
	clean \
	package
