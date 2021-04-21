package helm3

import (
	"context"
	"fmt"
	"log"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"

	"github.com/pkg/errors"
	"github.com/tinywell/baas/pkg/runtime"
)

// Client ...
type Client struct {
	chartOpts *action.ChartPathOptions
	settings  *cli.EnvSettings
	options   *options
}

// NewClient ...
func NewClient(opts ...Option) (*Client, error) {
	options := &options{}
	for _, opt := range opts {
		err := opt(options)
		if err != nil {
			return nil, err
		}
	}
	settings := cli.New()
	if len(options.kubeCfg.File) > 0 {
		settings.KubeConfig = options.kubeCfg.File
	}
	if len(options.kubeCfg.Context) > 0 {
		settings.KubeContext = options.kubeCfg.Context
	}

	client := &Client{
		options:   options,
		settings:  settings,
		chartOpts: options.chartOpts,
	}

	return client, nil
}

// Run ...
func (c *Client) Run(ctx context.Context, data runtime.ServiceMetadata) error {
	if data.RuntimeType() != TypeRuntimeHelm3 {
		return errors.Errorf("运行时类型不匹配，期望：%s 实际：%s", TypeRuntimeHelm3, data.RuntimeType())
	}
	switch data.ServiceType() {
	case TypeServiceInstall:
		if data, ok := data.(*MetaData); ok {
			return c.install(ctx, data)
		}
		return errors.New("数据内容错误")
	default:
		return errors.Errorf("不支持的类型：%s", data.ServiceType())
	}
}

func (c *Client) install(tx context.Context, data *MetaData) error {
	actionConfig := new(action.Configuration)
	if len(data.Namespace) == 0 {
		data.Namespace = c.settings.Namespace()
	}
	if err := actionConfig.Init(c.settings.RESTClientGetter(), data.Namespace, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		fmt.Printf("%+v", err)
		return errors.WithMessage(err, "action config 初始化失败")
	}
	client := action.NewInstall(actionConfig)
	client.ReleaseName = data.ReleaseName
	client.Namespace = data.Namespace
	cp, err := c.chartOpts.LocateChart(data.Chart, c.settings)
	if err != nil {
		return errors.WithMessagef(err, "chart 请求构建失败，chart=%s", data.Chart)
	}
	chartRequested, err := loader.Load(cp)
	if err != nil {
		return errors.WithMessagef(err, "获取 chart 失败,chart=%s", data.Chart)
	}
	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if client.DependencyUpdate {
				man := &downloader.Manager{
					ChartPath:        cp,
					Keyring:          client.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          getter.All(c.settings),
					RepositoryConfig: c.settings.RepositoryConfig,
					RepositoryCache:  c.settings.RepositoryCache,
				}
				if err := man.Update(); err != nil {
					return errors.WithMessage(err, "更新 chart 失败")
				}
			} else {
				return errors.WithMessage(err, "校验 chart 依赖失败")
			}
		}
	}

	// vals, err := mergeValues(data.Values)
	// if err != nil {
	// 	return err
	// }
	r, err := client.Run(chartRequested, data.Values)
	if err != nil {
		return err
	}
	fmt.Printf("release: %s\n", r.Name)
	return nil
}

func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = mergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}

// mergeValues ...
func mergeValues(values map[string]interface{}, oldConfigs ...map[string]interface{}) (map[string]interface{}, error) {
	base := map[string]interface{}{}

	if len(oldConfigs) != 0 {
		for _, cur := range oldConfigs {
			base = mergeMaps(base, cur)
		}
	}
	base = mergeMaps(base, values)
	return base, nil
}

// isChartInstallable validates if a chart can be installed
//
// Application chart type is only installable
func isChartInstallable(ch *chart.Chart) (bool, error) {
	switch ch.Metadata.Type {
	case "", "application":
		return true, nil
	}
	return false, errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}
