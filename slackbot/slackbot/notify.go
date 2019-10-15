package slackbot

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	cloudbuild "google.golang.org/api/cloudbuild/v1"
)

// Notify posts a notification to Slack that the build is complete.
func Notify(b *cloudbuild.Build, webhook string, project string) {
	url := fmt.Sprintf("https://console.cloud.google.com/cloud-build/builds/%s", b.Id)
	var i string
	switch b.Status {
	case "SUCCESS":
		i = ":white_check_mark:"
	case "FAILURE", "CANCELLED":
		i = ":x:"
	case "STATUS_UNKNOWN", "INTERNAL_ERROR":
		i = ":interrobang:"
	default:
		i = ":question:"
	}
	j := fmt.Sprintf(
		`{"text": "%s : %s - Build (%s) complete: %s %s",
		    "attachments": [
				{
					"fallback": "Open build details at %s",
					"actions": [
						{
							"type": "button",
							"text": "Open details",
							"url": "%s"
						}
					]
				}
			]}`, project, b.Source.RepoSource, b.Id[0:7], i, b.Status, url, url)

	r := strings.NewReader(j)
	resp, err := http.Post(webhook, "application/json", r)
	if err != nil {
		log.Fatalf("Failed to post to Slack: %v", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("Posted message to Slack: [%v], got response [%s]", j, body)
}
