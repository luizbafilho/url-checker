package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
)

func main() {
	githubAccessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	if githubAccessToken == "" {
		log.Fatal("Invalid GITHUB_ACCESS_TOKEN environment variable")
	}
	ghs := NewGHService(githubAccessToken)

	uc := NewURLChecker(2 * time.Second)

	e := echo.New()
	e.POST("/", func(c echo.Context) error {
		payload := WebhookPayload{}
		if err := c.Bind(&payload); err != nil {
			return err
		}

		if payload.PullRequest == nil {
			return c.NoContent(http.StatusNoContent)
		}

		statuses := uc.Check(payload.PullRequest)

		if err := ghs.CreatePullRequestComment(payload.PullRequest, payload.Repository, statuses); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("PR's comment creation failed: %s", err))
		}

		return c.NoContent(http.StatusNoContent)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
