package docker

// MetaData ...
type MetaData struct {
	Name       string
	Image      string
	Ports      []string // ["7050:7050","7051:7051"]
	Volumes    []string
	ENVs       []string // ["MYSQL_USER=baas"]
	CMDs       []string
	ExtraHosts []string // ["127.0.0.1:orderer.example.com",""]
	Network    string
	WorkDir    string
	svcType    string
}

// NewSingleServiceData ...
func NewSingleServiceData() *MetaData {
	return &MetaData{
		svcType: TypeServiceCreateSingle,
	}
}

// ServiceType ...
func (md *MetaData) ServiceType() string {
	return md.svcType
}

// RuntimeType ...
func (md *MetaData) RuntimeType() string {
	return TypeRuntimeDocker
}

// DataID ...
func (md *MetaData) DataID() string {
	return md.Name
}
