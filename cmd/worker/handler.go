package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/theluminousartemis/video-transcoder/internal/queue"
)

var (
	uploadsDir = "./uploads"
	outputsDir = "./outputs"
	scriptsDir = "./scripts"
)

//single container approach

// func (app *application) HandleTranscodeTask(ctx context.Context, t *asynq.Task) error {
// 	var payload queue.TranscodePayload
// 	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
// 		return err
// 	}

// 	// Run ffmpeg container
// 	containerInput := fmt.Sprintf("/uploads/%s", payload.Filename)
// 	cmd := exec.Command("docker", "run", "--rm",
// 		"-v", fmt.Sprintf("%s:/uploads:ro", uploadsDir),
// 		"-v", fmt.Sprintf("%s:/outputs", outputsDir),
// 		"-v", fmt.Sprintf("%s:/scripts:ro", scriptsDir),
// 		"--entrypoint", "bash",
// 		"jrottenberg/ffmpeg:latest",
// 		"/scripts/transcode.sh", containerInput,
// 	)

// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	slog.Info("docker command", "args", cmd.Args)

// 	if err := cmd.Run(); err != nil {
// 		return fmt.Errorf("docker run failed: %w", err)
// 	}
// 	return nil
// }

// cli orchestrator
func (app *application) HandleTranscodeTask(ctx context.Context, t *asynq.Task) error {
	var payload queue.TranscodePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	serviceName := fmt.Sprintf("transcode-%s", payload.VideoID)

	cmd := exec.Command(
		"docker", "service", "create",
		"--name", serviceName,
		"--restart-condition", "none",
		"--mount", fmt.Sprintf("type=bind,src=%s,dst=/uploads,ro", uploadsDir),
		"--mount", fmt.Sprintf("type=bind,src=%s,dst=/outputs", outputsDir),
		"--mount", fmt.Sprintf("type=bind,src=%s,dst=/scripts,ro", scriptsDir),
		"--entrypoint", "bash",
		"jrottenberg/ffmpeg:latest",
		"/scripts/transcode.sh", fmt.Sprintf("/uploads/%s", payload.Filename),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	slog.Info("docker command", "args", cmd.Args)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	for {
		out, _ := exec.Command(
			"docker", "service", "ps", serviceName,
			"--no-trunc", "--format", "{{.CurrentState}}",
		).Output()

		state := strings.TrimSpace(string(out))
		slog.Info("service state", "service", serviceName, "state", state)

		if strings.HasPrefix(state, "Complete") || strings.HasPrefix(state, "Failed") {
			break
		}
		time.Sleep(5 * time.Second)
	}

	exec.Command("docker", "service", "logs", serviceName).Run()

	exec.Command("docker", "service", "rm", serviceName).Run()

	return nil
}
