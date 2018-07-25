package main

import (
	"net"
	"net/url"
	"regexp"
	"time"
)

// PullRequest represents a pull request struct
type PullRequest struct {
	Body string
}

// URLStatus represents the status of a URL if it is reachable or not
type URLStatus struct {
	URL       string
	Reachable bool
}

type Checker interface {
	Check(pr *PullRequest) ([]URLStatus, error)
}

type UrlChecker struct {
}

// very naive regex I got on google, it won't cover all cases, but the basic ones should
var reURL = regexp.MustCompile(`(?im)https?:\/\/[A-Z\d\.-]{2,}\.[A-Z]{2,}(:\d{2,4})?`)

// Check checks for reachable urls in the Pull Request body
func (uc *UrlChecker) Check(pr *PullRequest) ([]URLStatus, error) {
	statuses := []URLStatus{}

	return statuses, nil
}

func urlReachable(u string) (bool, error) {
	timeout := 10 * time.Second
	parsedURL, err := url.Parse(u)
	if err != nil {
		return false, err
	}

	address := parsedURL.Host
	if parsedURL.Port() == "" {
		switch parsedURL.Scheme {
		case "http":
			address += ":80"
		case "https":
			address += ":443"
		}
	}

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false, err
	}
	conn.Close()

	return true, nil
}

func main() {

}
