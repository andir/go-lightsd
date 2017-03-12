package main

import (
    "time"
    "flag"
    "./core"
    _ "./operations"
)

func main() {
    configPath := flag.String("config", "config.yml", "The config file")

    flag.Parse()

    config, err := LoadConfig(*configPath)
    if err != nil {
        panic(err)
    }

    pipeline := core.NewPipeline()

    for _, op := range config.Operations {
        pipeline.NewOperation(op.Type, op.Name, config.Size, op.Config)
    }

    NewMqttConnection(config, pipeline)

    bc := StartDebug()

    sink := NewSHMOutput("/test", config.Size)

    for {
        s := time.Now()

        for _, p := range pipeline {
            p.Lock()
            p.Render()
            p.Unlock()
        }

        elapsed := time.Now().Sub(s)

        result := pipeline[len(pipeline) - 1].Stripe()

        sink.Render(result)
        bc.Broadcast(result)
        interval := time.Second / time.Duration(config.FPS)

        diff := interval - elapsed
        time.Sleep(diff)
    }
}
