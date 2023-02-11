package psg_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/retcon85/retcon-util-audio/model/psg"
)

func TestSearchForBlockSimple(t *testing.T) {
	s := bytes.NewBuffer([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 5, 6, 7, 8, 9, 10})
	buff := new(bytes.Buffer)
	psg.Compress(s, buff)
	got := buff.Bytes()
	want := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 0x0a, 0x04, 0x00}
	if len(want) != len(got) {
		t.Fatalf("result wrong length, expected:\n%d, got:\n%d", want, got)
	}
	for i, wanted := range want {
		if got[i] != wanted {
			t.Fatalf("failed at index %d, expected:\n%d, got:\n%d", i, want, got)
		}
	}
}

func TestSearchForBlockRepeated(t *testing.T) {
	s := bytes.NewBuffer([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 5, 6, 7, 8, 9, 10, 11, 12, 5, 6, 7, 8, 9, 10})
	buff := new(bytes.Buffer)
	psg.Compress(s, buff)
	got := buff.Bytes()
	want := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 0x0a, 0x04, 0x00, 11, 12, 0x0a, 0x04, 0x00}
	if len(want) != len(got) {
		t.Fatalf("result wrong length, expected:\n%d, got:\n%d", want, got)
	}
	for i, wanted := range want {
		if got[i] != wanted {
			t.Fatalf("failed at index %d, expected:\n%d, got:\n%d", i, want, got)
		}
	}
}

func TestSearchForBlockRepeatedPlusAnotherBlock(t *testing.T) {
	s := bytes.NewBuffer([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 2, 3, 4, 5, 6, 14, 2, 3, 4, 5, 6, 15, 5, 6, 7, 8, 9, 16})
	buff := new(bytes.Buffer)
	psg.Compress(s, buff)
	got := buff.Bytes()
	want := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 0x09, 0x01, 0x00, 14, 0x09, 0x01, 0x00, 15, 0x09, 0x04, 0x00, 16}
	if len(want) != len(got) {
		t.Fatalf("result wrong length, expected:\n%d, got:\n%d", want, got)
	}
	for i, wanted := range want {
		if got[i] != wanted {
			t.Fatalf("failed at index %d, expected:\n%d, got:\n%d", i, want, got)
		}
	}
}

func TestSearchForBlockDoesNotOverwriteTokens(t *testing.T) {
	s := bytes.NewBuffer([]byte{1, 2, 3, 4, 5, 1, 2, 3, 4, 1, 2, 3, 4, 5, 1, 2, 3, 4})
	buff := new(bytes.Buffer)
	psg.Compress(s, buff)
	got := buff.Bytes()
	want := []byte{1, 2, 3, 4, 5, 1, 2, 3, 4, 0x0d, 0x00, 0x00}
	if len(want) != len(got) {
		t.Fatalf("result wrong length, expected:\n%d, got:\n%d", want, got)
	}
	for i, wanted := range want {
		if got[i] != wanted {
			t.Fatalf("failed at index %d, expected:\n%d, got:\n%d", i, want, got)
		}
	}
}

func TestSearchForBlockDoesNotOverwriteTokensInMiddle(t *testing.T) {
	s := bytes.NewBuffer([]byte{8, 8, 8, 8, 1, 2, 3, 4, 8, 8, 8, 8, 8, 8, 8, 8, 5, 6, 8, 8, 8, 8, 7, 8, 9, 8, 8, 8, 8, 8, 8, 8, 8, 7, 6, 5, 4})
	buff := new(bytes.Buffer)
	psg.Compress(s, buff)
	got := buff.Bytes()
	want := []byte{8, 8, 8, 8, 1, 2, 3, 4, 8, 0x00, 0x00, 8, 0x00, 0x00, 5, 6, 0x08, 0x00, 0x00, 7, 8, 9, 0x08, 0x00, 0x00, 0x08, 0x00, 0x00, 7, 6, 5, 4}
	if len(want) != len(got) {
		t.Fatalf("result wrong length, expected:\n%d, got:\n%d", want, got)
	}
	for i, wanted := range want {
		if got[i] != wanted {
			t.Fatalf("failed at index %d, expected:\n%d, got:\n%d", i, want, got)
		}
	}
}

func TestRealFile(t *testing.T) {
	fn := "../../testdata/primates_reference.psg"
	stats, err := os.Stat(fn)
	if err != nil {
		t.Error("Could not stat reference file", err)
	}
	f, err := os.Open(fn)
	if err != nil {
		t.Error("Could not open reference file", err)
	}
	buff := new(bytes.Buffer)
	psg.Compress(f, buff)
	got := buff.Bytes()
	if int64(len(got)) > stats.Size() {
		t.Fatalf("output was not compressed, expected < \n%d, got:\n%d", stats.Size(), len(got))
	}
	want := 227 // the best I've got so far - target is 223 or less...
	if len(got) > want {
		t.Fatalf("compression worse than expected, wanted < \n%d, got:\n%d", want, len(got))
	}
}
