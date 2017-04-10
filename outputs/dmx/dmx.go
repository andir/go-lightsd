package dmx

// +build linux,cgo

// #cgo pkg-config: libftdi
// #include <ftdi.h>
//
//void render(struct ftdi_context* ctx, unsigned char* buf, int size) {
//  ftdi_set_line_property2(ctx, 8, 2, 0, 1);
//  usleep(10000);
//  ftdi_set_line_property2(ctx, 8, 2, 0, 0);
//  usleep(8000);
//  ftdi_write_data(ctx, buf, size);
//}
import "C"

import (
    "github.com/andir/lightsd/core"
    "github.com/andir/lightsd/outputs"
    "reflect"
    "errors"
    "unsafe"
)

type DMXOutputConfig struct {
    VendorId int `mapstructure:"vendor"`
    DeviceId int `mapstructure:"device"`

    Index int `mapstructure:"index"`

    Offset int `mapstructure:"offset"`
    Align  int `mapstructure:"align"`
}

type DMXOutput struct {
    source string

    count int

    offset int
    align  int

    buffer []byte

    ctx *C.struct_ftdi_context
}

func (this *DMXOutput) Render(stripe core.LEDStripeReader) {
    // Fill the buffer
    for i := 0; i < this.count; i++ {
        o := 1 + this.offset + i*this.align
        this.buffer[o+0], this.buffer[o+1], this.buffer[o+2] = stripe.Get(i).RGB255()
    }

    C.render(this.ctx, (*C.uchar)(unsafe.Pointer(&this.buffer[0])), C.int(len(this.buffer)))
}

func (this *DMXOutput) Source() string {
    return this.source
}

func init() {
    outputs.Register("dmx", &outputs.Factory{
        ConfigType: reflect.TypeOf(DMXOutputConfig{}),
        Create: func(count int, source string, rconfig interface{}) (core.Output, error) {
            config := rconfig.(*DMXOutputConfig)

            if config.Align == 0 {
                config.Align = 3
            }

            if count > (512-config.Offset)/config.Align {
                return nil, errors.New("To many LEDs for a DMX universe")
            }

            // Create FTDI context
            var ctx C.struct_ftdi_context
            if res := C.ftdi_init(&ctx); res != 0 {
                return nil, errors.New("Failed to create FTDI context: " + C.GoString(C.ftdi_get_error_string(&ctx)))
            }
            //defer C.ftdi_free(ctx)

            // Open and configure device
            if res := C.ftdi_usb_open_desc_index(&ctx, C.int(config.VendorId), C.int(config.DeviceId), nil, nil, 0 /*C.uint(config.Index)*/); res != 0 {
                return nil, errors.New("Failed to open FTDI device: " + C.GoString(C.ftdi_get_error_string(&ctx)))
            }

            if res := C.ftdi_usb_reset(&ctx); res != 0 {
                return nil, errors.New("Failed to configure FTDI device: reset: " + C.GoString(C.ftdi_get_error_string(&ctx)))
            }

            if res := C.ftdi_set_baudrate(&ctx, 250000); res != 0 {
                return nil, errors.New("Failed to configure FTDI device: baudrate: " + C.GoString(C.ftdi_get_error_string(&ctx)))
            }

            if res := C.ftdi_set_line_property2(&ctx, 8, 2, 0, 0); res != 0 {
                return nil, errors.New("Failed to configure FTDI device: line properties: " + C.GoString(C.ftdi_get_error_string(&ctx)))
            }

            if res := C.ftdi_setflowctrl(&ctx, 0); res != 0 {
                return nil, errors.New("Failed to configure FTDI device: flowcontrol: " + C.GoString(C.ftdi_get_error_string(&ctx)))
            }

            if res := C.ftdi_setrts(&ctx, 0); res != 0 {
                return nil, errors.New("Failed to configure FTDI device: RTS: " + C.GoString(C.ftdi_get_error_string(&ctx)))
            }

            if res := C.ftdi_usb_purge_buffers(&ctx); res != 0 {
                return nil, errors.New("Failed to configure FTDI device: purge: " + C.GoString(C.ftdi_get_error_string(&ctx)))
            }

            buffer := make([]byte, 1+config.Offset+count*config.Align) // The DMX universe to render
            buffer[0] = 0x00                                           // and an initial slot which is always zero

            return &DMXOutput{
                source: source,

                count: count,

                offset: config.Offset,
                align:  config.Align,

                buffer: buffer,

                ctx: &ctx,
            }, nil
        },
    })
}
