package shm

import (
    "github.com/andir/lightsd/core"
    "github.com/andir/lightsd/outputs"
    "unsafe"
    "os"
    "syscall"
    "runtime"
    "github.com/fabiokung/shm"
    "reflect"
    "log"
)

type SHMOutputConfig struct {
    Path string `mapstructure:"file"`
}

type SHMOutput struct {
    source   string
    fh       *os.File
    mmap     []byte
    count    int
    filename string
}

func destroySHMOutput(s *SHMOutput) {
    if s.mmap != nil {
        if err := syscall.Munmap(s.mmap); err != nil {
            log.Println("shm: munmap failed:", err)
        }
    }

    if s.fh != nil {
        if err := s.fh.Close(); err != nil {
            log.Println("shm: close failed:", err)
        }
    }
    shm.Unlink(s.filename)
}

func newSHMOutput(filename string, source string, count int) (*SHMOutput, error) {
    map_file, err := shm.Open(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
        return nil, err
    }

    if err := map_file.Truncate(int64(count * 4)); err != nil {
        return nil, err
    }

    mmap, err := syscall.Mmap(int(map_file.Fd()), 0, int(count * 4), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
    if err != nil {
        return nil, err
    }

    shm := &SHMOutput{
        source:   source,
        count:    count,
        fh:       map_file,
        mmap:     mmap,
        filename: filename,
    }

    runtime.SetFinalizer(shm, destroySHMOutput)

    return shm, nil
}

func (m *SHMOutput) Render(stripe core.LEDStripeReader) {
    map_array := (*[1 << 30]uint8)(unsafe.Pointer(&m.mmap[0]))[:m.count*4:m.count*4]

    // TODO: Add temporal dithering (see https://github.com/FastLED/FastLED/wiki/FastLED-Temporal-Dithering)
    // TODO: Make this available for other outputs, too

    for i := 0; i < m.count; i++ {
        map_array[i*4+0] = byte(0)
        map_array[i*4+1], map_array[i*4+2], map_array[i*4+3] = stripe.Get(i).RGB255()
    }
}

func (m *SHMOutput) Source() string {
    return m.source
}

func init() {
    outputs.Register("shm", &outputs.Factory{
        ConfigType: reflect.TypeOf(SHMOutputConfig{}),
        Create: func(count int, source string, rconfig interface{}) (core.Output, error) {
            config := rconfig.(*SHMOutputConfig)

            return newSHMOutput(config.Path, source, count)
        },
    })
}
