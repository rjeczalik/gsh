package sshutil

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"time"

	"github.com/rjeczalik/gsh"
)

func init() {
	rand.Seed(time.Now().UnixMicro())
}

type Dialer struct {
	Dialer   *net.Dialer
	MaxCount int
}

func DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	cfg := ctx.Value(gsh.ConfigKey).(*gsh.Config)

	d := &Dialer{
		Dialer: &net.Dialer{
			Timeout: cfg.Connection.Timeout,
		},
		MaxCount: cfg.Connection.MaxCount,
	}

	return d.DialContext(ctx, network, addr)
}

func (d *Dialer) DialContext(ctx context.Context, network, addr string) (c net.Conn, err error) {
	type temporary interface {
		Temporary() bool
	}

	for n := d.MaxCount; n > 0; n-- {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		c, err = d.Dialer.DialContext(ctx, network, addr)
		if err == nil {
			break
		}
		if tmp := temporary(nil); errors.As(err, &tmp) || !tmp.Temporary() {
			break
		}

		time.Sleep(time.Duration(rand.Int31n(230)+20) * time.Millisecond)
	}

	return c, err
}
