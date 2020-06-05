package pgbroadcaster

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// connection is a middleman between the websocket connection and the hub.
// it contains the actual websocket connection, a map where we keep the tables
// this connection is monitoring, and a send channel in which the  hub will
// write the outbound notifications.
type connection struct {
	ws *websocket.Conn
	send chan pgnotification
	h    *hub
}

// the reader listens for incoming messages on the websocket. the purpose is
// that the client sends a string message containing the tablename it wants to
// monitor.
func (c *connection) reader() {
	fmt.Println("entering reader")
	defer func() {
		fmt.Println("exiting reader")
		c.h.unregister <- c
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		t, m, err := c.ws.ReadMessage()
		if err != nil || t != websocket.TextMessage {
			break
		}
		//c.subscriptions[string(m)] = true
		fmt.Println("pgbroadcast: incoming [", string(m), "]")
	}
}

// writeMessage writes a message with the given message type and payload.
func (c *connection) writeMessage(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writeMessage marshals an interface to json and sends it.
func (c *connection) writeJSON(i interface{}) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteJSON(i)
}

// writer writes the pgnotificiations it receives from the hub to the websocket
// connection.
func (c *connection) writer() {
	fmt.Println("entering writer")
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		fmt.Println("exiting writer")
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case notification, ok := <-c.send:
			if !ok {
				c.writeMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.writeJSON(notification); err != nil {
				return
			}
		case <-ticker.C:
			fmt.Println("ticker ..")
			if err := c.writeMessage(websocket.PingMessage, []byte("tick")); err != nil {
				return
			}
		}
	}
}

// serverWs handles websocket requests from the peer.
func (pb *PgBroadcaster) WebsocketHandler(c echo.Context) error {

	fmt.Println("entering WebsocketHandler")

	// upgrade the connection to a websocket connection
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		return nil
	}

	// create a new connection struct and add it to the hub.
	cc := &connection{
		ws: ws,
		send: make(chan pgnotification, 1024),
		h:    pb.h,
	}
	pb.h.register <- cc

	// start the connections reader and writer.
	go cc.writer()
	cc.reader()

	return nil
}

// for testing front-end
func WebsocketHandlerOld(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			c.Logger().Error(err)
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
}
