package psg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

const (
	NTSC_CLK = 3579545
	PAL_CLK  = 3546893
)

type DebugOptions struct {
	PrintOffset bool
	PrintBytes  bool
	ShowFrames  bool
}

func write(w io.Writer, off int, b []byte, latch bool, msg string, opt DebugOptions) {
	if opt.PrintOffset {
		fmt.Fprintf(w, "%.4x:    ", off)
	}
	if latch {
		fmt.Fprintf(w, "* ")
	} else {
		fmt.Fprintf(w, "  ")
	}
	if opt.PrintBytes {
		fmt.Fprintf(w, "%- #15x", b)
	}
	fmt.Fprintln(w, msg)
}

func logRegisterChange(w io.Writer, off int, b byte, reg int, val int, opt DebugOptions) {
	latch := b&0x80 > 0
	switch {
	case reg&1 > 0:
		ch := reg >> 1
		attn := val & 0x0f
		if attn == 15 {
			write(w, off, []byte{b}, latch, fmt.Sprintf("channel %d attenuation => %5d (silent)", ch, attn), opt)
		} else {
			db := int(attn) * 2
			write(w, off, []byte{b}, latch, fmt.Sprintf("channel %d attenuation => %5d (%3d db)", ch, attn, db), opt)
		}
	case reg&0x6 == 0x6:
		var fb, nf string
		if val&(0x04) > 0 {
			fb = "white   "
		} else {
			fb = "periodic"
		}
		switch val & 0x3 {
		case 0:
			nf = "ϕ/2"
		case 1:
			nf = "ϕ/4"
		case 2:
			nf = "ϕ/8"
		case 3:
			nf = "ch 2 freq"
		}
		write(w, off, []byte{b}, latch, fmt.Sprintf("noise                 => %s @ %s", fb, nf), opt)
	default:
		ch := reg >> 1
		f := val & 0x3ff
		hz := NTSC_CLK / (32 * float32(f))
		hzk := " "
		if hz >= 1000 {
			hz /= 1000
			hzk = "k"
		}
		write(w, off, []byte{b}, latch, fmt.Sprintf("channel %d tone        => %5d (%6.2f %sHz)", ch, f, hz, hzk), opt)
	}
}

func Debug(src io.Reader, dst io.Writer, opt DebugOptions) error {
	s := new(bytes.Buffer)
	s.ReadFrom(src)
	buf := s.Bytes()
	s = bytes.NewBuffer(buf)
	off := -1
	var regs [8]int
	reg := 0
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
			reg = int((b >> 4) & 0x7)
			regs[reg] = (regs[reg] &^ 0x00f) | int(b&0x0f)
			logRegisterChange(dst, off, b, reg, regs[reg], opt)
		case b&0x40 == 0x40:
			regs[reg] = (regs[reg] & 0x00f) | int(b)<<4
			logRegisterChange(dst, off, b, reg, regs[reg], opt)
		case b&0x38 == 0x38:
			delay := b & 0x07
			pl := "s"
			if delay == 1 {
				pl = " "
			}
			if opt.ShowFrames {
				for i := 0; i <= int(delay); i++ {
					ln := strings.Repeat("-", 25)
					fmt.Fprintf(dst, "%s wait for %2d frame%s (%d of %d) %s\n", ln, delay, pl, i, delay, ln)
				}
			} else {
				write(dst, off, []byte{b}, false, fmt.Sprintf("wait for %d frame%s", delay, pl), opt)
			}
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
			write(dst, off, []byte{b, vec[0], vec[1]}, false, fmt.Sprintf("repeat block from %.4x:%.4x (%d bytes)", from, to, size), opt)
			off += 2
		}
	}
}
