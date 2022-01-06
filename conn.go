package gsh

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type SessionFunc func(context.Context, *ssh.Session) error

func (fn SessionFunc) Do(ctx context.Context, s *ssh.Session) error {
	return fn(ctx, s)
}

type Session interface {
	Do(context.Context, *ssh.Session) error
}

type Conn interface {
	Session(Session) error
	Context() context.Context
	Close() error
}

var _ Conn = (*conn)(nil)

type conn struct {
	ctx  context.Context
	conn ssh.Conn
	cli  *ssh.Client
}

func (c *conn) Session(s Session) error {
	sshSession, err := c.cli.NewSession()
	if err != nil {
		return fmt.Errorf("session error: %w", err)
	}
	defer sshSession.Close()

	return s.Do(c.Context(), sshSession)
}

func (c *conn) Context() context.Context {
	return c.ctx
}

func (c *conn) Close() error {
	return c.conn.Close()
}

func Command(cmd string, args ...string) *Cmd {
	for i, arg := range args {
		args[i] = strconv.Quote(arg)
	}

	return &Cmd{
		Path:   cmd,
		Args:   args,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

type Cmd struct {
	Path   string
	Args   []string
	Env    []string
	Dir    string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (c *Cmd) String() string {
	return strings.TrimSpace(c.Path + " " + strings.Join(c.Args, " "))
}

func (c *Cmd) Do(ctx context.Context, s *ssh.Session) error {
	s.Stdin = c.Stdin
	s.Stdout = c.Stdout
	s.Stderr = c.Stderr

	return s.Run(c.String())
}

func (c *Cmd) Shell() *Shell {
	return &Shell{Cmd: c}
}

type Shell struct {
	*Cmd
}

func (sh *Shell) Do(ctx context.Context, s *ssh.Session) error {
	s.Stdin = sh.Stdin
	s.Stdout = sh.Stdout
	s.Stderr = sh.Stderr

	width, height := 80, 25
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if stdin, ok := sh.Stdin.(*os.File); ok && terminal.IsTerminal(int(stdin.Fd())) {
		fd := int(stdin.Fd())

		if w, h, err := terminal.GetSize(fd); err == nil {
			width, height = w, h
		}

		if state, err := terminal.MakeRaw(fd); err == nil {
			defer terminal.Restore(fd, state)
		}
	}

	if err := s.RequestPty("xterm", height, width, modes); err != nil {
		return fmt.Errorf("request pty error: %w", err)
	}

	return s.Run(sh.Cmd.String())
}
