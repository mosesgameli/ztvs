---
description: Execute the phased testing strategy.
---

// turbo-all

1. Run unit tests with coverage report (Target: 95%+):
`go test -v -coverprofile=coverage.out ./internal/... ./pkg/...`

2. View coverage by function:
`go tool cover -func=coverage.out`

3. Run protocol integration tests:
`go test -v ./test/fixtures/...`

4. Run installer smoke test (Unix):
`chmod +x test/e2e/installer_smoke_test.sh && ./test/e2e/installer_smoke_test.sh`
