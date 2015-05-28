GO		:= $(shell which go)
GO_PKG		:= github.com/nesv/pm
GO_FILES	:= $(shell find . -type f -iname "*.go")

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
package: install-deps pm-bootstrap bin/pm ${METADATA_FILE}
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
