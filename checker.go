package main

import (
	"log"
	"net"
	"net/url"
	"regexp"
	"time"
)

// very naive regex I got on google, it won't cover all cases, but the basic ones should
var reURL = regexp.MustCompile(`(?im)https?:\/\/[A-Z\d\.-]{2,}\.[A-Z]{2,}(:\d{2,4})?`)

type Checker interface {
	Check(pr *PullRequest) ([]URLStatus, error)
}

type URLChecker struct {
	Timeout time.Duration
}

// NewURLChecker creates a new instace of URLChecker
func NewURLChecker(timeout time.Duration) *URLChecker {
	return &URLChecker{
		Timeout: timeout,
	}
}

// Check parses the pr's body looking for urls and checks their reachability
func (uc *URLChecker) Check(pr *PullRequest) []URLStatus {
	// extracting all urls from pr's body
	matches := reURL.FindAllString(pr.Body, -1)
	matchesLen := len(matches)

	statusesCh := make(chan URLStatus, matchesLen)
	for _, m := range matches {
		go func(url string) {
			status, err := uc.urlReachable(url)
			if err != nil {
				log.Printf("url=%s is not reachable: %s", url, err)
			}
			statusesCh <- URLStatus{url, status}
		}(m)
	}

	statuses := []URLStatus{}
	for i := 0; i < matchesLen; i++ {
		statuses = append(statuses, <-statusesCh)
	}

	return statuses
}

func (uc *URLChecker) urlReachable(u string) (bool, error) {
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

	conn, err := net.DialTimeout("tcp", address, uc.Timeout)
	if err != nil {
		return false, err
	}
	conn.Close()

	return true, nil
}
