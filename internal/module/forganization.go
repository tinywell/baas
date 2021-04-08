package module

// FOrganization fabric 组织
type FOrganization struct {
	BaaSData
	Name      string
	MSPID     string
	CACert    string
	CAKey     string
	TLSCACert string
	TLSCAKey  string
	AdminCert string
	Domian    string
}
