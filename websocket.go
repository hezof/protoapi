package protoapi

import (
	"github.com/hezof/protoapi/internal/websocket"
	"net/http"
)

type Upgrader = websocket.Upgrader

// newWebsocketUpgrader 根据配置生成websocket的upgrader
func newWebsocketUpgrader(c *Config) *Upgrader {
	upgrader := new(Upgrader)
	upgrader.ReadBufferSize = c.WbskReadBuffer
	upgrader.WriteBufferSize = c.WbskWriteBuffer
	if c.WbskOriginDisable {
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
	}
	return upgrader
}
