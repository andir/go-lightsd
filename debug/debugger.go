package debug

import (
    "github.com/andir/lightsd/core"
    "net/http"
    "golang.org/x/net/websocket"
    "fmt"
    "html/template"
)

type Debugger struct {
    pipelines []*core.Pipeline

    broadcasters map[string]*broadcaster
}

func (this *Debugger) createIndexHandler() http.Handler {
    template := template.Must(template.ParseFiles("debug/templates/index.html"))

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        err := template.Execute(w, this.pipelines)
        if err != nil {
            panic(fmt.Errorf("Failed to render template: index.html: %v", err))
        }
    })
}

func (this *Debugger) createPipelineHandler() http.Handler {
    template := template.Must(template.ParseFiles("debug/templates/pipeline.html"))

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        name := r.URL.Path

        var pipeline *core.Pipeline
        for _, p := range this.pipelines {
            if p.Name == name {
                pipeline = p
                break
            }
        }

        err := template.Execute(w, pipeline)
        if err != nil {
            panic(fmt.Errorf("Failed to render template: pipeline.html: %v", err))
        }
    })
}

func (this *Debugger) createStreamHandler() http.Handler {
    return websocket.Handler(func(ws *websocket.Conn) {
        name := ws.Request().URL.Path

        broadcaster := this.broadcasters[name]

        broadcaster.Add(ws)
        defer broadcaster.Remove(ws)

        for {
            var msg string
            if err := websocket.Message.Receive(ws, &msg); err != nil {
                fmt.Println("Recv error: ", err.Error())
                break
            }
        }
    })
}

func StartDebug(port int, pipelines []*core.Pipeline) *Debugger {
    broadcasters := make(map[string]*broadcaster, len(pipelines))
    for _, pipeline := range pipelines {
        broadcasters[pipeline.Name] = &broadcaster{}
    }

    debugger := &Debugger{
        pipelines:    pipelines,
        broadcasters: broadcasters,
    }

    mux := http.NewServeMux()
    mux.Handle("/pipeline/", http.StripPrefix("/pipeline/", debugger.createPipelineHandler()))
    mux.Handle("/stream/", http.StripPrefix("/stream/", debugger.createStreamHandler()))
    mux.Handle("/", debugger.createIndexHandler())

    go func() {
        err := http.ListenAndServe(":9000", mux)
        if err != nil {
            panic(fmt.Errorf("Failed to server debug interface: %v", err))
        }
    }()

    return debugger
}

func (this *Debugger) Broadcast(context *core.RenderContext) {
    msg := make([]byte, context.Count()*3)

    for i := 0; i < context.Count(); i++ {
        msg[i*3+0], msg[i*3+1], msg[i*3+2] = context.Results[context.Pipeline.Output.Source()].Get(i)
    }

    this.broadcasters[context.Pipeline.Name].Broadcast(msg)
}
