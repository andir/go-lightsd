package main

import (
    "flag"
    "time"
    "runtime/pprof"
    "os"
    "fmt"
    "os/signal"
    _ "github.com/andir/lightsd/operations/gradient"
    _ "github.com/andir/lightsd/operations/raindrop"
    _ "github.com/andir/lightsd/operations/rotation"
    _ "github.com/andir/lightsd/operations/blackout"
    _ "github.com/andir/lightsd/outputs/shm"
    "github.com/andir/lightsd/debug"
    "github.com/andir/lightsd/core"
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

    var debugger *debug.Debugger = nil
    if config.Debug.Enable {
        debugger = debug.StartDebug(config.Debug.Port, pipelines)
    }

    for _, pipeline := range pipelines {
        go func(pipeline *core.Pipeline) {
            interval := time.Second / time.Duration(config.FPS)
            lastTime := time.Now()

            for {
                currTime := time.Now()
                duration := currTime.Sub(lastTime)

                context := pipeline.Render(duration)

                if debugger != nil {
                    debugger.Broadcast(context)
                }

                // Wait until next frame should start
                time.Sleep(lastTime.Add(interval).Sub(time.Now()))

                lastTime = currTime
            }
        }(pipeline)
    }

    signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, os.Interrupt)
    for range signalChan {
        fmt.Println("Interrupt received. Stopping")
        return
    }
}
