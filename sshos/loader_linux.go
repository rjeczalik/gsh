// +build linux

package sshos

import (
	"os/user"
	"path/filepath"
)

var DefaultLoader = &Loader{
	Dir:              filepath.Join(home, ".ssh", "config"),
	UserConfig:       filepath.Join(home, ".ssh", "config"),
	UserKnownHosts:   filepath.Join(home, ".ssh", "known_hosts"),
	SystemConfig:     filepath.FromSlash("/etc/ssh/ssh_config"),
	SystemKnownHosts: filepath.FromSlash("/etc/ssh/known_hosts"),
	Identity: []string{
		filepath.Join(home, ".ssh", "id_dsa"),
		filepath.Join(home, ".ssh", "id_ecdsa"),
		filepath.Join(home, ".ssh", "id_ed25519"),
		filepath.Join(home, ".ssh", "id_rsa"),
	},
}

var home = currentUserHomeDir()

func currentUserHomeDir() string {
	u, err := user.Current()
	if err != nil {
		panic("unexpected error reading home dir: " + err.Error())
	}

	return u.HomeDir
}
