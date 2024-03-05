package builder

import (
	"strings"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/common/log"
)

type ContainerBuilder struct {
	opts types.CreateDockerContainerOptions
}

func NewContainerOpts() *ContainerBuilder {
	return &ContainerBuilder{}
}

func (b *ContainerBuilder) Build() types.CreateDockerContainerOptions {
	return b.opts
}

func (b *ContainerBuilder) WithName(name string) *ContainerBuilder {
	b.opts.ContainerName = name
	return b
}

func (b *ContainerBuilder) WithImage(image string) *ContainerBuilder {
	b.opts.ImageName = image
	return b
}

func (b *ContainerBuilder) WithCommand(cmd *string) *ContainerBuilder {
	if cmd != nil {
		b.opts.Cmd = strings.Split(*cmd, " ")
	}
	return b
}

func (b *ContainerBuilder) WithEnv(env types.EnvVariables) *ContainerBuilder {
	for _, e := range env {
		b.opts.Env = append(b.opts.Env, e.Name+"="+e.Value)
	}
	return b
}

func (b *ContainerBuilder) WithCaps(caps types.Capabilities) *ContainerBuilder {
	for _, cp := range caps {
		b.opts.CapAdd = append(b.opts.CapAdd, cp.Name)
	}
	return b
}

func (b *ContainerBuilder) WithSysctls(sysctls types.Sysctls) *ContainerBuilder {
	if len(sysctls) == 0 {
		return b
	}
	b.opts.Sysctls = make(map[string]string)
	for _, sysctl := range sysctls {
		b.opts.Sysctls[sysctl.Name] = sysctl.Value
	}
	return b
}

func (b *ContainerBuilder) WithPorts(ports types.Ports) *ContainerBuilder {
	var all []string
	for _, p := range ports {
		out := p.Out
		in := p.In
		all = append(all, out+":"+in)
	}

	var err error
	b.opts.ExposedPorts, b.opts.PortBindings, err = nat.ParsePortSpecs(all)
	if err != nil {
		log.Error(err)
	}

	return b
}

func (b *ContainerBuilder) WithVolumes(volumes types.Volumes) *ContainerBuilder {
	for _, v := range volumes {
		if v.Type == types.VolumeTypeBind {
			b.opts.Binds = append(b.opts.Binds, v.Out+":"+v.In)
		} else {
			b.opts.Mounts = append(b.opts.Mounts, mount.Mount{
				Type:   mount.TypeVolume,
				Source: v.Out,
				Target: v.In,
			})
		}
	}
	return b
}
