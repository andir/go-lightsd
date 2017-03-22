package debug

import (
    "log"
    "golang.org/x/net/websocket"
    "sync"
)

type broadcaster struct {
    sync.RWMutex
    clients []*websocket.Conn
}

func (b *broadcaster) Add(ws *websocket.Conn) {
    b.Lock()
    defer b.Unlock()

    b.clients = append(b.clients, ws)
}

func (b *broadcaster) Remove(ws *websocket.Conn) {
    b.Lock()
    defer b.Unlock()

    for i, c := range b.clients {
        if c == ws {
            b.clients = append(b.clients[:i], b.clients[i+1:]...)
            return
        }
    }
}

func (b *broadcaster) Broadcast(msg []byte) {
    b.Lock()
    defer b.Unlock()

    for _, c := range b.clients {
        if err := websocket.Message.Send(c, msg); err != nil {
            log.Println(err.Error())
        }
    }
}
