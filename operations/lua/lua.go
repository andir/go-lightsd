package lua

import (
    "github.com/andir/lightsd/core"
    "github.com/andir/lightsd/operations"
    lua "github.com/yuin/gopher-lua"
    "reflect"
    "sync"
    "fmt"
)

type LuaScriptConfig struct {
    Path string `mapstructure:"file"`
}

type LuaScript struct {
    sync.RWMutex

    name string

    state *lua.LState
    fn    lua.LValue
}

func (this *LuaScript) Name() string {
    return this.name
}

func (this *LuaScript) Render(context *core.RenderContext) {
    if err := this.state.CallByParam(lua.P{
        Fn:      this.fn,
        Protect: true,
    }); err != nil {
        fmt.Printf("Error in script / render: %v", err)
    }
}

func init() {
    operations.Register("lua", &operations.Factory{
        ConfigType: reflect.TypeOf(LuaScriptConfig{}),
        Create: func(name string, count int, rconfig interface{}) (core.Operation, error) {
            config := rconfig.(*LuaScriptConfig)

            state := lua.NewState()
            state.SetGlobal("count", lua.LNumber(count))
            //state.SetGlobal("put", state.NewFunction(func(l *lua.LState) int {
            //    i := int(l.ToNumber(1))
            //    stripe[i].R = uint8(l.ToNumber(2))
            //    stripe[i].G = uint8(l.ToNumber(3))
            //    stripe[i].B = uint8(l.ToNumber(4))
            //    stripe[i].A = 0
            //
            //    return 0
            //}))

            err := state.DoFile(config.Path)
            if err != nil {
                return nil, err
            }

            fn := state.GetGlobal("render")

            return &LuaScript{
                name:   name,

                state: state,
                fn: fn,
            }, nil
        },
    })
}
