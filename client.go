package gsh

import (
	"context"
	"errors"
	"fmt"
	"net"

	"golang.org/x/crypto/ssh"
)

var ErrConfigNotFound = errors.New("config not found")

type contextKey struct {
	string
}

var (
	ConfigKey = contextKey{"ssh-config"}
	ConnKey   = contextKey{"ssh-conn"}
	ClientKey = contextKey{"ssh-client"}
)

type (
	ConfigCallback func(ctx context.Context, network, addr string) (*Config, error)
	DialContext    func(ctx context.Context, network, addr string) (net.Conn, error)
)

type Client struct {
	ConfigCallback
	DialContext
}

func (c *Client) Connect(ctx context.Context, network, address string) (Conn, error) {
	cfg, err := c.ConfigCallback(ctx, network, address)
	if err != nil {
		return nil, fmt.Errorf("config callback error: %w", err)
	}

	ctx = context.WithValue(ctx, ConfigKey, cfg)

	tcpConn, err := c.DialContext(ctx, cfg.Network, cfg.Address)
	if err != nil {
		return nil, fmt.Errorf("dial error: %w", err)
	}

	sshConn, chans, reqs, err := ssh.NewClientConn(tcpConn, cfg.Address, &cfg.ClientConfig)
	if err != nil {
		return nil, fmt.Errorf("ssh connection error: %w", err)
	}

	ctx = context.WithValue(ctx, ConnKey, sshConn)

	sshCli := ssh.NewClient(sshConn, chans, reqs)

	ctx = context.WithValue(ctx, ClientKey, sshCli)

	return &conn{
		ctx:  ctx,
		conn: sshConn,
		cli:  sshCli,
	}, nil
}
