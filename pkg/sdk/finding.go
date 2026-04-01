package sdk

type Finding struct {
	ID          string
	Severity    string
	Title       string
	Description string
	Evidence    map[string]interface{}
	Remediation string
}
