package sshutil

import (
	"context"
	"errors"
	"fmt"

	"github.com/tectumsh/gsh"
)

func Callback(callbacks ...gsh.ConfigCallback) gsh.ConfigCallback {
	return func(ctx context.Context, network, address string) (*gsh.Config, error) {
		for _, cb := range callbacks {
			cfg, err := cb(ctx, network, address)
			if err == nil {
				return cfg, nil
			}
			if errors.Is(err, gsh.ErrConfigNotFound) {
				continue
			}
			return nil, err
		}
		return nil, gsh.ErrConfigNotFound
	}
}

func LazyCallback(fns ...func() gsh.ConfigCallback) gsh.ConfigCallback {
	return func(ctx context.Context, network, address string) (*gsh.Config, error) {
		var callbacks []gsh.ConfigCallback

		for _, fn := range fns {
			callbacks = append(callbacks, fn())
		}

		return Callback(callbacks...)(ctx, network, address)
	}
}

func PatchCallback(cb gsh.ConfigCallback, patch func(context.Context, *gsh.Config) error) gsh.ConfigCallback {
	return func(ctx context.Context, network, address string) (*gsh.Config, error) {
		cfg, err := cb(ctx, network, address)
		if err != nil {
			return nil, err
		}

		if err := patch(ctx, cfg); err != nil {
			return nil, fmt.Errorf("failed to patch config: %w", err)
		}

		return cfg, nil
	}
}
