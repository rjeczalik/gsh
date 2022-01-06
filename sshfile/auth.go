package sshfile

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh"
)

var NoAuthMethods = errors.New("no auth methods could be used")

func IdentityAuth(files ...string) (ssh.AuthMethod, error) {
	var signers []ssh.Signer

	for _, file := range files {
		p, err := ioutil.ReadFile(file)
		if os.IsNotExist(err) && len(files) != 1 {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("error reading %q file: %w", file, err)
		}

		signer, err := ssh.ParsePrivateKey(p)
		if err != nil {
			return nil, fmt.Errorf("error parsing %q file: %w", file, err)
		}

		signers = append(signers, signer)
	}

	if len(signers) == 0 {
		return nil, NoAuthMethods
	}

	return ssh.PublicKeys(signers...), nil
}
