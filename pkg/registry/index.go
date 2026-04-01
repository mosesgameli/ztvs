package registry

type PluginMetadata struct {
	Name          string `json:"name"`
	LatestVersion string `json:"latest_version"`
	Repo          string `json:"repo"`
	Checksum      string `json:"checksum"`
	Signature     string `json:"signature"`
	AuditStatus   string `json:"audit_status"`
}

type Index struct {
	Version string           `json:"version"`
	Plugins []PluginMetadata `json:"plugins"`
}
