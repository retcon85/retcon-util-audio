package psg

import (
	"bytes"
	"encoding/binary"
	"io"
)

func Decompress(src io.Reader, dst io.Writer) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(src)
	s := buf.Bytes()

	for i := 0; i < len(s); i++ {
		if s[i] >= 0x08 && s[i] <= 0x37 {
			from := int(binary.LittleEndian.Uint16(s[i+1 : i+3]))
			dst.Write(s[from : from+int(s[i]-4)])
			i += 2
			continue
		}
		dst.Write(s[i : i+1])
	}
}
