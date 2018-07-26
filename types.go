package main

type WebhookPayload struct {
	Repository  *Repo        `json:"repository"`
	PullRequest *PullRequest `json:"pull_request"`
}

type Repo struct {
	Name  string
	Owner Owner
}

type Owner struct {
	Login string
}

type PullRequest struct {
	Number int
	Body   string
}

// URLStatus represents the status of a URL if it is reachable or not
type URLStatus struct {
	URL       string
	Reachable bool
}
