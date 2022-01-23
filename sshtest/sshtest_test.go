package sshtest_test

import (
	"io/ioutil"
	"testing"

	"github.com/rjeczalik/gsh/sshtest"
)

func TestFixture(t *testing.T) {
	p, err := ioutil.ReadFile("")
	if err != nil {
		t.Fatalf("failed reading: %s", err)
	}

	_ = p // raw private key

	f, err := sshtest.NewFixture()
	if err != nil {
		t.Fatalf("fixture error: %s", err)
	}
	defer f.Close()

	box, err := f.RunBox("sshtest-ubuntu", "latest")
	if err != nil {
		t.Fatalf("box error: %s", err)
	}

	err = box.Pool.Retry(func() error {
		// todo: ssh connect to the box.Address()
		return nil
	})
	if err != nil {
		t.Fatalf("connection error: %s", err)
	}

	// todo: rest of the test
}
