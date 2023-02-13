package psg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	NTSC_CLK = 3579545
	PAL_CLK  = 3546893
)

type DebugOptions struct {
	PrintOffset bool
	PrintBytes  bool
}

func write(w io.Writer, off int, b []byte, data_byte bool, msg string, opt DebugOptions) {
	data_str := "    "
	if data_byte {
		data_str = "  * "
	}
	if opt.PrintOffset {
		fmt.Fprintf(w, "%.4x:\t", off)
	}
	if opt.PrintBytes {
		fmt.Fprintf(w, "%- #12x\t", b)
	}
	fmt.Fprintf(w, "%s%s\n", data_str, msg)
}

func log2ByteCommand(w io.Writer, off int, b1 byte, b2 byte, data_byte bool, opt DebugOptions) {
	b := b1
	if data_byte {
		b = b2
	}
	switch {
	case b1&(0x80|0x10) == 0x80|0x10:
		b1 |= (b2 & 0x0f)
		ch := (b1 >> 5) & 0x03
		write(w, off, []byte{b}, data_byte, fmt.Sprintf("channel %d attenuation => %2d", ch, b1&0x0f), opt)
	case b1&(0x80|0x60) == 0x80|0x60:
		b1 |= (b2 & 0x0f)
		var fb, nf string
		if b1&(0x04) > 0 {
			fb = "white noise   "
		} else {
			fb = "periodic noise"
		}
		switch b1 & 0x3 {
		case 0:
			nf = "[clk/2]"
		case 1:
			nf = "[clk/4]"
		case 2:
			nf = "[clk/8]"
		case 3:
			nf = "[channel 2]"
		}
		write(w, off, []byte{b}, data_byte, fmt.Sprintf("play %s    with frequency %18s", fb, nf), opt)
	case b1&0x80 == 0x80:
		ch := (b1 >> 5) & 0x03
		f := int(b2&0x3f)<<4 + int(b1&0x0f)
		hz := NTSC_CLK / (32 * float32(f))
		hzk := " "
		if hz >= 1000 {
			hz /= 1000
			hzk = "k"
		}
		write(w, off, []byte{b}, data_byte, fmt.Sprintf("play tone on channel %d with frequency %5d (%6.2f %sHz)", ch, f, hz, hzk), opt)
	}
}

func Debug(src io.Reader, dst io.Writer, opt DebugOptions) error {
	s := new(bytes.Buffer)
	s.ReadFrom(src)
	buf := s.Bytes()
	s = bytes.NewBuffer(buf)
	off := -1
	var b1, b2 byte
	for {
		b, err := s.ReadByte()
		off++
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		switch {
		case b&0x80 == 0x80:
			b1 = b
			log2ByteCommand(dst, off, b1, b2, false, opt)
		case b&0x40 == 0x40:
			b2 = b
			log2ByteCommand(dst, off, b1, b2, true, opt)
		case b&0x38 == 0x38:
			delay := b & 0x07
			pl := "s"
			if delay == 1 {
				pl = ""
			}
			write(dst, off, []byte{b}, false, fmt.Sprintf("wait for %d frame%s", delay, pl), opt)
		case b == 0:
			write(dst, off, []byte{b}, false, "end of file", opt)
		case b == 1:
			write(dst, off, []byte{b}, false, "loop marker", opt)
		default:
			size := uint16(b - 4)
			if size < 4 || size > 51 {
				panic(fmt.Sprintf("compression size was out of range, got %d reading byte %#.2x at offset %.8x", size, b, off))
			}
			vec := make([]byte, 2)
			for i := range vec {
				b, err := s.ReadByte()
				if err == io.EOF {
					return fmt.Errorf("unexpected EOF")
				} else if err != nil {
					return err
				}
				vec[i] = b
			}
			from := binary.LittleEndian.Uint16(vec)
			to := from + size
			write(dst, off, []byte{b, vec[0], vec[1]}, false, fmt.Sprintf("repeat block from %.4x:%.4x (%d bytes <= % x)", from, to, size, buf[from:to]), opt)
			off += 2
		}
	}
}
