package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

func (app *application) wsHandler(w http.ResponseWriter, r *http.Request) {
	videoID := r.URL.Query().Get("id")
	if videoID == "" {
		app.logger.Error("missing video id")
		app.badRequestError(w, r, errors.New("missing video id for websockets"))
		return
	}
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		app.logger.Error("upgrade to websocket failed", "err", err)
		app.internalServerError(w, r, err)
		return
	}

	defer conn.Close()

	ctx := context.Background()
	sub := app.rdb.Subscribe(ctx, fmt.Sprintf("video:%s", videoID))
	defer sub.Close()

	ch := sub.Channel()

	for msg := range ch {
		var payload map[string]interface{}
		if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
			app.logger.Error("invalid message from redis", "err", err)
			continue
		}

		// Send to websocket client
		if err := conn.WriteJSON(payload); err != nil {
			app.logger.Error("write to websocket failed", "err", err)
			return
		}
	}
}
