package main

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "fmt"
)

type Config struct {
    FPS uint
    Size int

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

    Outputs []struct {
        Name string
        Size uint16
    }

    Operations []struct {
        Name string
        Type string
        Config map[string]interface{}
    }
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
