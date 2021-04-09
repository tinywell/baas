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
