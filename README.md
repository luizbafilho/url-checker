# URL Checker

This project waits for webhooks requests from Github and checks the reachability of URLs in PR's descripions.

# Running

Download the project:
```shell
go get github.com/luizbafilho/url-checker
```

Generate a new [Github Personal Token](https://github.com/settings/tokens)

```shell
export GITHUB_ACCESS_TOKEN=<token>

cd $GOPATH/src/github.com/luizbafilho/url-checker
go build && ./url-checker
```