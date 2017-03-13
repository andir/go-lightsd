package lua

import (
    "github.com/andir/lightsd/core"
    lua "github.com/yuin/gopher-lua"
    "reflect"
    "sync"
    "time"
    "fmt"
)

type LuaScriptConfig struct {
    Path string `mapstructure:"path"`
}

type LuaScript struct {
    sync.RWMutex

    name   string
    stripe core.LEDStripe

    state *lua.LState

    fnUpdate lua.LValue
    fnRender lua.LValue
}

func (this *LuaScript) Name() string {
    return this.name
}

func (this *LuaScript) Stripe() core.LEDStripe {
    return this.stripe
}

func (this *LuaScript) Update(duration time.Duration) {
    if err := this.state.CallByParam(lua.P{
        Fn:      this.fnUpdate,
        Protect: true,
    }, lua.LNumber(duration.Seconds())); err != nil {
        fmt.Printf("Error in script / update: %v", err)
    }
}

func (this *LuaScript) Render() {
    if err := this.state.CallByParam(lua.P{
        Fn:      this.fnRender,
        Protect: true,
    }); err != nil {
        fmt.Printf("Error in script / render: %v", err)
    }
}

func init() {
    core.RegisterOperation("lua", core.OperationFactory{
        ConfigType: reflect.TypeOf(LuaScriptConfig{}),
        Create: func(pipeline *core.Pipeline, name string, count int, rconfig interface{}) (core.Operation, error) {
            config := rconfig.(*LuaScriptConfig)

            stripe := core.NewLEDStripe(count)

            state := lua.NewState()
            state.SetGlobal("count", lua.LNumber(count))
            state.SetGlobal("put", state.NewFunction(func(l *lua.LState) int {
                i := int(l.ToNumber(1))
                stripe[i].R = uint8(l.ToNumber(2))
                stripe[i].G = uint8(l.ToNumber(3))
                stripe[i].B = uint8(l.ToNumber(4))
                stripe[i].A = 0

                return 0
			}))

            err := state.DoFile(config.Path)
            if err != nil {
                return nil, err
            }

            fnInit := state.GetGlobal("init")
            fnUpdate := state.GetGlobal("update")
            fnRender := state.GetGlobal("render")

            if err := state.CallByParam(lua.P{
                Fn:      fnInit,
                Protect: true,
            }); err != nil {
                return nil, err
            }

            return &LuaScript{
                name:   name,
                stripe: stripe,

                state: state,

                fnUpdate: fnUpdate,
                fnRender: fnRender,
            }, nil
        },
    })
}
