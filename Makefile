.PHONY: build build_host build_plugins sync_manifests init clean run-scan

build: build_host build_plugins sync_manifests

build_host:
	go build -o zt ./cmd/zt

build_plugins:
	@mkdir -p plugins/plugin-os
	go build -C plugins/plugin-os -o plugin-os .
	go build -C plugins/plugin-axios-mitigation -o plugin-axios-mitigation .

sync_manifests:
	go run ./tools/manifest-sync ./plugins/plugin-os
	go run ./tools/manifest-sync ./plugins/plugin-axios-mitigation

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
