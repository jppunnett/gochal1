package drum

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
)

var (
	// ErrEmptySpliceFile occurs when opening a splice files that has no bytes.
	ErrEmptySpliceFile = errors.New("Not a SPLICE file")

	// ErrBadFileType indicates that the first few bytes of a file does not contain
	// a splice file ID i.e. "SPLICE". 
	ErrBadFileType = errors.New("Not a SPLICE file")

	// ErrNoRemBytesFld indicates that a SPLICE file does not contain a field
	// representing the number of bytes remaining.
	ErrNoRemBytesFld = errors.New("Missing remaing-bytes field.")

	// ErrInvalidNumBytes indicates the field containing the number of bytes
	// remaining is not correct i.e. there should be at least that number of 
	// bytes remaining in the file.
	ErrInvalidNumBytes = errors.New("Bytes remaining is zero")
)

const (
	szFileIDFld = 6
	szPlatFld   = 11

	posRemBytesFld = 13
	posStartOfData = 14

	posStartTempo = 32
	szTempoFld    = 4

	posFirstInstrument = 36

	numSteps = 16
)

// track represents a single track in a Pattern. A Pattern may have many tracks. 
type track struct {
	id    uint
	name  string
	steps [numSteps]byte
}

// Pattern is the high level representation of the
// drum pattern contained in a .splice file.
type Pattern struct {
	hwver    string

	tempo  float32
	tracks []track
}

func (p *Pattern) String() string {
	return fmt.Sprintf("Saved with HW Version: %v\nTempo: %v\n%v",
		p.hwver,
		p.tempo,
		p.tracksAsString(),
	)
}

func (p *Pattern) tracksAsString() string {
	format := "(%d) %s\t|%s\n"
	trkstr := ""

	for _, t := range p.tracks {
		var stepstr string
		for j, s := range t.steps {
			if s == 1 {
				stepstr += "x"
			} else {
				stepstr += "-"
			}

			if (j+1) % 4 == 0 {
				stepstr += "|"	
			}
		}

		trkstr += fmt.Sprintf(format, t.id, t.name, stepstr)
	}
	
	return trkstr
}

func isValidateSplice(contents []byte) (error) {

	nbytes := len(contents)
	// Check that we have at least *something* to decode.
	if nbytes == 0 {
		return ErrEmptySpliceFile
	}

	// Check if this is a splice file.
	filetype := string(contents[:szFileIDFld])
	if filetype != "SPLICE" {
		return ErrBadFileType
	}

	// Check that we have at least the number-of-bytes-remaining field.
	if nbytes < posRemBytesFld + 1 {
		return ErrNoRemBytesFld
	}

	return nil
}

// DecodeFile decodes the drum machine file found at the provided path
// and returns a pointer to a parsed pattern which is the entry point to the
// rest of the data.
func DecodeFile(path string) (*Pattern, error) {
	
	//	The structure of the file appears to be:
	//	File type/id e.g. "SPLICE"
	//	Number of bytes remaining e.g. 197
	//	Hardware version e.g. "0.808-alpha"
	// 	The tempo e.g. 120 or 98.4 so must be a float
	//	The remaining bytes define the tracks and are laid out like this:
	//		Track ID
	//		Length of instrument name 
	//		Instrument name
	//		Steps
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err = isValidateSplice(contents); err != nil {
		return nil, err
	}

	// Check the number of bytes left to decode is *at least* what is reported in
	// the remaining-bytes field.
	reptdBytesRem := uint(contents[posRemBytesFld])
	actualBytesRem := len(contents) - posRemBytesFld - 1
	if uint(actualBytesRem) < reptdBytesRem {
		return nil, ErrInvalidNumBytes
	}

	//	Decode the remaining contents of the splice file.
	remaining := contents[posStartOfData : posStartOfData + reptdBytesRem]

	p := &Pattern{}
	
	// The hardware version. remaining[:szPlatFld] could end up with trailing 
	// zeros so we trim before converting to a string.
	hwver := bytes.TrimRight(remaining[:szPlatFld], string([]byte{0}))
	p.hwver = string(hwver)

	// The tempo is a float so we cannot simply cast the byte array to a float.
	buf := bytes.NewReader(remaining[posStartTempo : posStartTempo + szTempoFld])
	if err = binary.Read(buf, binary.LittleEndian, &p.tempo); err != nil {
		return nil, err
	}

	// Decode each track
	for byteptr := posFirstInstrument; byteptr < len(remaining); {
		t := track{}

		//	The instrument Id
		t.id = uint(remaining[byteptr])
		byteptr += 4

		//	szname is the size (in bytes) of the instrument name that follows
		szname := int(remaining[byteptr])
		byteptr++

		//	The instrument name
		t.name = string(remaining[byteptr : byteptr+szname])
		byteptr += szname

		//	The steps
		for i := 0; i < numSteps; i++ {
			t.steps[i] = remaining[byteptr+i]
		}

		byteptr += numSteps

		p.tracks = append(p.tracks, t)
	}

	return p, err
}
