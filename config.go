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

    return config, nil
}

func BuildOutput(config PipelineConfig) (core.Output, error) {
    f := outputs.Get(config.Output.Type)

    c := reflect.New(f.ConfigType).Interface()
    err := mapstructure.Decode(config.Output.Config, c)
    if err != nil {
        return nil, fmt.Errorf("config: Unknown output type: %s", config.Output.Type)
    }

    output, err := f.Create(config.Count, config.Output.Operation, c)
    if err != nil {
        return nil,fmt.Errorf("config: Failed to build output config: %v", err)
    }

    return output, nil
}

func BuildOperation(config PipelineConfig, i int) (core.Operation, error) {
    f := operations.Get(config.Operations[i].Type)
    if f == nil {
        return nil, fmt.Errorf("config: Unknown operation type: %s", config.Operations[i].Type)
    }

    c := reflect.New(f.ConfigType).Interface()
    err := mapstructure.Decode(config.Operations[i].Config, c)
    if err != nil {
        return nil,fmt.Errorf("config: Failed to build operation config: %s: %v", config.Operations[i].Name, err)
    }

    output, err := f.Create(config.Operations[i].Name, config.Count, c)
    if err != nil {
        return nil,fmt.Errorf("config: Failed to create operation: %v", err)
    }

    return output, nil
}

func BuildPipeline(name string, config PipelineConfig) (*core.Pipeline, error) {
    output, err := BuildOutput(config)
    if err != nil {
        return nil, err
    }

    operations := make([]core.Operation, 0, len(config.Operations))
    for i := range config.Operations {
        op, err := BuildOperation(config, i)
        if err != nil {
            return nil, err
        }

        operations = append(operations, op)
    }

    return core.NewPipeline(name, config.Count, output, operations), nil
}

func BuildPipelines(config map[string]PipelineConfig) ([]*core.Pipeline, error) {
    pipelines := make([]*core.Pipeline, 0, len(config))
    for name, config := range config {
        pipeline, err := BuildPipeline(name, config)
        if err != nil {
            return nil, err
        }

        pipelines = append(pipelines, pipeline)
    }

    return pipelines, nil
}
