# ZTVS Go Starter Codebase Skeleton (Implements RFC-001 + RFC-002)

**Language:** Go
**Architecture:** core CLI + out-of-process JSON-RPC plugins
**Goal:** implementation-ready engineering bootstrap

---

# 1. Repository Layout

```text
ztvs/
├── go.mod
├── go.sum
├── Makefile
├── README.md
│
├── cmd/
│   └── zt/
│       └── main.go
│
├── internal/
│   ├── cli/
│   │   └── root.go
│   │
│   ├── engine/
│   │   ├── engine.go
│   │   ├── scheduler.go
│   │   └── scan.go
│   │
│   ├── pluginhost/
│   │   ├── host.go
│   │   ├── process.go
│   │   ├── rpc.go
│   │   └── discovery.go
│   │
│   ├── report/
│   │   ├── json.go
│   │   ├── sarif.go
│   │   └── terminal.go
│   │
│   └── config/
│       └── config.go
│
├── pkg/
│   ├── sdk/
│   │   ├── sdk.go
│   │   ├── plugin.go
│   │   ├── check.go
│   │   ├── finding.go
│   │   └── rpc_server.go
│   │
│   ├── rpc/
│   │   ├── messages.go
│   │   └── protocol.go
│   │
│   └── types/
│       └── severity.go
│
├── plugins/
│   └── plugin-os/
│       ├── main.go
│       └── ssh_check.go
│
└── examples/
    └── plugin-template/
        └── main.go
```

---

# 2. go.mod

```go
module github.com/your-org/ztvs

go 1.24
```

---

# 3. CLI Entrypoint

`cmd/zt/main.go`

```go
package main

import "github.com/your-org/ztvs/internal/cli"

func main() {
    cli.Execute()
}
```

---

# 4. CLI Bootstrap

`internal/cli/root.go`

```go
package cli

import (
    "fmt"
    "os"

    "github.com/your-org/ztvs/internal/engine"
)

func Execute() {
    if len(os.Args) < 2 {
        fmt.Println("usage: zt scan")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "scan":
        eng := engine.New()
        if err := eng.Scan(); err != nil {
            fmt.Println("scan failed:", err)
            os.Exit(1)
        }
    default:
        fmt.Println("unknown command")
        os.Exit(1)
    }
}
```

---

# 5. Engine

`internal/engine/engine.go`

```go
package engine

import (
    "context"

    "github.com/your-org/ztvs/internal/pluginhost"
)

type Engine struct {
    host *pluginhost.Host
}

func New() *Engine {
    return &Engine{
        host: pluginhost.New(),
    }
}

func (e *Engine) Scan() error {
    ctx := context.Background()

    plugins, err := e.host.Discover(ctx)
    if err != nil {
        return err
    }

    for _, p := range plugins {
        _, err := e.host.RunCheck(ctx, p, "ssh_config")
        if err != nil {
            return err
        }
    }

    return nil
}
```

---

# 6. Plugin Host

`internal/pluginhost/host.go`

```go
package pluginhost

import (
    "context"
)

type Host struct{}

func New() *Host {
    return &Host{}
}

func (h *Host) Discover(ctx context.Context) ([]string, error) {
    return []string{
        "./plugins/plugin-os/plugin-os",
    }, nil
}
```

---

# 7. Plugin Process Runner

`internal/pluginhost/process.go`

```go
package pluginhost

import (
    "bytes"
    "context"
    "encoding/json"
    "os/exec"

    "github.com/your-org/ztvs/pkg/rpc"
)

func (h *Host) RunCheck(
    ctx context.Context,
    pluginPath string,
    checkID string,
) (*rpc.RunCheckResponse, error) {

    req := rpc.Request{
        JSONRPC: "2.0",
        ID:      "1",
        Method:  "run_check",
        Params: rpc.RunCheckRequest{
            CheckID: checkID,
        },
    }

    payload, _ := json.Marshal(req)

    cmd := exec.CommandContext(
        ctx,
        pluginPath,
        "--rpc",
    )

    cmd.Stdin = bytes.NewReader(payload)

    out, err := cmd.Output()
    if err != nil {
        return nil, err
    }

    var resp rpc.Response[rpc.RunCheckResponse]
    err = json.Unmarshal(out, &resp)

    return &resp.Result, err
}
```

---

# 8. RPC Protocol Types

`pkg/rpc/messages.go`

