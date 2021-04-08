package module

// Peer ...
type Peer struct {
	BaaSData
	HFNode
	Name       string `json:"name,omitempty" db:"name"`
	DomainName string
	Endpoint   string
	Port       int
	EXTPort    int
	Image      string
	StateDB    string `json:"state_db,omitempty" db:"state_db"`
}
