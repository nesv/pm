GO		:= $(shell which go)
GO_PKG		:= github.com/nesv/pm
GO_FILES	:= $(shell find . -type f -iname "*.go")

TEST_DIR	:= ./test
TEST_CACHE_DIR	:= ${TEST_DIR}/cache
TEST_BASE_DIR	:= ${TEST_DIR}/pkgs
TEST_BIN_DIR	:= ${TEST_DIR}/bin
PM_TEST_FLAGS	:= -C ${TEST_CACHE_DIR} -D ${TEST_BASE_DIR} -B ${TEST_BIN_DIR} -v

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
PM_PACKAGE_TAR_GZ := pm-${version}-${platform}-${arch}.tar.gz

package: ${PM_PACKAGE_TAR_GZ}

${PM_PACKAGE_TAR_GZ}: \
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
	live-test-link \
	live-test-install \
	live-test-list-cached \
	live-test-list-linked \
	live-test-list-unpacked \
	live-test-list-cached-unpacked-linked \
	live-test-clean

${TEST_DIR}/%:
	mkdir -p $@

live-test-fetch: bin/pm clean-test ${PM_PACKAGE_TAR_GZ}
	@echo "=== Testing fetch"
	bin/pm ${PM_TEST_FLAGS} fetch ${PM_PACKAGE_TAR_GZ}

live-test-unpack: bin/pm live-test-fetch
	@echo "=== Testing unpack"
	bin/pm ${PM_TEST_FLAGS} unpack pm-0.1.0

live-test-link: bin/pm live-test-unpack
	@echo "=== Testing link"
	bin/pm ${PM_TEST_FLAGS} link pm-${version}

live-test-install: clean-test bin/pm ${PM_PACKAGE_TAR_GZ}
	@echo "=== Testing install"
	bin/pm ${PM_TEST_FLAGS} install ${PM_PACKAGE_TAR_GZ}

live-test-clean: bin/pm ${PM_PACKAGE_TAR_GZ}
	@echo "=== Testing clean --all"
	bin/pm ${PM_TEST_FLAGS} clean --all

live-test-list-linked: bin/pm
	@echo "=== Testing list --linked"
	bin/pm ${PM_TEST_FLAGS} list --linked

live-test-list-cached: bin/pm
	@echo "=== Testing list --cached"
	bin/pm ${PM_TEST_FLAGS} list --cached

live-test-list-unpacked: bin/pm
	@echo "=== Testing list --unpacked"
	bin/pm ${PM_TEST_FLAGS} list --unpacked

live-test-list-cached-unpacked-linked: bin/pm
	@echo "=== Testing list -cxi"
	bin/pm ${PM_TEST_FLAGS} list -cxi

clean-test:
	rm -rvf ${TEST_BASE_DIR}/*
	rm -rvf ${TEST_CACHE_DIR}/*
	rm -rvf ${TEST_BIN_DIR}/*

clean:
	@echo "=== Cleaning up"
	@rm -rf bin
	@rm -f pm-bootstrap
	@rm -f pm-*-*-*.tar.gz
	@rm -f metadata-*.json
	@rm -f ${METADATA_FILE}

clean-all: clean clean-test

.PHONY: \
	install-deps \
	clean \
	clean-test \
	package \
	live-test \
	live-test-fetch \
	live-test-unpack \
	live-test-install
