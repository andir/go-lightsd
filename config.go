package main

import (
    "fmt"
    "reflect"
    "io/ioutil"
    "gopkg.in/yaml.v2"
    "github.com/mitchellh/mapstructure"
    "github.com/andir/lightsd/core"
    "github.com/andir/lightsd/operations"
    "github.com/andir/lightsd/outputs"
)

type OperationConfig struct {
    Name string
    Type string

    Config map[string]interface{}
}

type OutputConfig struct {
    Type      string
    Operation string

    Config map[string]interface{}
}

type PipelineConfig struct {
    Count int

    Output     OutputConfig
    Operations []OperationConfig
}

type Config struct {
    FPS uint

    MQTT struct {
        Host string
        Port int

        ClientID string

        Realm string
    }

    Debug struct {
        Enable bool
        Port   int
    }

    Pipelines map[string]PipelineConfig
}

func LoadConfig(path string) (*Config, error) {
    b, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }

    config := new(Config)
    err = yaml.Unmarshal(b, config)
    if err != nil {
        return nil, err
    }

    out, err := yaml.Marshal(config)
    fmt.Printf("Config: %s\n", out)

    return config, nil
}

func BuildOutput(config PipelineConfig) core.Output {
    f := outputs.Get(config.Output.Type)

    c := reflect.New(f.ConfigType).Interface()
    err := mapstructure.Decode(config.Output.Config, c)
    if err != nil {
        panic(err)
    }

    output, err := f.Create(config.Count, config.Output.Operation, c)
    if err != nil {
        panic(err)
    }

    return output
}

func BuildOperation(config PipelineConfig, i int) core.Operation {
    f := operations.Get(config.Operations[i].Type)

    c := reflect.New(f.ConfigType).Interface()
    err := mapstructure.Decode(config.Operations[i].Config, c)
    if err != nil {
        panic(err)
    }

    output, err := f.Create(config.Operations[i].Name, config.Count, c)
    if err != nil {
        panic(err)
    }

    return output
}

func BuildPipeline(name string, config PipelineConfig) *core.Pipeline {
    output := BuildOutput(config)

    operations := make([]core.Operation, 0, len(config.Operations))
    for i := range config.Operations {
        operations = append(operations, BuildOperation(config, i))
    }

    return core.NewPipeline(name, config.Count, output, operations)
}

func BuildPipelines(config map[string]PipelineConfig) []*core.Pipeline {
    pipelines := make([]*core.Pipeline, 0, len(config))
    for name, config := range config {
        pipelines = append(pipelines, BuildPipeline(name, config))
    }

    return pipelines
}
