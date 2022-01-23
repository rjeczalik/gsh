package main

import (
	"context"
	"log"

	"github.com/rjeczalik/gsh"
	"github.com/rjeczalik/gsh/sshfile"
	"github.com/rjeczalik/gsh/sshutil"
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
