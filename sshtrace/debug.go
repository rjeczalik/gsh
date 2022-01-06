package sshtrace

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/tectumsh/gsh"
	"github.com/tectumsh/gsh/sshfile"
)

var _ = Debug("/tmp")

func Debug(tmpdir string) *ClientTrace {
	if err := os.MkdirAll(tmpdir, 0755); err != nil {
		panic("unexpected error: " + err.Error())
	}

	dir, err := ioutil.TempDir(tmpdir, "glaucos-ssh")
	if err != nil {
		panic("unexpected error: " + err.Error())
	}

	return &ClientTrace{
		GotFileConfig: func(cfg sshfile.Configs) {
			dump("GotFileConfig", "sshfile-config-*.json", dir, cfg)
		},
		GotConfig: func(cfg *gsh.Config) {
			dump("GotConfig", "ssh-config-*.json", dir, cfg)
		},
	}
}

func dump(prefix, pattern, dir string, v interface{}) {
	f, err := ioutil.TempFile(dir, pattern)
	if err != nil {
		log.Printf("failed to create temporary file: %s", err)
		return
	}

	log.Printf("%s: dumping %T to file: %s", prefix, v, f.Name())

	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")

	if err := nonil(enc.Encode(v), f.Close()); err != nil {
		log.Printf("failed to dump %T to a file: %s", v, err)
	}
}
