module github.com/mosesgameli/ztvs/plugins/plugin-os

go 1.26.1

replace (
	github.com/mosesgameli/ztvs => ../..
	github.com/mosesgameli/ztvs-sdk-go => ../../../sdk/go
)

require github.com/mosesgameli/ztvs-sdk-go v0.0.0-20260401141939-b860acbcb67c
