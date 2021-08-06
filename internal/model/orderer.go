package model

// Orderer ..
type Orderer struct {
	BaaSData
	HFNode
	Name       string `json:"name,omitempty" db:"name"`
	DomainName string
	Endpoint   string
	Port       int
	Image      string
	Tag        string
}
