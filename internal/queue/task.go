package queue

import (
	"encoding/json"
	"log/slog"
	"path/filepath"

	"github.com/hibiken/asynq"
)

type TranscodePayload struct {
	VideoID   string `json:"video_id"`
	Filename  string `json:"filename"`
	InputPath string `json:"input_path"`
}

func EnqueueTranscode(client *asynq.Client, videoID, filePath string) error {
	slog.Info("file metadata", "videoID", videoID, "filePath", filePath)
	payload, err := json.Marshal(TranscodePayload{
		VideoID:   videoID,
		Filename:  filepath.Base(filePath),
		InputPath: filePath,
	})
	if err != nil {
		return err
	}

	info, err := client.Enqueue(asynq.NewTask(TypeVideoTranscode, payload))
	jsonInfo, _ := json.Marshal(info)
	slog.Info("task information", "info", jsonInfo)
	return err
}
