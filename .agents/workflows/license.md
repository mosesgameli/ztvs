---
description: Verify that all source files contain the Apache License 2.0 header.
---

// turbo-all

1. Find files missing the license header (Go, Shell, PowerShell):
`grep -L "Licensed under the Apache License, Version 2.0" $(find . -name "*.go" -o -name "*.sh" -o -name "*.ps1" | grep -v "vendor")`

2. Apply header to all .go files if missing (requires manual confirmation):
`find . -name "*.go" | xargs grep -L "Apache License, Version 2.0" | xargs -I {} sed -i '1i // Licensed under the Apache License, Version 2.0\n// ...' {}`
