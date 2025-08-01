package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func (app *application) uploadVideo(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(100 * 1024 * 1024)
	if err != nil {
		app.logger.Info("file exceeds 100MB")
		app.badRequestError(w, r, err)
	}

	video, fileHeader, err := r.FormFile("video")
	if err != nil {
		app.logger.Info("Unable to parse file")
		app.internalServerError(w, r, err)
	}

	defer video.Close()

	id := uuid.New().String()
	ext := filepath.Ext(fileHeader.Filename)
	filename := fmt.Sprintf("%s%s", id, ext)
	savePath := filepath.Join("uploads", filename)

	dst, err := os.Create(savePath)
	if err != nil {
		app.logger.Info("Unable to create file")
		app.internalServerError(w, r, err)
	}

	defer dst.Close()

	_, err = io.Copy(dst, video)
	if err != nil {
		app.logger.Info("Unable to copy file")
		app.internalServerError(w, r, err)
	}

	WriteJSON(w, http.StatusOK, "file uploaded successfully")
}
