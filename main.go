package main

import (
    "flag"
    "time"
    
    "github.com/andir/lightsd/core"
    _ "github.com/andir/lightsd/operations"
    "github.com/andir/lightsd/outputs/shm"
    "runtime/pprof"
    "os"
    "fmt"
    "os/signal"
)

func main() {
    configPath := flag.String("config", "config.yml", "The config file")
    profileOutput := flag.String("cpuprofile", "", "Output file for profile output")

    flag.Parse()

    if *profileOutput != "" {
        fh, err := os.Create(*profileOutput)

        if err != nil {
            panic(err)
        }
        defer fh.Close()

        pprof.StartCPUProfile(fh)
        defer pprof.StopCPUProfile()
    }

    config, err := LoadConfig(*configPath)
    if err != nil {
        panic(err)
    }

    pipeline := core.NewPipeline()

    for _, op := range config.Operations {
        pipeline.NewOperation(op.Type, op.Name, config.Size, op.Config)
    }

    //NewMqttConnection(config, pipeline)

    bc := StartDebug()
    sink := shm.NewSHMOutput("/test", config.Size)

    interval := time.Second / time.Duration(config.FPS)

    go func() {
        for {
            elapsed := pipeline.Render()

            result := pipeline.Result()

            sink.Render(result)
            bc.Broadcast(result)

            diff := interval - elapsed
            time.Sleep(diff)
        }
    }()

    signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, os.Interrupt)
    for _ = range signalChan {
        fmt.Println("Interrupt received. Stopping")
        return
    }
}
