description: Build the ZTVS host and all first-party plugins for the current platform.

1. Compile the host and first-party plugins:
// turbo
`make build`

2. Verify all binaries are generated:
// turbo
`ls -l zt plugins/plugin-os/plugin-os`

3. Verify cross-compilation support:
// turbo
`GOOS=linux GOARCH=amd64 go build -o zt-linux ./cmd/zt`
