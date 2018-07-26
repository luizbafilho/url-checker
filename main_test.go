package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	pr := &PullRequest{
		Body: `Lorem ipsum dolor sit ame http://www.google.com
		llo inventore veritatis https://yahoo.com et quasi
		architecto beatae vitae dicta sunt explicabo.  consequuntur magni
		dolores eos qui http://foo.local:8888
		`,
	}

	uc := URLChecker{
		Timeout: 1 * time.Second,
	}
	assert.Equal(t, []URLStatus{
		{"http://www.google.com", true},
		{"https://yahoo.com", true},
		{"http://foo.local:8888", false},
	}, uc.Check(pr))
}

func TestUrlReachable(t *testing.T) {
	uc := URLChecker{
		Timeout: 1 * time.Second,
	}
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
			reachable, err := uc.urlReachable(c.url)
			assert.Equal(t, c.noError, err == nil)
			assert.Equal(t, c.status, reachable)
		})
	}
}
