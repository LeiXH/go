package controller

import (
	"context"
	"encoding/json"
	"github.com/labstack/echo"
	"net/http"
	"pkg/server/controller/response"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/atomic"
	"pkg/logger"
	"pkg/models"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 1 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 1 * time.Second
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var webSocketConned = atomic.NewBool(false)

var userInfoChan = make(chan models.UserInfo)

func PushFaceDetectResultToFront(c echo.Context) error {

	if webSocketConned.CAS(true, true) {
		return c.JSON(http.StatusLocked, "ws has been established")
	}
	webSocketConned.Store(true)

	ctx, cancle := context.WithCancel(context.Background())
	defer cancle()

	// 升级 get 请求为 webSocket 协议
	ws, err := upGrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logger.Errorf("ws upgrade failed: %s", err)
		return c.JSON(http.StatusInternalServerError, response.RenderError(500, err))
	}
	defer ws.Close()

	// handle cliet close the connection
	ws.SetCloseHandler(func(code int, text string) error {
		logger.Warn("ws connection closed by client")
		cancle()
		webSocketConned.Toggle()
		return nil
	})

	tick := time.NewTicker(3 * time.Second)
	// pop and clean client msg to receive control-message
	go func() {
	ListenWebSock:
		for {
			select {
			case <-tick.C:
				{
					_, _, _ = ws.ReadMessage()
				}
			case <-ctx.Done():
				break ListenWebSock
			}
		}
	}()

WSEND:
	for {
		select {
		// get user info
		case info := <-userInfoChan:
			message, _ := json.Marshal(info)
			// 写入 ws 数据
			_ = ws.SetWriteDeadline(time.Now().Add(writeWait))
			err = ws.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				userInfoChan <- info
				logger.Errorf("write user-info message to front error, %s", err)
				break WSEND
			}
		case <-ctx.Done():
			break WSEND
		}
	}

	logger.Debug("websocket connection broken")
	return nil
}
