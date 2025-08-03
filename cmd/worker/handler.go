package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/hibiken/asynq"
	"github.com/theluminousartemis/video-transcoder/internal/queue"
)

var (
	uploadsDir     = "./uploads"
	outputsDirBase = "./outputs"
	scriptsDir     = "./scripts"
)

func (app *application) HandleTranscodeTask(ctx context.Context, t *asynq.Task) error {
	var payload queue.TranscodePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	// Run ffmpeg container
	containerInput := fmt.Sprintf("/uploads/%s", payload.Filename)
	cmd := exec.Command("docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/uploads:ro", uploadsDir),
		"-v", fmt.Sprintf("%s:/outputs", outputsDirBase),
		"-v", fmt.Sprintf("%s:/scripts:ro", scriptsDir),
		"--entrypoint", "bash",
		"jrottenberg/ffmpeg:latest",
		"/scripts/transcode.sh", containerInput,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	slog.Info("docker command", "args", cmd.Args)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker run failed: %w", err)
	}
	return nil
}
