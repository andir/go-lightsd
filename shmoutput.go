package main

import (
	"unsafe"
	"os"
	"fmt"
	"syscall"
	"runtime"
	"github.com/fabiokung/shm"
	"image/color"
)

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

func NewSHMOutput(filename string, size int) *SHMOutput {

	size *= 4

	map_file, err := shm.Open(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := map_file.Truncate(int64(size)); err != nil {
		fmt.Println("truncate:", err)
		os.Exit(1)
	}

	mmap, err := syscall.Mmap(int(map_file.Fd()), 0, int(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		fmt.Println("mmap:", err)
		os.Exit(1)
	}
	//map_array := (*[n]int)(unsafe.Pointer(&mmap[0]))
	shm := &SHMOutput{
		size: size,
		fh:   map_file,
		mmap: mmap,
		filename: filename,
	}

	runtime.SetFinalizer(shm, destroySHMOutput)

	return shm
}

func minInt(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func (m *SHMOutput) Render(l []color.RGBA) {
	map_array := (*[1<<30]uint8)(unsafe.Pointer(&m.mmap[0]))[:m.size:m.size]
	n := minInt(len(l), m.size)

	for i, p := range l {
		if i > n {
			break
		}

		map_array[i+0] = byte(p.R)
		map_array[i+1] = byte(p.G)
		map_array[i+2] = byte(p.B)
		map_array[i+3] = byte(0)
	}
}