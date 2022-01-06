package sshtest

import (
	"fmt"
	"io"
	"net"

	dt "github.com/ory/dockertest"
	dc "github.com/ory/dockertest/docker"
	"github.com/spf13/pflag"
)

type Fixture struct {
	Pool *dt.Pool
	Keep bool

	closers []io.Closer
}

func (f *Fixture) Register(fs *pflag.FlagSet) {
	fs.BoolVarP(&f.Keep, "keep-box", "k", false, "Do not auto-remove containers")
}

func (f *Fixture) Close() error {
	var errs []error
	for _, c := range f.closers {
		if err := c.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return fmt.Errorf("close failed: %v", errs)
	}

	return nil
}

// image = repo:tag
func (f *Fixture) RunBox(repo, tag string) (*Box, error) {
	res, err := f.Pool.RunWithOptions(
		&dt.RunOptions{
			Repository: repo,
			Tag:        tag,
		},
		func(hc *dc.HostConfig) {
			hc.AutoRemove = !f.Keep
		},
	)
	if err != nil {
		return nil, fmt.Errorf("container execution failed: %w", err)
	}

	b := &Box{
		Resource: res,
	}

	f.closers = append(f.closers, b)

	return b, nil
}

type Box struct {
	Pool     *dt.Pool
	Resource *dt.Resource
}

func (b *Box) Address() string {
	return net.JoinHostPort("localhost", b.Resource.GetPort("22/tcp"))
}

func (b *Box) Close() error {
	// todo: b.Resource.Expire?
	return b.Pool.Purge(b.Resource)
}

func NewFixture() (*Fixture, error) {
	pool, err := dt.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("docker connection failed: %w", err)
	}

	return &Fixture{
		Pool: pool,
	}, nil
}
