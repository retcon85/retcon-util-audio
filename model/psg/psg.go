package psg

import (
	"encoding/binary"

	"github.com/willbritton/gocli"
)

var log = gocli.Log
var dbg = gocli.Dbg

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
