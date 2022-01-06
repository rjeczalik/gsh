package sshfile_test

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"testing"
)

var updateGolden = flag.Bool("update-golden", false, "Update golden files")

func UnmarshalFile(file string, v interface{}) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	return dec.Decode(v)
}

func MarshalFile(v interface{}, file string) error {
	p, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, p, 0644)
}

func TestMain(m *testing.M) {
	flag.Parse()

	os.Exit(m.Run())
}
