package psg_test

import (
	"bytes"
	"testing"

	"github.com/retcon85/retcon-util-audio/model/psg"
)

func TestSearchForBlockSimple(t *testing.T) {
	s := bytes.NewBuffer([]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', '5', '6', '7', '8', '9', 'a'})
	buff := new(bytes.Buffer)
	psg.Compress(s, buff)
	got := buff.Bytes()
	want := []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 0x0a, 0x04, 0x00}
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
	s := bytes.NewBuffer([]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', '5', '6', '7', '8', '9', 'a', 'b', 'c', '5', '6', '7', '8', '9', 'a'})
	buff := new(bytes.Buffer)
	psg.Compress(s, buff)
	got := buff.Bytes()
	want := []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 0x0a, 0x04, 0x00, 'b', 'c', 0x0a, 0x04, 0x00}
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
	s := bytes.NewBuffer([]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', '2', '3', '4', '5', '6', 'e', '2', '3', '4', '5', '6', 'f', '5', '6', '7', '8', '9', 'g'})
	buff := new(bytes.Buffer)
	psg.Compress(s, buff)
	got := buff.Bytes()
	want := []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 0x09, 0x01, 0x00, 'e', 0x09, 0x01, 0x00, 'f', 0x09, 0x04, 0x00, 'g'}
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
	s := bytes.NewBuffer([]byte{'1', '2', '3', '4', '5', '1', '2', '3', '4', '1', '2', '3', '4', '5', '1', '2', '3', '4'})
	buff := new(bytes.Buffer)
	psg.Compress(s, buff)
	got := buff.Bytes()
	want := []byte{'1', '2', '3', '4', '5', '1', '2', '3', '4', 0x0d, 0x00, 0x00}
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
	s := bytes.NewBuffer([]byte{'8', '8', '8', '8', '1', '2', '3', '4', '8', '8', '8', '8', '8', '8', '8', '8', '5', '6', '8', '8', '8', '8', '7', '8', '9', '8', '8', '8', '8', '8', '8', '8', '8', '7', '6', '5', '4'})
	buff := new(bytes.Buffer)
	psg.Compress(s, buff)
	got := buff.Bytes()
	want := []byte{'8', '8', '8', '8', '1', '2', '3', '4', 0x08, 0x00, 0x00, 0x08, 0x00, 0x00, '5', '6', 0x08, 0x00, 0x00, '7', '8', '9', 0x08, 0x00, 0x00, 0x08, 0x00, 0x00, '7', '6', '5', '4'}
	if len(want) != len(got) {
		t.Fatalf("result wrong length, expected:\n%d, got:\n%d", want, got)
	}
	for i, wanted := range want {
		if got[i] != wanted {
			t.Fatalf("failed at index %d, expected:\n%d, got:\n%d", i, want, got)
		}
	}
}
