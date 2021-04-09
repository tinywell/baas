package helm3

// runtime type
const (
	TypeRuntimeHelm3 = "HELM3"

	TypeServiceInstall = "Install"
)

// Meta ...
type Meta struct {
	Name string
}

// RuntimeType 运行时类型
func (r *Meta) RuntimeType() string {
	return TypeRuntimeHelm3
}

// DataID 数据标识
func (r *Meta) DataID() string {
	return r.Name
}

// MetaData 资源创建用数据类型
type MetaData struct {
	Meta
	Chart       string
	Values      map[string]interface{}
	ReleaseName string
	Namespace   string
	svcType     string
}

// ServiceType ...
func (md *MetaData) ServiceType() string {
	return md.svcType
}

// NewInstallData ...
func NewInstallData() *MetaData {
	return &MetaData{
		svcType: TypeServiceInstall,
	}
}
