package debug

//go:generate go-bindata -pkg $GOPACKAGE -o templates.go templates/

import (
    "github.com/andir/lightsd/core"
    "net/http"
    "golang.org/x/net/websocket"
    "html/template"
    "log"
    "io"
)

type Debugger struct {
    pipelines []*core.Pipeline

    broadcasters map[string]*broadcaster
}

func (this *Debugger) createIndexHandler() http.Handler {
    template := template.Must(template.New("index.html").Parse(string(MustAsset("templates/index.html"))))

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        err := template.Execute(w, this.pipelines)
        if err != nil {
            log.Panicf("debugger: failed to render template: index.html: %v", err)
        }
    })
}

func (this *Debugger) createPipelineHandler() http.Handler {
    template := template.Must(template.New("pipeline.html").Parse(string(MustAsset("templates/pipeline.html"))))

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
            log.Panic("debugger: failed to render template: pipeline.html:", err)
        }
    })
}

func (this *Debugger) createStreamHandler() http.Handler {
    return websocket.Handler(func(ws *websocket.Conn) {
        name := ws.Request().URL.Path

        broadcaster := this.broadcasters[name]

        broadcaster.Add(ws)
        defer broadcaster.Remove(ws)

        var msg string
        for {
            err := websocket.Message.Receive(ws, &msg)
            if err != nil {
                if err != io.EOF {
                    log.Print("debuger: recv error:", err.Error())
                }

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
            log.Panic("debugger: failed to listen:", err)
        }
    }()

    return debugger
}

func (this *Debugger) Broadcast(context *core.RenderContext) {
    msg := make([]byte, context.Count()*3)

    for i := 0; i < context.Count(); i++ {
        c := context.Results[context.Pipeline.Output.Source()].Get(i)
        msg[i*3+0], msg[i*3+1], msg[i*3+2] = c.RGB255()
    }

    this.broadcasters[context.Pipeline.Name].Broadcast(msg)
}
