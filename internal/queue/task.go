package queue

import (
	"encoding/json"
	"log/slog"
	"path/filepath"

	"github.com/hibiken/asynq"
)

var (
	ResolutionsSlice = []string{"360p", "480p", "720p"}
)

type TranscodePayload struct {
	VideoID     string   `json:"video_id"`
	Filename    string   `json:"filename"`
	InputPath   string   `json:"input_path"`
	OutputDir   string   `json:"output_dir"`
	Resolutions []string `json:"resolutions"`
}

func EnqueueTranscode(client *asynq.Client, videoID, filePath string) error {
	slog.Info("file metadata", "videoID", videoID, "filePath", filePath)
	payload, err := json.Marshal(TranscodePayload{
		VideoID:     videoID,
		Filename:    filepath.Base(filePath),
		InputPath:   filePath,
		OutputDir:   filepath.Join("outputs", videoID),
		Resolutions: ResolutionsSlice,
	})
	if err != nil {
		return err
	}

	info, err := client.Enqueue(asynq.NewTask(TypeVideoTranscode, payload))
	slog.Info("task information", "task id", info.ID)
	return err
}
