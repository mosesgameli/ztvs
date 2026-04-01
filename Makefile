VERSION ?= dev
BIN_EXT ?=
LDFLAGS = -s -w -X github.com/mosesgameli/ztvs/internal/cli.Version=$(VERSION)

.PHONY: build build_host build_plugins sync_manifests init clean run-scan

build: build_host build_plugins sync_manifests

build_host:
	go build -ldflags="$(LDFLAGS)" -o zt$(BIN_EXT) ./cmd/zt

build_plugins:
	@mkdir -p plugins/plugin-os
	go build -C plugins/plugin-os -ldflags="-s -w" -o plugin-os$(BIN_EXT) .
	go build -C plugins/plugin-axios-mitigation -ldflags="-s -w" -o plugin-axios-mitigation$(BIN_EXT) .

sync_manifests:
	GOOS="" GOARCH="" go run ./tools/manifest-sync ./plugins/plugin-os $(BIN_EXT)
	GOOS="" GOARCH="" go run ./tools/manifest-sync ./plugins/plugin-axios-mitigation $(BIN_EXT)

init: build_host
	@mkdir -p $(HOME)/.ztvs/plugins
	@if [ ! -f $(HOME)/.ztvs/config.yaml ]; then \
		./zt plugin init || echo "Run './zt' to initialize config"; \
	fi

dev-setup: build init

clean:
	rm -f zt
	rm -rf plugins/plugin-os/plugin-os
	rm -rf plugins/plugin-axios-mitigation/plugin-axios-mitigation

run-scan: build
	./zt scan
