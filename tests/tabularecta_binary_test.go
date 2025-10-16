/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Tests for BinaryTabulaRecta (encrypt binary files instead of text)
 *-----------------------------------------------------------------*/
package tests

import (
	"lordofscripts/caesarx/ciphers"
	"testing"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *					T e s t s :: BinaryTabulaRecta
 *-----------------------------------------------------------------*/

// The IRuneFinder
func Test_Binary_FindRune(t *testing.T) {
	btr := ciphers.NewBinaryTabulaRecta()
	if btr.GetName() != "Binary" {
		t.Errorf("expected 'binary' as name")
	}

	for i := range 255 {
		if _, pos, err := btr.FindRune(rune(i)); err != nil {
			t.Errorf("at %d got %v", i, err)
		} else if pos != i {
			t.Errorf("the ASCII code %02x appeared as %02x", i, pos)
		}
	}
}

func Test_Binary_HasRune(t *testing.T) {
	btr := ciphers.NewBinaryTabulaRecta()

	for i := range 255 {
		if exists, pos := btr.HasRune(byte(i)); !exists {
			t.Errorf("at %d got that it doesn't exist", i)
		} else if pos != i {
			t.Errorf("the ASCII code %02x appeared as %02x", i, pos)
		}
	}
}

func Test_Binary_PrintTape(t *testing.T) {
	t.Skip()
	t.Helper()
	const KEY byte = 128
	btr := ciphers.NewBinaryTabulaRecta()
	btr.PrintTape(KEY)
}

func Test_Binary_PrintTabula(t *testing.T) {
	t.Skip()
	t.Helper()
	btr := ciphers.NewBinaryTabulaRecta()
	btr.PrintTabulaRecta(false)
}

func Test_Binary_EncodeRune(t *testing.T) {
	stdASCII := make([]byte, 128) // single-byte runes
	extASCII := make([]byte, 128) // 2-byte runes

	for v := range 128 {
		stdASCII[v] = byte(v)       // 0..127
		extASCII[v] = byte(128 + v) // 128..255
	}

	const KEY byte = 128
	btr := ciphers.NewBinaryTabulaRecta()
	for index, v := range stdASCII {
		if btr.EncodeRune(v, KEY) != extASCII[index] {
			t.Errorf("at %d Std %c (%d) expected Ext %c (%d)", index,
				v, v,
				extASCII[index], extASCII[index])
		}
	}
}

func Test_Binary_EncodeRune_Many(t *testing.T) {
	allCases := []struct {
		In  byte
		Key byte
		Out byte
	}{
		{0x11, 0x00, 0x11}, // unchanged with Key shift=0
		{0x12, 0x0b, 0x1d},
		{0x46, 0x53, 0x99},
		{0xff, 0xff, 0xfe},
		{0x00, 0xff, 0xff},
	}

	// one tabula recta is good for all input values
	btr := ciphers.NewBinaryTabulaRecta()
	for i, tc := range allCases {
		if v := btr.EncodeRune(tc.In, tc.Key); v != tc.Out {
			t.Errorf("#%d In:%02x Key:%02x Exp:%02x Got:%02x", i+1, tc.In, tc.Key, tc.Out, v)
		}
	}
}
