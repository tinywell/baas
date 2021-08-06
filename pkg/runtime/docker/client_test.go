package docker

import (
	"context"
	"testing"

	"baas/pkg/runtime"

	"github.com/docker/docker/client"
)

func TestClient_Run(t *testing.T) {
	type fields struct {
		opts options
		dcli *client.Client
	}
	type args struct {
		ctx  context.Context
		data runtime.ServiceMetadata
	}
	dcli, _ := newDockerClient(dockerConfig{})
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "testcontainer",
			fields: fields{
				dcli: dcli,
			},
			args: args{
				ctx:  context.Background(),
				data: testData(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				opts: tt.fields.opts,
				dcli: tt.fields.dcli,
			}
			if err := c.Run(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Client.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		{
			name: "testcontainer",
			// args:    args{opts: []Option{WithConfig("/var/run/docker.sock", "")}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func testData() runtime.ServiceMetadata {
	data := NewSingleServiceData()
	data.Image = "mysql:5.7.27"
	data.Name = "testmysql"
	data.ENVs = []string{
		"MYSQL_ROOT_PASSWORD=baas",
		"MYSQL_DATABASE=baas",
		"MYSQL_USER=baas",
		"MYSQL_PASSWORD=baas"}
	data.Ports = []string{"3306:3306"}
	data.Volumes = []string{
		"/tmp/baas/test/:/etc/mysql/conf.d/",
		"/Users/zfh/Documents/WORK/workspace/esfe/db/scripts/dbtrc.sql:/docker-entrypoint-initdb.d/dbtrc.sql",
		"/tmp/baas/test/data:/var/lib/mysql",
	}
	data.CMDs = []string{
		"pwd",
	}

	return data
}

func TestClient_checkNetwork(t *testing.T) {
	type fields struct {
		opts options
		dcli *client.Client
	}
	type args struct {
		ctx context.Context
		net string
	}
	dcli, _ := newDockerClient(dockerConfig{})
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "testnetwork",
			fields: fields{
				dcli: dcli,
			},
			args: args{
				ctx: context.Background(),
				net: "testbaasnet",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				opts: tt.fields.opts,
				dcli: tt.fields.dcli,
			}
			if err := c.checkNetwork(tt.args.ctx, tt.args.net); (err != nil) != tt.wantErr {
				t.Errorf("Client.checkNetwork() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
