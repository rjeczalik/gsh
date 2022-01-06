package main

import (
	"context"
	"log"

	"github.com/tectumsh/gsh"
	"github.com/tectumsh/gsh/sshfile"
	"github.com/tectumsh/gsh/sshutil"
)

func main() {
	ctx := context.Background()

	cfg := &sshfile.Config{
		User:         "sshtest",
		Hostname:     "127.0.0.1",
		Port:         2222,
		IdentityFile: "sshtest/key.pem",
	}

	client := &gsh.Client{
		ConfigCallback: cfg.Callback(),
		DialContext:    sshutil.DialContext,
	}

	conn, err := client.Connect(ctx, "tcp", "")
	if err != nil {
		log.Fatal(err)
	}

	shell := gsh.Command("/bin/bash").Shell()

	if err := conn.Session(shell); err != nil {
		log.Fatal(err)
	}
}
