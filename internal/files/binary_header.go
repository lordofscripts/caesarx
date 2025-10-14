/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Magic Headers for binary encrypted files
 *-----------------------------------------------------------------*/
package files

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"lordofscripts/caesarx"
	"lordofscripts/caesarx/app/mlog"
	"path"
	"strings"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	FILEHEADER_MAJOR uint8  = 0x01
	FILEHEADER_MINOR uint8  = 0x00
	FILEHEADER_START uint32 = 0xBABEF007
	FILEHEADER_END   uint16 = 0xDEAD
	CAE              uint16 = 0xCAE5
	DID              uint16 = 0xD1D1
	FIB              uint16 = 0xF1B0
	BEL              uint16 = 0xBE50
	VIG              uint16 = 0xB16E
	AFI              uint16 = 0xAF1E
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// (internal) FileHeaderStart will be the first bytes of the encrypted binary file
// for our internal identification. This is a fixed-size structure.
type FileHeaderStart struct {
	Magic        uint32 // a magic number to identify our files
	MajorVersion uint8
	MinorVersion uint8
	AlgorithmA   caesarx.CipherVariant
	AlgorithmB   uint16
}

// (internal) file header ending marker. This is a variable-size structure.
type FileHeaderEnd struct {
	ExtLen    byte
	Extension string
	Trailer   uint16
}

// A binary file header that contains basic information about the
// encrypted file that may be useful when recovering it (decode).
type FileHeader struct {
	Start   *FileHeaderStart
	End     *FileHeaderEnd
	isValid bool
}

/* ----------------------------------------------------------------
 *						C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// returns a file header for WRITING an encrypted binary file. The header
// contains information about the encryption algorithm and the
// original filename's extension.
// This file header is an optional feature. Sometimes you may wish
// to sacrifice something to ensure something can be decrypted
// after you pass away.
func NewFileHeader(cipherId caesarx.CipherVariant, filename string) (*FileHeader, error) {
	if fhE, err := newFileHeaderEnd(filename); err != nil {
		return nil, err
	} else {
		fhS, err := newFileHeaderStart(cipherId)
		fh := &FileHeader{
			Start:   fhS,
			End:     fhE,
			isValid: true,
		}
		return fh, err
	}
}

// returns an empty instance that can be used for reading a file header
// from a binary file stream
func NewEmptyFileHeader() *FileHeader {
	return &FileHeader{
		Start:   new(FileHeaderStart),
		End:     new(FileHeaderEnd),
		isValid: false,
	}
}

// create a pre-baked fileheader for the selected algorithm. The
// file type is given in the Write method.
func newFileHeaderStart(cipherId caesarx.CipherVariant) (*FileHeaderStart, error) {
	var id uint16
	switch cipherId {
	case caesarx.AffineCipher:
		id = AFI
	case caesarx.BellasoCipher:
		id = BEL
	case caesarx.CaesarCipher:
		id = CAE
	case caesarx.DidimusCipher:
		id = DID
	case caesarx.FibonacciCipher:
		id = FIB
	case caesarx.VigenereCipher:
		id = VIG
	default:
		return nil, fmt.Errorf("invalid algorithm ID by file header")
	}

	return &FileHeaderStart{
		Magic:        FILEHEADER_START,
		MajorVersion: FILEHEADER_MAJOR,
		MinorVersion: FILEHEADER_MINOR,
		AlgorithmA:   cipherId,
		AlgorithmB:   id,
	}, nil
}

// (internal) create the end of a file header, it contains some
// information about the file
func newFileHeaderEnd(filename string) (*FileHeaderEnd, error) {
	ext := path.Ext(filename)
	if len(ext) >= 2 && strings.HasPrefix(ext, ".") {
		ext = ext[1:]
	}
	size := len(ext)
	if size > 255 {
		const MSG = "extension too long"
		mlog.ErrorT(MSG, mlog.String("Ext", ext), mlog.Int("Size", size))
		err := fmt.Errorf("%s '%s' is %d", MSG, ext, size)
		return nil, err
	}

	return &FileHeaderEnd{
		ExtLen:    byte(size),
		Extension: ext,
		Trailer:   FILEHEADER_END,
	}, nil
}

/* ----------------------------------------------------------------
 *				M e t h o d s :: FileHeaderEnd
 *-----------------------------------------------------------------*/

// implements fmt.Stringer for FileHeaderEnd
func (fe *FileHeaderEnd) String() string {
	return fmt.Sprintf("EndFH (%d) %s %0x", fe.ExtLen, fe.Extension, fe.Trailer)
}

func (fe *FileHeaderEnd) Equals(other *FileHeaderEnd) bool {
	result := false

	if other != nil {
		if fe.ExtLen == other.ExtLen &&
			fe.Extension == other.Extension &&
			fe.Trailer == other.Trailer {
			result = true
		}
	}

	return result
}

/* ----------------------------------------------------------------
 *				M e t h o d s :: FileHeaderStart
 *-----------------------------------------------------------------*/

// implements fmt.Stringer for FileHeaderStart
func (fs *FileHeaderStart) String() string {
	return fmt.Sprintf("StartFH %0x v%x.%02x %s (%0x)", fs.Magic, fs.MajorVersion, fs.MinorVersion, fs.AlgorithmA, fs.AlgorithmB)
}

