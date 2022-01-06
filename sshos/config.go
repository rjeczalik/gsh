package sshos

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

func UserConfig(username string) (*os.File, error) {
	usr, err := (*user.User)(nil), error(nil)

	if username == "" {
		usr, err = user.Current()
	} else {
		usr, err = user.Lookup(username)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to lookup user: %w", err)
	}

	dir := filepath.Join(usr.HomeDir, ".ssh")

	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create %q directory: %w", dir, err)
	}

	file := filepath.Join(dir, "config")

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open or create %q: %w", file, err)
	}

	return f, nil
}
