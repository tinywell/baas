package docker

// runtime type
const (
	TypeRuntimeDocker = "DOCKER"

	TypeServiceCreateSingle = "CreateSingle"
)

// WithConfig ...
func WithConfig(host, version string) Option {
	return func(opts *options) error {
		opts.dockerConfig.host = host
		opts.dockerConfig.version = version
		return nil
	}
}

// WithTLS ...
func WithTLS(capath, certpath, keypath string) Option {
	return func(opts *options) error {
		opts.dockerConfig.tls = true
		opts.dockerConfig.capath = capath
		opts.dockerConfig.certpath = certpath
		opts.dockerConfig.keypath = keypath
		return nil
	}
}

// NotDockerRuntimeErr ...
type NotDockerRuntimeErr struct{}

func (e NotDockerRuntimeErr) Error() string {
	return "not for runtime: " + TypeRuntimeDocker
}