func (fs *FileHeaderStart) Equals(other *FileHeaderStart) bool {
	result := false

	if other != nil {
		result = (fs.Magic == other.Magic) &&
			(fs.MajorVersion == other.MajorVersion) &&
			(fs.MinorVersion == other.MinorVersion) &&
			(fs.AlgorithmA == other.AlgorithmA) &&
			(fs.AlgorithmB == other.AlgorithmB)
	}

	return result
}

/* ----------------------------------------------------------------
 *				M e t h o d s :: FileHeader
 *-----------------------------------------------------------------*/

// implements fmt.Stringer
func (fh *FileHeader) String() string {
	const NL = "\n"
	var sb strings.Builder
	sb.WriteString("Header:Prologue" + NL)
	sb.WriteString(fmt.Sprintf("\tVersion: %x.%02x", fh.Start.MajorVersion, fh.Start.MinorVersion) + NL)
	sb.WriteString("\tCipher: " + fh.Start.AlgorithmA.String() + NL)
	sb.WriteString(fmt.Sprintf("\tAlgo: %02x\n", fh.Start.AlgorithmB))

	sb.WriteString("Header:Epilogue" + NL)
	sb.WriteString("\tExtension: " + fh.End.Extension + NL)
	return sb.String()
}

// whether this file header is valid, else it is incomplete
func (fh *FileHeader) IsValid() bool {
	return fh.isValid
}

// returns the file extension recorded in the header BUT
// with the leading "."
func (fh *FileHeader) FileExtension() string {
	const EXT_DIVIDER string = "."
	var result string = ""

	if fh.End.ExtLen > 0 {
		if !strings.HasPrefix(fh.End.Extension, EXT_DIVIDER) {
			result = EXT_DIVIDER + fh.End.Extension
		}
	}

	return result
}

// compares the equality of two file header instances
// by their values
func (fh *FileHeader) Equals(other *FileHeader) bool {
	result := false
	if other != nil {
		result = fh.isValid && other.isValid &&
			fh.Start.Equals(other.Start) &&
			fh.End.Equals(other.End)
	}

	return result
}

// writes the current header information to a binary file
func (fh *FileHeader) Write(w io.Writer) error {
	var err error
	buf := new(bytes.Buffer)

	if err = binary.Write(buf, binary.LittleEndian, fh.Start); err == nil {
		if err = binary.Write(buf, binary.LittleEndian, fh.End.ExtLen); err == nil {
			extSlice := []byte(fh.End.Extension)
			if err = binary.Write(buf, binary.LittleEndian, extSlice); err == nil {
				err = binary.Write(buf, binary.LittleEndian, fh.End.Trailer)
			}
		}
	}

	if err == nil {
		err = binary.Write(w, binary.LittleEndian, buf.Bytes())
	}

	return err
}

// reads the entire header (fixed and variable part) from a binary
// file. If it is valid, you can read the real file extension using
// the Extension() method.
func (fh *FileHeader) Read(r io.Reader) error {
	var err error

	// read fixed header
	if err = binary.Read(r, binary.LittleEndian, fh.Start); err == nil {
		// read extension length as part of the variable-size header epilogue
		if err = binary.Read(r, binary.LittleEndian, &fh.End.ExtLen); err == nil {
			if fh.End.ExtLen == 0 {
				fh.End.Extension = ""
			} else {
				extData := make([]byte, fh.End.ExtLen)
				if err = binary.Read(r, binary.LittleEndian, extData); err == nil {
					fh.End.Extension = string(extData)
				} else {
					return err
				}
			}

			// whether ext="" or "extension" we must read the trailer
			err = binary.Read(r, binary.LittleEndian, &fh.End.Trailer)
		}
	}

	// check file header
	if err == nil {
		// check the trailer is as expected
		if fh.End.Trailer != FILEHEADER_END {
			return fmt.Errorf("invalid file header trailer %v", fh.End.Trailer)
		}
		// check the header is correct
		if fh.Start.Magic != FILEHEADER_START {
			return fmt.Errorf("invalid file header magic %v", fh.Start.Magic)
		}
		// check compatibility of major version
		if fh.Start.MajorVersion != FILEHEADER_MAJOR {
			return fmt.Errorf("incompatible file header major version, exp:%x got:%x", FILEHEADER_MAJOR, fh.Start.MajorVersion)
		}

		fh.isValid = true
	}

	return err
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

/*
func demo() {
	const DUMMY string = "filename.bin"
	file, err := os.OpenFile(DUMMY, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		mlog.Fatal(-1, "could not open binary file for writing", err)
	}

	if headerW, err := NewFileHeader(caesarx.CaesarCipher, "secret.txt"); err != nil {
		fmt.Println("Error", err)
	} else {
		// now we have both the fixed and variable part of the header
		if err := headerW.Write(file); err != nil {
			fmt.Println("Error", err)
		}
	}
	file.Close()

	if fileIn, err := os.Open(DUMMY); err != nil {
		mlog.Fatal(-1, "could not open binary file for reading", err)
	} else {
		headerR := NewEmptyFileHeader()
		if err := headerR.Read(fileIn); err == nil {
			fmt.Println("Header: ", headerR)
		} else {
			mlog.Fatal(-1, "invalid header", err)
		}
	}
	defer file.Close()
}
*/
