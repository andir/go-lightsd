package main

import (
    "image/color"
    "log"
    "net/http"
    "golang.org/x/net/websocket"
    "sync"
    "fmt"
    _ "./operations"
)

func CreateStreamHandler(broadcaster *WebsocketBroadcaster) websocket.Handler {
    return func(ws *websocket.Conn) {
        broadcaster.Add(ws)

        for {
            var msg string
            if err := websocket.Message.Receive(ws, &msg); err != nil {
                fmt.Println("Recv error: ", err.Error())
                break
            }
        }
        broadcaster.Remove(ws)
    }
}

type WebsocketBroadcaster struct {
    sync.RWMutex
    clients []*websocket.Conn
}

func (b *WebsocketBroadcaster) Add(ws *websocket.Conn) {
    b.Lock()
    defer b.Unlock()

    b.clients = append(b.clients, ws)
}

func (b *WebsocketBroadcaster) Remove(ws *websocket.Conn) {
    b.Lock()
    defer b.Unlock()

    for i, c := range b.clients {
        if c == ws {
            b.clients = append(b.clients[:i], b.clients[i+1:]...)
            return
        }
    }
}

func (b *WebsocketBroadcaster) Broadcast(l []color.RGBA) {
    msg := make([]byte, len(l)*3)

    for i, p := range l {
        msg[i*3+0] = byte(p.R)
        msg[i*3+1] = byte(p.G)
        msg[i*3+2] = byte(p.B)
    }

    for _, c := range b.clients {
        if err := websocket.Message.Send(c, msg); err != nil {
            log.Println(err.Error())
        }
    }
}

func StartDebug() *WebsocketBroadcaster {
    bc := &WebsocketBroadcaster{}

    go func() {
        http.Handle("/stream", websocket.Handler(CreateStreamHandler(bc)))
        http.Handle("/", http.FileServer(http.Dir("web")))
        err := http.ListenAndServe(":9000", nil)
        if err != nil {
            panic("ListenAndServe: " + err.Error())
        }
    }()

    return bc
}
