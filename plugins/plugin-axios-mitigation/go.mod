module github.com/mosesgameli/ztvs/plugins/plugin-axios-mitigation

go 1.26.1

require github.com/mosesgameli/ztvs-sdk-go v0.0.0-00010101000000-000000000000

replace (
	github.com/mosesgameli/ztvs => ../..
	github.com/mosesgameli/ztvs-sdk-go => ../../../sdk/go
)
