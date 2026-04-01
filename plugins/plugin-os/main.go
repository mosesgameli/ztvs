package main

import (
	"github.com/mosesgameli/ztvs-sdk-go/sdk"
)

func main() {
	sdk.Run(sdk.Metadata{
		Name:       "plugin-os",
		Version:    "1.0.0",
		APIVersion: 1,
	}, []sdk.Check{
		&SSHCheck{},
	})
}
