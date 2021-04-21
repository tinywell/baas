package helm3

import (
	"context"
	"testing"

	"github.com/tinywell/baas/pkg/runtime"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
)

func TestClient_Run(t *testing.T) {
	type fields struct {
		chartOpts *action.ChartPathOptions
		settings  *cli.EnvSettings
		options   *options
	}
	type args struct {
		ctx  context.Context
		data runtime.ServiceMetadata
	}
	setting := cli.New()
	setting.KubeConfig = "/Users/zfh/.kube/config"
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "testinstall",
			fields: fields{
				chartOpts: &action.ChartPathOptions{
					InsecureSkipTLSverify: false,
					RepoURL:               "https://apphub.aliyuncs.com",
				},
				settings: setting,
				options: &options{
					kubeCfg: kubeConfig{
						File: "/Users/zfh/.kube/config",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				data: &MetaData{
					Meta: Meta{
						Name: "testbaas",
					},
					Chart:       "nginx",
					Values:      map[string]interface{}{},
					ReleaseName: "testbaas",
					Namespace:   "default",
					svcType:     TypeServiceInstall,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				chartOpts: tt.fields.chartOpts,
				settings:  tt.fields.settings,
				options:   tt.fields.options,
			}
			if err := c.Run(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Client.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
