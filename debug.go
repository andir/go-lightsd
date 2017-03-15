package main

import (
    "github.com/andir/lightsd/core"
    "log"
    "net/http"
    "golang.org/x/net/websocket"
    "sync"
    "fmt"
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

func (b *WebsocketBroadcaster) Broadcast(pipeline *core.Pipeline, context *core.RenderContext) {
    msg := make([]byte, pipeline.Count()*3)

    stripe := context.Results[pipeline.Output().Source()]

    for i := 0; i < stripe.Count(); i++ {
        msg[i*3+0], msg[i*3+1], msg[i*3+2] = stripe.Get(i)
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
