package helm3

import "helm.sh/helm/v3/pkg/action"

type options struct {
	kubeCfg   kubeConfig
	chartOpts *action.ChartPathOptions
}

// Option ...
type Option func(opts *options) error

type kubeConfig struct {
	File    string
	Context string
}

// RepoConfig ...
type RepoConfig struct {
	RepoURL    string
	TLS        bool
	CACertFile string
	CertFile   string
	KeyFile    string
	Private    bool
	Username   string
	Password   string
}

// WithKubeConfig ...
func WithKubeConfig(kubefile, context string) Option {
	return func(opts *options) error {
		opts.kubeCfg.File = kubefile
		opts.kubeCfg.Context = context
		return nil
	}
}

// WithRepoConfig ...
func WithRepoConfig(config RepoConfig) Option {
	return func(opts *options) error {
		opts.chartOpts = &action.ChartPathOptions{}
		opts.chartOpts.RepoURL = config.RepoURL
		if config.TLS {
			opts.chartOpts.CaFile = config.CACertFile
			opts.chartOpts.CertFile = config.CertFile
			opts.chartOpts.KeyFile = config.KeyFile
		}
		if config.Private {
			opts.chartOpts.Username = config.Username
			opts.chartOpts.Password = config.Password
		}
		return nil
	}
}
