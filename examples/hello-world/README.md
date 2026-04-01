# Hello World Go Plugin Example

This is a comprehensive example of a ZTVS plugin implemented in Go. It demonstrates the use of the Go SDK and the new polyglot-ready manifest structure.

## Overview

A ZTVS plugin consists of two main components:
1.  **The Manifest (`plugin.yaml`)**: Metadata that tells the ZTVS host how to run the plugin.
2.  **The Executable**: A binary (or script) that implements the ZTVS RPC protocol.

## Project Structure

- `main.go`: The plugin implementation.
- `go.mod`: Dependency management.
- `plugin.yaml`: The plugin manifest.
- `Makefile`: Build automation.

## Building the Plugin

To build the plugin binary, run:

```bash
make build
```

This will create a `hello-world` executable in the current directory.

## Running with ZTVS

To test this plugin with the `zt` scanner:

1.  **Create the plugin directory**:
    ```bash
    mkdir -p ~/.ztvs/plugins/hello-world
    ```

2.  **Install the plugin**:
    Copy the binary and the manifest to the plugin directory:
    ```bash
    cp hello-world ~/.ztvs/plugins/hello-world/
    cp plugin.yaml ~/.ztvs/plugins/hello-world/
    ```

3.  **Run a scan**:
    ```bash
    zt scan
    ```

The scanner will discover the plugin, verify its manifest, and execute the `hello_world` check via the `BinaryRunner`.

## SDK Usage correctly

In `main.go`, we use `sdk.Run` to start the RPC server:

```go
sdk.Run(meta, []sdk.Check{&HelloWorldCheck{}})
```

The SDK handles:
-   The JSON-RPC handshake.
-   Dispatching `run_check` requests to your implementation.
-   Formatting and sending responses back to the host.
