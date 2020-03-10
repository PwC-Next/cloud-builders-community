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
func Notify(b *cloudbuild.Build, t *cloudbuild.BuildTrigger, webhook string, project string) {
	url := fmt.Sprintf("https://console.cloud.google.com/cloud-build/builds/%s?project=%s", b.Id, project)
	var i string
	var branch_string string
	branch_string = "Trigger details not found - see logs"

	if t.Name != "" {
		branch_string = fmt.Sprintf("Trigger: %s - %s.", t.Name, t.Description )
	}

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
		`{"text": "%s - Build (%s) complete: %s %s \n %s ",
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
			]}`, project, b.Id[0:7], i, b.Status, branch_string, url, url)

	r := strings.NewReader(j)
	resp, err := http.Post(webhook, "application/json", r)
	if err != nil {
		log.Fatalf("Failed to post to Slack: %v", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("Posted message to Slack: [%v], got response [%s]", j, body)
}
