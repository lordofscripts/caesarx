/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Tests for FileHeader
 *-----------------------------------------------------------------*/
package tests

import (
	"fmt"
	"lordofscripts/caesarx"
	"lordofscripts/caesarx/internal/files"
	"os"
	"testing"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *					T e s t s :: AffineHelper
 *-----------------------------------------------------------------*/

// File Header test. It writes a valid header to a binary file, then
// reads it back and checks they are the same and valid.
func Test_FileHeader_ReadWrite(t *testing.T) {
	const OUTPUT_BIN_FILE = "testdata/header.bin"

	// I. create the (temporary) target binary file
	fileOut, err := os.OpenFile(OUTPUT_BIN_FILE, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Errorf("could not open binary file for writing: %v", err)
	}

	// II. write the header to the file, this should be successful
	fhOut, err := files.NewFileHeader(caesarx.DidimusCipher, OUTPUT_BIN_FILE)
	if err != nil {
		t.Error(err)
	} else {
		err = fhOut.Write(fileOut)
		if err != nil {
			t.Errorf("FileHeader.Write() failed: %v", err)
		}
	}
	fileOut.Close()

	// III. read the header from the binary file
	fhIn := files.NewEmptyFileHeader()
	// 3.a) the file stream reader
	fileIn, err := os.Open(OUTPUT_BIN_FILE)
	if err != nil {
		t.Error(err)
	} else {
		// 3.b) read the header
		err = fhIn.Read(fileIn)
		if err != nil {
			t.Errorf("FileHeader.Read() failed: %v", err)
		}
	}
	fileIn.Close()

	// IV. Compare written vs. read
	if !fhIn.Equals(fhOut) {
		fmt.Println("Written", fhIn)
		fmt.Println("Read", fhOut)
		t.Errorf("file header read vs. write do not match")
	}

	os.Remove(OUTPUT_BIN_FILE)
}
