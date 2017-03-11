package main

import (
    "image/color"
    "log"
    "time"
    "net/http"
    "golang.org/x/net/websocket"
    "sync"
    "fmt"
    "flag"
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

func main() {
    broker := flag.String("broker", "tcp://whisky.w17.io:1883", "The broker URI. ex: tcp://whisky.w17.io:1883")
    id := flag.String("id", "super-lightsd", "The ClientID (optional)")

    stripe := NewLEDStripe(1000)

    fps := 60

    pipeline := map[string]Operation{
        //"rainbow":   NewRainbow(),
        "raindrops": NewRaindrop(),
        //"rotation":  NewRotation(60.0),
    }

    NewMqttConnection(*broker, *id, pipeline)

    bc := WebsocketBroadcaster{}

    go func() {
        http.Handle("/stream", websocket.Handler(CreateStreamHandler(&bc)))
        http.Handle("/", http.FileServer(http.Dir("web")))
        err := http.ListenAndServe(":9000", nil)
        if err != nil {
            panic("ListenAndServe: " + err.Error())
        }
    }()

    sink := NewSHMOutput("/test", len(stripe))

    for {
        s := time.Now()
        for i := range pipeline {
            //log.Printf("%v", i)
            pipeline[i].Render(stripe)
        }

        elapsed := time.Now().Sub(s)

        sink.Render(stripe)
        bc.Broadcast(stripe)
        interval := time.Second / time.Duration(fps)

        diff := interval - elapsed
        time.Sleep(diff)
    }
}
