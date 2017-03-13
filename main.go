package main

import (
    "flag"
    "time"
    
    "github.com/andir/lightsd/core"
    _ "github.com/andir/lightsd/operations"
    "github.com/andir/lightsd/outputs/shm"
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
    sink := shm.NewSHMOutput("/test", config.Size)

    interval := time.Second / time.Duration(config.FPS)

    for {
        elapsed := pipeline.Render()

        result := pipeline.Result()

        sink.Render(result)
        bc.Broadcast(result)

        diff := interval - elapsed
        time.Sleep(diff)
    }
}
