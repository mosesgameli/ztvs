VERSION ?= dev
BIN_EXT ?=
LDFLAGS = -s -w -X github.com/mosesgameli/ztvs/internal/cli.Version=$(VERSION)

.PHONY: build build_host build_plugins sync_manifests init clean run-scan

build: build_host

init: build_host
	@mkdir -p $(HOME)/.ztvs/plugins
	@if [ ! -f $(HOME)/.ztvs/config.yaml ]; then \
		./zt plugin init || echo "Run './zt' to initialize config"; \
	fi

dev-setup: build init

clean:
	rm -f zt zt.exe

run-scan: build
	./zt scan
