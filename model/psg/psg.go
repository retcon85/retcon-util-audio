package psg

import (
	"encoding/binary"

	"github.com/retcon85/retcon-util-audio/logger"
)

var l = logger.DefaultLogger()

const (
	MIN = 4
	MAX = 51
)

type Block struct {
	Bytes []byte
	Start int
	Len   int
	At    []int
}

func (b Block) Encode() []byte {
	e := []byte{byte(b.Len + 4), 0, 0}
	binary.LittleEndian.PutUint16(e[1:], uint16(b.Start))
	return e
}

func (b Block) Score() int {
	return len(b.At) * (b.Len - 3)
}
