package sshfile_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/rjeczalik/gsh/sshfile"
	"github.com/google/go-cmp/cmp"
)

func TestConfig(t *testing.T) {
	want := &sshfile.Config{
		Port:                  22,
		StrictHostKeyChecking: sshfile.Boolean(true),
		GlobalKnownHostsFile:  "/dev/null",
		UserKnownHostsFile:    "/dev/null",
		TcpKeepAlive:          sshfile.Boolean(true),
		ConnectTimeout:        sshfile.Duration(10 * time.Second),
		ConnectionAttempts:    3,
		ServerAliveInterval:   sshfile.Duration(5 * time.Second),
		ServerAliveCountMax:   10,
	}

	p, err := json.MarshalIndent(want, "", "\t")
	if err != nil {
		t.Fatalf("json.Marshal()=%s", err)
	}

	got := new(sshfile.Config)

	if err := json.Unmarshal(p, got); err != nil {
		t.Fatalf("json.Unmarshal()=%s", err)
	}

	if !cmp.Equal(got, want) {
		t.Fatalf("got != want:\n%s\n", cmp.Diff(got, want))
	}
}

func TestParseConfig(t *testing.T) {
	got, err := sshfile.ParseConfigFile("testdata/config")
	if err != nil {
		t.Fatalf("ParseConfig()=%s", err)
	}

	if *updateGolden {
		if err := MarshalFile(got, "testdata/config.golden"); err != nil {
			t.Fatalf("MarshalFile()=%s", err)
		}

		return
	}

	var want sshfile.Configs

	if err := UnmarshalFile("testdata/config.golden", &want); err != nil {
		t.Fatalf("UnmarshalFile()=%s", err)
	}

	if got := got; !cmp.Equal(got, want) {
		t.Fatalf("got != want:\n%s\n", cmp.Diff(got, want))
	}
}

func TestParseArgs(t *testing.T) {
	var tests [][]string

	if err := UnmarshalFile("testdata/flags.json", &tests); err != nil {
		t.Fatalf("UnmarshalFile()=%s", err)
	}

	var wants []*sshfile.Config

	if !*updateGolden {
		if err := UnmarshalFile("testdata/flags.golden.json", &wants); err != nil {
			t.Fatalf("UnmarshalFile()=%s", err)
		}
	}

	for i, flags := range tests {
		t.Run("", func(t *testing.T) {
			got, err := sshfile.ParseArgs(flags)
			if err != nil {
				t.Fatalf("ParseArgs()=%s", err)
			}

			if *updateGolden {
				wants = append(wants, got)
				return
			}

			if want := wants[i]; !cmp.Equal(got, want) {
				t.Fatalf("got != want:\n%s\n", cmp.Diff(got, want))
			}
		})
	}

	if *updateGolden {
		if err := MarshalFile(wants, "testdata/flags.golden.json"); err != nil {
			t.Fatalf("MarshalFile()=%s", err)
		}
	}
}
