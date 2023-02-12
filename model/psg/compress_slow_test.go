package psg_test

import (
	"bytes"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/retcon85/retcon-util-audio/model/psg"
)

func TestRoundTrip(t *testing.T) {
	tdpath := "../../testdata"
	des, err := os.ReadDir(tdpath)
	if err != nil {
		t.Error("could not read testdata directory")
	}
	for _, de := range des {
		test := path.Join(tdpath, de.Name())
		if de.IsDir() || !strings.HasSuffix(test, "_decompressed.psg") {
			continue
		}

		var f []byte
		f, err = os.ReadFile(test)
		if err != nil {
			t.Fatal("could not open test file", err)
		}
		buf := new(bytes.Buffer)
		psg.Compress(bytes.NewReader(f), buf)
		compressed := buf.Bytes()
		buf.Reset()
		psg.Decompress(bytes.NewReader(compressed), buf)
		got := buf.Bytes()

		if !bytes.Equal(got, f) {
			t.Errorf("round-trip of '%s' did not match original file bytes", path.Base(test))
		}
	}
}

func TestCompressionComparesWellToPsgcomp(t *testing.T) {
	factor := 1.005 // worst result so far is 0.4% larger than psgcomp

	tdpath := "../../testdata"
	des, err := os.ReadDir(tdpath)
	if err != nil {
		t.Error("could not read testdata directory")
	}
	for _, de := range des {
		test := path.Join(tdpath, de.Name())
		if de.IsDir() || !strings.HasSuffix(test, "_decompressed.psg") {
			continue
		}

		ref := strings.Replace(test, "_decompressed.psg", "_psgcomp.psg", 1)
		fstat, err := os.Stat(ref)
		if err != nil {
			t.Fatal("could not stat reference file", err)
		}
		var f *os.File
		f, err = os.Open(test)
		if err != nil {
			t.Fatal("could not open test file", err)
		}
		bufgot := new(bytes.Buffer)
		psg.Compress(f, bufgot)
		got := bufgot.Len()
		want := int(float64(fstat.Size()) * factor)

		if got > want {
			t.Errorf("compression of '%s' was significantly worst than psgcomp - got %d, needed <= %d", path.Base(test), got, want)
		}
	}
}

func TestCompressionComparesWellToPsgtool(t *testing.T) {
	factor := 1.051 // worst result so far is 5.07% larger than psgtool

	tdpath := "../../testdata"
	des, err := os.ReadDir(tdpath)
	if err != nil {
		t.Error("could not read testdata directory")
	}
	for _, de := range des {
		test := path.Join(tdpath, de.Name())
		if de.IsDir() || !strings.HasSuffix(test, "_decompressed.psg") {
			continue
		}

		ref := strings.Replace(test, "_decompressed.psg", "_psgtool.psg", 1)
		fstat, err := os.Stat(ref)
		if err != nil {
			t.Fatal("could not stat reference file", err)
		}
		var f *os.File
		f, err = os.Open(test)
		if err != nil {
			t.Fatal("could not open test file", err)
		}
		bufgot := new(bytes.Buffer)
		psg.Compress(f, bufgot)
		got := bufgot.Len()
		want := int(float64(fstat.Size()) * factor)

		if got > want {
			t.Errorf("compression of '%s' was significantly worst than psgcomp - got %d, needed <= %d", path.Base(test), got, want)
		}
	}
}
