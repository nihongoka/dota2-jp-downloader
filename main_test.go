package main

import (
	"testing"
)

func TestEscape(t *testing.T) {
	sd := map[string]string{
		"abcdefg":      "abcdefg",
		"あいうえお":        "あいうえお",
		"s\\q\"t\tn\n": "s\\\\q\\\"t\\tn\\n",
	}

	for s, d := range sd {
		if d != escape(s) {
			t.Error(s)
		}
	}
}
