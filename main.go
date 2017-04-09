package main

import (
    "log"
    "flag"
    "os"
    "os/signal"
    "runtime/pprof"
    "time"
    _ "github.com/andir/lightsd/operations/gradient"
    _ "github.com/andir/lightsd/operations/raindrop"
    _ "github.com/andir/lightsd/operations/rotation"
    _ "github.com/andir/lightsd/operations/blackout"
    _ "github.com/andir/lightsd/operations/dimmer"
    _ "github.com/andir/lightsd/operations/adjustment"
    _ "github.com/andir/lightsd/outputs/shm"
    "github.com/andir/lightsd/operations"
    "github.com/andir/lightsd/outputs"
    "github.com/andir/lightsd/debug"
    "github.com/andir/lightsd/core"
)

func buildPipelines(config map[string]PipelineConfig) ([]*core.Pipeline, error) {
    pipelines := make([]*core.Pipeline, 0, len(config))
    for name, config := range config {
        out, err := outputs.Make(config.Output.Type, config.Count, config.Output.Source, config.Output.Config)
        if err != nil {
            return nil, err
        }

        ops := make([]core.Operation, 0, len(config.Operations))
        for i := range config.Operations {
            op, err := operations.Make(config.Operations[i].Type, config.Operations[i].Name, config.Count, config.Operations[i].Config)
            if err != nil {
                return nil, err
            }

            ops = append(ops, op)
        }

        pipelines = append(pipelines, core.NewPipeline(name, config.Count, out, ops))
    }

    return pipelines, nil
}

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

    pipelines, err := buildPipelines(config.Pipelines)
    if err != nil {
        panic(err)
    }

    mqtt, err := NewMqttConnection(config)
    if err != nil {
        panic(err)
    }

    for _, pipeline := range pipelines {
        err = mqtt.Register(pipeline)
        if err != nil {
            panic(err)
        }
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
        log.Print("Interrupt received - stopping")
        return
    }
}
