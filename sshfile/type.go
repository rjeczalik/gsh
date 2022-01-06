package sshfile

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Duration time.Duration

var (
	_ json.Marshaler   = new(Duration)
	_ json.Unmarshaler = new(Duration)
)

func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.Itoa(int(time.Duration(d) / time.Second)))
}

func (d *Duration) UnmarshalJSON(p []byte) error {
	var s string
	if err := json.Unmarshal(p, &s); err != nil {
		return err
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*d = Duration(time.Duration(n) * time.Second)
	return nil
}

type Bool bool

var (
	_ json.Marshaler   = new(Bool)
	_ json.Unmarshaler = new(Bool)
)

func Boolean(b bool) *Bool {
	d := Bool(b)
	return &d
}

func (b *Bool) Bool() bool {
	return b != nil && bool(*b)
}

func (b Bool) MarshalJSON() ([]byte, error) {
	if b.Bool() {
		return []byte(`"yes"`), nil
	}
	return []byte(`"no"`), nil
}

func (b *Bool) UnmarshalJSON(p []byte) error {
	var s string
	if err := json.Unmarshal(p, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	case "yes":
		*b = true
	case "no":
		*b = false
	default:
		return fmt.Errorf("unexpected boolean value: %q", s)
	}
	return nil
}
