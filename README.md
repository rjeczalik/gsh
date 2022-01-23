# gsh

`tectum.sh/gsh` is yet another wrapper for the golang.org/x/crypto/ssh package. `gsh` is both a Go client API and a command line tool.

Since it understand OpenSSH's configuration files, it can be used as a drop-in replacement for ssh for simple use-case (it is not feature-complete yet).

### Getting started with the command line tool

The quickest way to give the `gsh` a try is to use provided Docker images. Once the container is built and started, it will serve as OpenSSH server that the `gsh` can connect to.

In order to build and start the container, run the following commands:

```bash
./sshtest/build.sh
```
```bash
docker run --rm -d -ti -p 2222:22 -P sshtest-ubuntu:latest
```

Then build the `gsh` command line tool and connect to the container:

```bash
go build github.com/rjeczalik/gsh/cmd/gsh
```
```bash
gsh -F sshtest/ssh.config sshtest
```
```
To run a command as administrator (user "root"), use "sudo <command>".
See "man sudo_root" for details.

sshtest@9ddf71ccfd81:~$ 
```

### Getting started with the client API

Equivalent code for the above `gsh` command looks like the following:

```go
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
```

Give it a try by running:

```bash
go run examples/sshtest.go
```
```
To run a command as administrator (user "root"), use "sudo <command>".
See "man sudo_root" for details.

sshtest@9ddf71ccfd81:~$ 
```
