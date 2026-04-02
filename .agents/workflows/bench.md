---
description: Execute the engine performance benchmark suite.
---

// turbo-all

1. Run all benchmarks in the engine package:
`go test -v -bench . -benchmem ./internal/engine/...`

1. Compare results (if benchstat is installed):
`go test -v -bench . ./internal/engine/... > old.txt && (your-change) && go test -v -bench . ./internal/engine/... > new.txt && benchstat old.txt new.txt`
