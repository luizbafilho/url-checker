package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlReachable(t *testing.T) {
	var cases = []struct {
		url     string
		status  bool
		noError bool
	}{
		{"https://google.com", true, true},
		{"http://google.com", true, true},
		{"http://foo.local", false, false},
		{"https://yahoo.com:443", true, true},
	}

	for _, c := range cases {
		t.Run(c.url, func(t *testing.T) {
			t.Parallel()
			reachable, err := urlReachable(c.url)
			assert.Equal(t, c.noError, err == nil)
			assert.Equal(t, c.status, reachable)
		})
	}
}
