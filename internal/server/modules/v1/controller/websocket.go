package controller

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	_logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/consts"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
)

const (
	result = "result"
	outPut = "output"
)

var (
	ch chan int
)

type WsCtrl struct {
	Namespace         string
	*websocket.NSConn `stateless:"true"`

	WebSocketService *service.WebSocketService `inject:""`
}

func NewWsCtrl() *WsCtrl {
	inst := &WsCtrl{Namespace: serverConsts.WsDefaultNameSpace}
	return inst
}

func (c *WsCtrl) OnNamespaceConnected(msg websocket.Message) error {
	c.WebSocketService.SetConn(c.Conn)

	_logUtils.Infof("WebSocket OnNamespaceConnected: ConnID=%s, Room=%s", c.Conn.ID(), msg.Room)

	data := map[string]string{"msg": "from server: connected to websocket"}
	c.WebSocketService.Broadcast(msg.Namespace, msg.Room, "OnVisit", data)
	return nil
}

// OnNamespaceDisconnect
// This will call the "OnVisit" event on all clients, except the current one,
// it can't because it's left but for any case use this type of design
func (c *WsCtrl) OnNamespaceDisconnect(msg websocket.Message) error {
	_logUtils.Infof("WebSocket OnNamespaceDisconnect: ConnID=%s", c.Conn.ID())

	data := map[string]string{"msg": "from server: disconnected to websocket"}
	c.WebSocketService.Broadcast(msg.Namespace, msg.Room, "OnVisit", data)
	return nil
}

// OnChat This will call the "OnVisit" event on all clients, including the current one, with the 'newCount' variable.
func (c *WsCtrl) OnChat(msg websocket.Message) (err error) {
	ctx := websocket.GetContext(c.Conn)

	_logUtils.Infof("WebSocket OnChat: remote address=%s, room=%s, msg=%s", ctx.RemoteAddr(), msg.Room, string(msg.Body))

	return
}

func (c *WsCtrl) TestWs(ctx iris.Context) {
	data := map[string]interface{}{"action": "taskUpdate", "taskId": 1, "msg": ""}
	c.WebSocketService.Broadcast(serverConsts.WsDefaultNameSpace, serverConsts.WsDefaultRoom, serverConsts.WsChatEvent, data)

	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: nil, Msg: _domain.NoErr.Msg})
}
