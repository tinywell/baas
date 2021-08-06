package docker

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"

	"baas/pkg/runtime"
)

// Client ...
type Client struct {
	opts options
	dcli *client.Client
}

// Opts ...
type options struct {
	dockerConfig dockerConfig
}

type dockerConfig struct {
	host     string
	version  string
	tls      bool
	capath   string
	certpath string
	keypath  string
}

// Option ...
type Option func(opts *options) error

// NewClient ...
func NewClient(opts ...Option) (*Client, error) {
	cfgs := &options{}
	for _, opt := range opts {
		err := opt(cfgs)
		if err != nil {
			return nil, err
		}
	}
	client := &Client{
		opts: *cfgs,
	}
	dcli, err := newDockerClient(cfgs.dockerConfig)
	if err != nil {
		return nil, err
	}
	client.dcli = dcli
	return client, nil
}

// Run ...
func (c *Client) Run(ctx context.Context, data runtime.ServiceMetadata) error {
	if data.RuntimeType() != TypeRuntimeDocker {
		return NotDockerRuntimeErr{}
	}
	switch data.ServiceType() {
	case TypeServiceCreateSingle:
		if d, ok := data.(*MetaData); ok {
			_, err := c.runSingle(ctx, d)
			if err != nil {
				return err
			}
		} else {
			return errors.New("参数 data 与其类型不匹配")
		}
	default:
		return errors.New("数据类型暂不支持")
	}
	return nil
}

func (c *Client) runSingle(ctx context.Context, data *MetaData) (*container.ContainerCreateCreatedBody, error) {
	containerConfig, err := c.parseContainerConfig(data)
	if err != nil {
		return nil, err
	}
	hostConfig, err := c.parseHostConfig(data)
	if err != nil {
		return nil, err
	}
	err = c.checkNetwork(ctx, data.Network)
	if err != nil {
		return nil, err
	}
	body, err := c.dcli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, data.Name)
	if err != nil {
		return nil, errors.WithMessage(err, "创建容器失败")
	}
	err = c.dcli.ContainerStart(ctx, body.ID, types.ContainerStartOptions{})
	if err != nil {
		return nil, errors.WithMessage(err, "启动容器失败")
	}
	return &body, nil
}

func (c *Client) parseContainerConfig(data *MetaData) (*container.Config, error) {
	config := &container.Config{
		Env:        data.ENVs,
		Cmd:        strslice.StrSlice(data.CMDs),
		Image:      data.Image,
		WorkingDir: data.WorkDir,
		Volumes:    make(map[string]struct{}),
	}
	for _, v := range data.Volumes {
		vs := strings.Split(v, ":")
		if len(vs) != 2 {
			return nil, errors.Errorf("挂载信息格式有误: %s", v)
		}
		d := vs[1]
		config.Volumes[d] = struct{}{}
	}

	ps, _, err := nat.ParsePortSpecs(data.Ports)
	if err != nil {
		return nil, errors.WithMessage(err, "解析端口信息失败")
	}
	config.ExposedPorts = ps

	return config, nil
}

func (c *Client) parseHostConfig(data *MetaData) (*container.HostConfig, error) {
	config := &container.HostConfig{
		NetworkMode: container.NetworkMode(data.Network),
		ExtraHosts:  data.ExtraHosts,
	}
	for _, v := range data.Volumes {
		vs := strings.Split(v, ":")
		if len(vs) != 2 {
			return nil, errors.Errorf("挂载信息格式有误: %s", v)
		}
		// s, d := vs[0], vs[1]
		// mount := mount.Mount{
		// 	Type: mount.TypeBind, Source: s, Target: d,
		// 	BindOptions: &mount.BindOptions{Propagation: mount.PropagationRPrivate},
		// }
		// config.Mounts = append(config.Mounts, mount)
		bind := strings.Join([]string{v, "rw"}, ":")
		config.Binds = append(config.Binds, bind)
	}
	_, pbs, err := nat.ParsePortSpecs(data.Ports)
	if err != nil {
		return nil, errors.WithMessage(err, "解析 Ports 信息失败")
	}
	config.PortBindings = pbs

	return config, nil
}

func (c *Client) checkNetwork(ctx context.Context, net string) error {
	filters := filters.NewArgs(filters.Arg("name", net))
	r, err := c.dcli.NetworkList(ctx, types.NetworkListOptions{Filters: filters})
	if err != nil {
		return err
	}
	if len(r) > 0 {
		return nil
	}
	_, err = c.dcli.NetworkCreate(ctx, net, types.NetworkCreate{})
	if err != nil {
		return err
	}
	return nil
}

func newDockerClient(cfg dockerConfig) (*client.Client, error) {
	opts := []client.Opt{}
	if len(cfg.host) > 0 {
		opts = append(opts, client.WithHost(cfg.host))
	}
	if len(cfg.version) == 0 {
		opts = append(opts, client.WithAPIVersionNegotiation())
	}
	if cfg.tls {
		opts = append(opts, client.WithTLSClientConfig(cfg.capath, cfg.certpath, cfg.keypath))
	}
	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return nil, err
	}
	return cli, nil
}
