package server

import (
	"net/http"
	"sync"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/twinj/uuid"
)

type WebSocketServer struct {
	mu      sync.Mutex
	clients map[string]*Client
}

type Client struct {
	id     string
	socket *websocket.Conn
}

func NewWebSocketServer() IServer {
	return &WebSocketServer{
		clients: make(map[string]*Client, 128),
	}
}

func (s *WebSocketServer) Run() {

}

func (s *WebSocketServer) wsHandle(c *gin.Context) {
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{c.GetHeader("Sec-WebSocket-Protocol")},
	}

	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logs.Error("websocket conect error: %v", err)
	}

	client := &Client{
		id:     uuid.NewV4().String(),
		socket: conn,
	}

}

func (s *WebSocketServer) registerClient(client *Client) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clients[client.id] = client
	return nil
}

func (s *WebSocketServer) unRegisterClient(client *Client) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.clients, client.id)
	return nil
}
