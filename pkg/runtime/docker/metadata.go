package docker

// Meta ...
type Meta struct {
	Name string
}

// RuntimeType 运行时类型
func (r *Meta) RuntimeType() string {
	return TypeRuntimeDocker
}

// DataID 数据标识
func (r *Meta) DataID() string {
	return r.Name
}

// MetaData 资源创建用数据类型
type MetaData struct {
	Meta
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
