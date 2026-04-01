package main

import (
	"github.com/mosesgameli/ztvs-sdk-go/sdk"
)

func main() {
	meta := sdk.Metadata{
		Name:       "Axios Mitigation Plugin",
		Version:    "1.0.0",
		APIVersion: 1,
	}

	checks := []sdk.Check{
		&DependencyCheck{},
		&HostIOCCheck{},
		&NetworkCheck{},
	}

	sdk.Run(meta, checks)
}
