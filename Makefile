BIN_DIR = $(CURDIR)/bin
GOLANGCI_LINT = $(BIN_DIR)/golangci-lint
GOSEC = $(BIN_DIR)/gosec
GOVULNCHECK = $(BIN_DIR)/govulncheck

.PHONY: build build_host build_plugins sync_manifests init clean run-scan

build: build_host

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

GOLANGCI_LINT_VER = v1.64.5
GOSEC_VER = v2.22.1
GOVULNCHECK_VER = latest

init: build_host
	@mkdir -p $(HOME)/.ztvs/plugins
	@if [ ! -f $(HOME)/.ztvs/config.yaml ]; then \
		./zt plugin init || echo "Run './zt' to initialize config"; \
	fi

dev-setup: build init

clean:
	rm -f zt zt.exe
	rm -rf $(BIN_DIR)

run-scan: build
	./zt scan
