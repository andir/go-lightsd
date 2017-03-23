package main

import (
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

type OperationConfig struct {
    Name string
    Type string

    Config map[string]interface{}
}

type OutputConfig struct {
    Type   string
    Source string

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