```go
package rpc

type Request struct {
    JSONRPC string      `json:"jsonrpc"`
    ID      string      `json:"id"`
    Method  string      `json:"method"`
    Params  interface{} `json:"params"`
}

type Response[T any] struct {
    JSONRPC string `json:"jsonrpc"`
    ID      string `json:"id"`
    Result  T      `json:"result"`
}

type RunCheckRequest struct {
    CheckID string `json:"check_id"`
}

type RunCheckResponse struct {
    Status  string   `json:"status"`
    Finding *Finding `json:"finding"`
}

type Finding struct {
    ID          string                 `json:"id"`
    Severity    string                 `json:"severity"`
    Title       string                 `json:"title"`
    Description string                 `json:"description"`
    Evidence    map[string]interface{} `json:"evidence"`
    Remediation string                 `json:"remediation"`
}
```

---

# 9. SDK Interfaces

`pkg/sdk/check.go`

```go
package sdk

import "context"

type Check interface {
    ID() string
    Name() string
    Run(ctx context.Context) (*Finding, error)
}
```

---

`pkg/sdk/finding.go`

```go
package sdk

type Finding struct {
    ID          string
    Severity    string
    Title       string
    Description string
    Evidence    map[string]interface{}
    Remediation string
}
```

---

# 10. SDK Runtime

`pkg/sdk/sdk.go`

```go
package sdk

import (
    "context"
    "encoding/json"
    "os"

    "github.com/your-org/ztvs/pkg/rpc"
)

func Run(checks []Check) {
    var req rpc.Request
    _ = json.NewDecoder(os.Stdin).Decode(&req)

    switch req.Method {
    case "run_check":
        params := map[string]interface{}{}
        b, _ := json.Marshal(req.Params)
        _ = json.Unmarshal(b, &params)

        checkID := params["check_id"].(string)

        for _, c := range checks {
            if c.ID() == checkID {
                finding, _ := c.Run(context.Background())

                resp := rpc.Response[rpc.RunCheckResponse]{
                    JSONRPC: "2.0",
                    ID:      req.ID,
                    Result: rpc.RunCheckResponse{
                        Status: "fail",
                        Finding: &rpc.Finding{
                            ID: finding.ID,
                            Severity: finding.Severity,
                            Title: finding.Title,
                            Description: finding.Description,
                            Evidence: finding.Evidence,
                            Remediation: finding.Remediation,
                        },
                    },
                }

                _ = json.NewEncoder(os.Stdout).Encode(resp)
                return
            }
        }
    }
}
```

---

# 11. First Plugin

`plugins/plugin-os/main.go`

```go
package main

import "github.com/your-org/ztvs/pkg/sdk"

func main() {
    sdk.Run([]sdk.Check{
        &SSHCheck{},
    })
}
```

---

# 12. Example Check

`plugins/plugin-os/ssh_check.go`

```go
package main

import (
    "context"

    "github.com/your-org/ztvs/pkg/sdk"
)

type SSHCheck struct{}

func (c *SSHCheck) ID() string {
    return "ssh_config"
}

func (c *SSHCheck) Name() string {
    return "SSH Config Check"
}

func (c *SSHCheck) Run(
    ctx context.Context,
) (*sdk.Finding, error) {

    return &sdk.Finding{
        ID: "F-001",
        Severity: "high",
        Title: "Root login enabled",
        Description: "PermitRootLogin yes found",
        Evidence: map[string]interface{}{
            "file": "/etc/ssh/sshd_config",
            "value": "PermitRootLogin yes",
        },
        Remediation: "Set PermitRootLogin no",
    }, nil
}
```

---

# 13. Build Commands

```bash
go build -o zt ./cmd/zt
go build -o plugin-os ./plugins/plugin-os
```

---

# 14. Example Run

```bash
./zt scan
```

Expected output (next implementation step):

```text
HIGH: Root login enabled
```

---

# 15. Immediate Next Engineering Tasks

Priority order:

### P0

* handshake RPC
* enumerate RPC
* timeout support
* stderr logging
* plugin discovery directories

### P1

* concurrency worker pool
* JSON reporter
* terminal reporter
* config loader

### P2

* SARIF output
* signed plugin manifest
* capability enforcement
* policy engine

---

# 16. Strong Recommendation

This skeleton is enough for an engineering team to start building immediately.

The **core contracts are now stable**:

* SDK
* wire protocol
* plugin process model
* CLI
* engine bootstrap

This is a strong MVP foundation.
