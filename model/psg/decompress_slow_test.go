package psg_test

import (
	"bytes"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/retcon85/retcon-util-audio/model/psg"
)

func TestSamples(t *testing.T) {
	tdpath := "../../testdata"
	des, err := os.ReadDir(tdpath)
	if err != nil {
		t.Error("could not read testdata directory")
	}
	for _, de := range des {
		test := path.Join(tdpath, de.Name())
		if de.IsDir() || !strings.HasSuffix(test, "_psgcomp.psg") {
			continue
		}

		ref := strings.Replace(test, "_psgcomp.psg", "_decompressed.psg", 1)
		bref, err := os.ReadFile(ref)
		if err != nil {
			t.Fatal("could not open reference file", err)
		}
		var f *os.File
		f, err = os.Open(test)
		if err != nil {
			t.Fatal("could not open test file", err)
		}
		bufgot := new(bytes.Buffer)
		psg.Decompress(f, bufgot)
		got := bufgot.Bytes()

		if !bytes.Equal(got, bref) {
			t.Errorf("decompression of '%s' did not match reference file '%s'", path.Base(test), path.Base(ref))
		}
	}
}
