package psg

import (
	"bytes"
	"io"
)

func Compress(src io.Reader, dst io.Writer) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(src)
	s := buf.Bytes()

	l.Printf("uncompressed file size = %d bytes\n", len(s))

	blkmap := make([]*Block, len(s)) // holds an image of the input file, to keep track of already converted blocks
	blks := make([]*Block, 0)        // holds a simple list of all known blocks
	// keep looping until we can't find anything to compress
	for {
		dict := make(map[string]Block) // for speed, keep track of all blocks found this pass
		best := Block{}
		// phase 1 - loop through the whole input file to find the best block
		for i := 0; i < len(s)-MIN; i++ {
			// try every possible length of block from 4 to 51
			for w := MIN; w <= MAX && w+i < len(s); w++ {
				needle := make([]byte, w)
				copy(needle, s[i:i+w])
				_, exists := dict[string(needle)]
				// if we've already seen this block, don't process it again
				if !exists {
					foundAt := []int{}
					from := i + w
					// find all subsequent occurrences of this block
					for len(s) >= len(needle) {
						at := bytes.Index(s[from:], needle)
						if at < 0 {
							break
						}
						// exclude any occurrences which overlap with already replaced data
						skip := false
						for _, l := range blkmap[from+at : from+at+len(needle)] {
							if l != nil {
								skip = true
								break
							}
						}
						if !skip {
							foundAt = append(foundAt, from+at)
							from += at + len(needle)
						} else {
							from += at + 1
						}
					}
					if len(foundAt) == 0 {
						break
					}
					blk := Block{needle, i, len(needle), foundAt}
					dict[string(needle)] = blk
					// if the block's compression potential is higher than the others, choose it as the best
					if blk.Score() > best.Score() {
						best = blk
					}
				}
			}
		}
		l.Printf("found %d repeated blocks\n", len(dict))
		if best.Score() <= 0 {
			l.Println("finishing")
			break
		}
		l.Printf("best block is at %.4x, length %d with score %d (%x)\n", best.Start, best.Len, best.Score(), best.Bytes)
		// we are going to do some replacements with our winning block - update blks and blkmap to protect the block from future updates
		blks = append(blks, &best)
		for i := best.Start; i < best.Start+best.Len; i++ {
			blkmap[i] = &best
		}
		// phase 2 - now action all the replacements we noted we could make with the block
		for _, at := range best.At {
			skipped := best.Len - 3
			l.Printf("replacing block of %d bytes at %.4x with vector to %.4x, saving %d bytes (%x)\n", best.Len, at, best.Start, skipped, best.Bytes)
			copy(s[at:], best.Encode())
			s = append(s[:at+3], s[at+best.Len:]...)
			copy(blkmap[at:], []*Block{&best, &best, &best})
			blkmap = append(blkmap[:at+3], blkmap[at+best.Len:]...)
			// adjust any offsets of replacements we have already made, if they occur after the point at which we just shortened the file
			for _, blk := range blks {
				rebased := false
				if blk.Start > at+3 {
					blk.Start -= skipped
					rebased = true
				}
				for i, a := range blk.At {
					if a > at+3 {
						blk.At[i] -= skipped
					}
					if rebased {
						l.Printf("correcting vector at %.4x from %x to %x as target moved\n", blk.At[i], s[blk.At[i]:blk.At[i]+3], blk.Encode())
						copy(s[blk.At[i]:], blk.Encode())
					}
				}
			}
		}
	}

	l.Printf("compressed file size = %d bytes\n", len(s))

	dst.Write(s)
}
