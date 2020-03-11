package slackbot

import (
	"context"
	"fmt"
	"log"

	cloudbuild "google.golang.org/api/cloudbuild/v1"
)

// Trigger starts an independent watcher build.
func Trigger(ctx context.Context, build string, webhook string) {
	svc := gcbClient(ctx)
	b := &cloudbuild.Build{
		Steps: []*cloudbuild.BuildStep{
			&cloudbuild.BuildStep{
				Name: "gcr.io/acuit-dev/slackbot",
				Args: []string{
					fmt.Sprintf("--build=%s", build),
					fmt.Sprintf("--webhook=%s", webhook),
					"--mode=monitor",
				},
				Timeout: "3600",
			},
		},
		Tags: []string{"slackbot"},
		Timeout: "3610",
	}

	project, err := getProject()
	if err != nil {
		log.Fatalf("Failed to get project: %v", err)
	}

	cr := svc.Projects.Builds.Create(project, b)
	_, err = cr.Do()
	if err != nil {
		log.Fatalf("Failed to create watcher build: %v", err)
	} else {
		log.Printf("Triggered watcher build.")
	}
}
