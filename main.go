package main

import (
    "flag"
    "time"
    "runtime/pprof"
    "os"
    "fmt"
    "os/signal"
    _ "github.com/andir/lightsd/operations/lua"
    _ "github.com/andir/lightsd/operations/rainbow"
    _ "github.com/andir/lightsd/operations/raindrop"
    _ "github.com/andir/lightsd/operations/rotation"
    _ "github.com/andir/lightsd/outputs/shm"
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

    pipelines := BuildPipelines(config.Pipelines)

    mqtt := NewMqttConnection(config)
    for _, pipeline := range pipelines {
        mqtt.Register(pipeline)
    }

    bc := StartDebug()

    go func() {
        interval := time.Second / time.Duration(config.FPS)
        lastTime := time.Now()
        for {
            currTime := time.Now()
            duration := currTime.Sub(lastTime)

            for _, pipeline := range pipelines {
                pipeline.Render(duration)
            }

            bc.Broadcast(pipelines[0])

            // Wait until next frame should start
            time.Sleep(lastTime.Add(interval).Sub(time.Now()))

            lastTime = currTime
        }
    }()

    signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, os.Interrupt)
    for range signalChan {
        fmt.Println("Interrupt received. Stopping")
        return
    }
}