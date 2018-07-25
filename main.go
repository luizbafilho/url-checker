package main

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

// Check checks for reachable urls in the Pull Request body
func (uc *UrlChecker) Check(pr *PullRequest) ([]URLStatus, error) {
	return nil, nil
}

func urlReachable(url string) bool {
	return false
}

func main() {
}
