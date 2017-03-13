package shm

import (
    "github.com/andir/lightsd/core"
    "github.com/andir/lightsd/outputs"
    "unsafe"
    "os"
    "fmt"
    "syscall"
    "runtime"
    "github.com/fabiokung/shm"
    "reflect"
)

type SHMOutputConfig struct {
	Path string `mapstructure:"file"`
}

type SHMOutput struct {
	fh   *os.File
	mmap []byte
	size int
	filename string
}

func destroySHMOutput(s *SHMOutput) {
	if s.mmap != nil {
		if err := syscall.Munmap(s.mmap); err != nil {
			fmt.Println("munmap:", err)
		}
	}

	if s.fh != nil {
		if err := s.fh.Close(); err != nil {
			fmt.Println("close:", err)
		}
	}
	shm.Unlink(s.filename)
}

func newSHMOutput(filename string, count int) (*SHMOutput, error) {
	size := count * 4

	map_file, err := shm.Open(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, err
	}

	if err := map_file.Truncate(int64(size)); err != nil {
		return nil, err
	}

	mmap, err := syscall.Mmap(int(map_file.Fd()), 0, int(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	shm := &SHMOutput{
		size: size,
		fh:   map_file,
		mmap: mmap,
		filename: filename,
	}

	runtime.SetFinalizer(shm, destroySHMOutput)

	return shm, nil
}

func minInt(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func (m *SHMOutput) Render(stripe core.LEDStripe) {
	map_array := (*[1<<30]uint8)(unsafe.Pointer(&m.mmap[0]))[:m.size:m.size]
	n := minInt(len(stripe), m.size)

	for i, p := range stripe {
		if i > n {
			break
		}
		map_array[i*4+0] = byte(0)
		map_array[i*4+1] = byte(p.G)
		map_array[i*4+2] = byte(p.R)
		map_array[i*4+3] = byte(p.B)
	}
}

func init() {
	outputs.Register("shm", &outputs.Factory{
		ConfigType: reflect.TypeOf(SHMOutputConfig{}),
		Create: func(count int, rconfig interface{}) (core.Output, error) {
			config := rconfig.(*SHMOutputConfig)

			return newSHMOutput(config.Path, count)
		},
	})
}