package drum

import (
	"fmt"
	"path"
	"io/ioutil"
	"os"
	"testing"
)

func TestDecodeFile(t *testing.T) {
	tData := []struct {
		path   string
		output string
	}{
		{"pattern_1.splice",
			`Saved with HW Version: 0.808-alpha
Tempo: 120
(0) kick	|x---|x---|x---|x---|
(1) snare	|----|x---|----|x---|
(2) clap	|----|x-x-|----|----|
(3) hh-open	|--x-|--x-|x-x-|--x-|
(4) hh-close	|x---|x---|----|x--x|
(5) cowbell	|----|----|--x-|----|
`,
		},
		{"pattern_2.splice",
			`Saved with HW Version: 0.808-alpha
Tempo: 98.4
(0) kick	|x---|----|x---|----|
(1) snare	|----|x---|----|x---|
(3) hh-open	|--x-|--x-|x-x-|--x-|
(5) cowbell	|----|----|x---|----|
`,
		},
		{"pattern_3.splice",
			`Saved with HW Version: 0.808-alpha
Tempo: 118
(40) kick	|x---|----|x---|----|
(1) clap	|----|x---|----|x---|
(3) hh-open	|--x-|--x-|x-x-|--x-|
(5) low-tom	|----|---x|----|----|
(12) mid-tom	|----|----|x---|----|
(9) hi-tom	|----|----|-x--|----|
`,
		},
		{"pattern_4.splice",
			`Saved with HW Version: 0.909
Tempo: 240
(0) SubKick	|----|----|----|----|
(1) Kick	|x---|----|x---|----|
(99) Maracas	|x-x-|x-x-|x-x-|x-x-|
(255) Low Conga	|----|x---|----|x---|
`,
		},
		{"pattern_5.splice",
			`Saved with HW Version: 0.708-alpha
Tempo: 999
(1) Kick	|x---|----|x---|----|
(2) HiHat	|x-x-|x-x-|x-x-|x-x-|
`,
		},
	}

	for _, exp := range tData {
		decoded, err := DecodeFile(path.Join("fixtures", exp.path))
		if err != nil {
			t.Fatalf("something went wrong decoding %s - %v", exp.path, err)
		}
		if fmt.Sprint(decoded) != exp.output {
			t.Logf("decoded:\n%#v\n", fmt.Sprint(decoded))
			t.Logf("expected:\n%#v\n", exp.output)
			t.Fatalf("%s wasn't decoded as expect.\nGot:\n%s\nExpected:\n%s",
				exp.path, decoded, exp.output)
		}
	}
}

func TestFileNotFound(t *testing.T) {
	_, err := DecodeFile("sillyfilename.splice")
	if err == nil {
		t.Log("Expecting an error when file not found, but err is nil.")
		t.Fail()
	}
}

func TestEmptyFile(t *testing.T) {
	name := path.Join("fixtures", "nobytes.splice")
	file, err := os.Create(name)
	if err != nil {
		t.Fatalf("Could not create %v: %v", name, err)
	}

	if err = file.Close(); err != nil {
		t.Fatalf("Could not close %v: %v", name, err)
	}

	if _, err = DecodeFile(name); err != ErrEmptySpliceFile {
		t.Logf("Expecting an error when file is empty, but err is: %v", err)
		t.Fail()
	}
}

func TestEmptySPLICEFile(t *testing.T) {
	// Create an empty SPLICE file (One that contains only "SPLICE" at start)
	name := path.Join("fixtures", "empty_SPLICE.splice")
	err := ioutil.WriteFile(name, []byte("SPLICE"), 0644)
	if err != nil {
		t.Fatalf("Could not write to file %v: %v", name, err)
	}

	_, err = DecodeFile(name)
	if err == nil {
		t.Log("Expecting an error when SPLICE file is empty, but err is nil.")
		t.Fail()
	}
}
