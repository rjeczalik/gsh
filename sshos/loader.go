package sshos

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/rjeczalik/gsh"
	"github.com/rjeczalik/gsh/sshfile"
	"github.com/rjeczalik/gsh/sshtrace"
	"github.com/rjeczalik/gsh/sshutil"

	xssh "golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func NewClient() (*gsh.Client, error) {
	return DefaultLoader.NewClient()
}

type Loader struct {
	Dir              string
	UserConfig       string
	UserKnownHosts   string
	SystemConfig     string
	SystemKnownHosts string
	Identity         []string
	Options          []string
}

func (l *Loader) NewClient() (*gsh.Client, error) {
	var mixin *sshfile.Config

	if len(l.options()) != 0 {
		var err error
		if mixin, err = sshfile.ParseOptions(l.options()); err != nil {
			return nil, fmt.Errorf("failed to parse options: %w", err)
		}
	}

	usr, err := sshfile.ParseConfigFile(l.userConfig())
	if err != nil && !is(err, os.ErrNotExist, os.ErrPermission) {
		return nil, fmt.Errorf("failed to parse %q user config: %w", l.userConfig(), err)
	}

	sys, err := sshfile.ParseConfigFile(l.systemConfig())
	if err != nil && !is(err, os.ErrNotExist, os.ErrPermission) {
		return nil, fmt.Errorf("failed to parse %q system config: %w", l.systemConfig(), err)
	}

	cfgfile := usr.Merge(sys)

	if mixin != nil {
		for i := range cfgfile {
			if err := cfgfile[i].Merge(mixin); err != nil {
				return nil, fmt.Errorf("%d: unable to apply custom options: %w", i, err)
			}
		}
	}

	cb := cfgfile.Callback()

	auth, err := sshfile.IdentityAuth(l.identity()...)
	if err != nil && !is(err, sshfile.NoAuthMethods) {
		return nil, fmt.Errorf("failed to build identity auth: %w", err)
	}

	if auth != nil {
		cb = sshutil.PatchCallback(cb, func(_ context.Context, cfg *gsh.Config) error {
			cfg.Auth = append(cfg.Auth, auth)
			return nil
		})
	}

	known, err := knownhosts.New(l.userKnownHosts(), l.systemKnownHosts())
	if err != nil && !is(err, os.ErrNotExist, os.ErrPermission) {
		return nil, fmt.Errorf("failed to parse %q, %q known hosts files: %w", l.userKnownHosts(), l.systemKnownHosts(), err)
	}

	if known != nil {
		cb = sshutil.PatchCallback(cb, func(_ context.Context, cfg *gsh.Config) error {
			if cfg.HostKeyCallback == nil {
				cfg.HostKeyCallback = known
			}
			return nil
		})
	} else {
		cb = sshutil.PatchCallback(cb, func(_ context.Context, cfg *gsh.Config) error {
			if cfg.HostKeyCallback == nil {
				cfg.HostKeyCallback = xssh.InsecureIgnoreHostKey()
			}
			return nil
		})
	}

	var once sync.Once

	cb = sshutil.PatchCallback(cb, func(ctx context.Context, cfg *gsh.Config) error {
		if ct := sshtrace.ContextClientTrace(ctx); ct != nil {
			once.Do(func() { ct.GotFileConfig(cfgfile) })
			ct.GotConfig(cfg)
		}
		return nil
	})

	c := &gsh.Client{
		ConfigCallback: cb,
		DialContext:    sshutil.DialContext,
	}

	return c, nil
}

func (l *Loader) copy() *Loader {
	lCopy := *l

	if len(l.Options) != 0 {
		lCopy.Options = make([]string, len(l.Options))
		copy(lCopy.Options, l.Options)
	}

	return &lCopy
}

func (l *Loader) dir() string {
	if l.Dir != "" {
		return l.Dir
	}
	return DefaultLoader.Dir
}

func (l *Loader) userConfig() string {
	if l.UserConfig != "" {
		return l.UserConfig
	}
	return DefaultLoader.UserConfig
}

func (l *Loader) userKnownHosts() string {
	if l.UserKnownHosts != "" {
		return l.UserKnownHosts
	}
	return DefaultLoader.UserKnownHosts
}

func (l *Loader) systemConfig() string {
	if l.SystemConfig != "" {
		return l.SystemConfig
	}
	return DefaultLoader.SystemConfig
}

func (l *Loader) systemKnownHosts() string {
	if l.SystemKnownHosts != "" {
		return l.SystemKnownHosts
	}
	return DefaultLoader.SystemKnownHosts
}

func (l *Loader) identity() []string {
	if len(l.Identity) != 0 {
		return l.Identity
	}
	return DefaultLoader.Identity
}

func (l *Loader) options() []string {
	if len(l.Options) != 0 {
		return l.Options
	}
	return DefaultLoader.Options
}
