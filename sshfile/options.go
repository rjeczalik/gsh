package sshfile

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

func ParseArgs(args []string) (*Config, error) {
	var options []string

	f := pflag.NewFlagSet("ssh", pflag.ContinueOnError)
	f.StringArrayVarP(&options, "option", "o", nil, "")

	if err := f.Parse(args); err != nil {
		return nil, fmt.Errorf("unable to parse flags: %w", err)
	}

	return ParseOptions(options)
}

func ParseOptions(options []string) (*Config, error) {
	tmp := make(map[string]string)

	for _, kv := range options {
		k, v, err := parsekv(kv)
		if err != nil {
			return nil, fmt.Errorf("unexpected %q flag: %w", kv, err)
		}

		tmp[strings.ToLower(k)] = v
	}

	hc := new(Config)

	if err := merge(hc, tmp); err != nil {
		return nil, fmt.Errorf("unexpected flags: %w", err)
	}

	return hc, nil
}
