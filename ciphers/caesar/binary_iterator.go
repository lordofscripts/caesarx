package caesar

import (
	"fmt"
	"lordofscripts/caesarx"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/ciphers"
	"lordofscripts/caesarx/internal/crypto"
	"strings"
)

// BinaryIterator is used for streaming through a binary file.
type BinaryIterator struct {
	tabulas   []ciphers.IGTabulaRecta[byte]
	sequencer crypto.IKeySequencer
	sb        strings.Builder
	pos       int
	dataPtr   *([]byte)
	max       int
}

type byteRune struct {
	Rune  byte
	Shift int
}

/**
 * (ctor) regular constructor when the rune has a defined position
 */
func newByteRune(r byte, at int) *byteRune {
	return &byteRune{Rune: r, Shift: at}
}

/**
 * (ctor) extraordinary constructor when the rune isn't present
 */
func newByteRuneNotFound(r byte) *byteRune {
	return &byteRune{Rune: r, Shift: -1}
}

// NewBinaryIterator creates a new instance of a BinaryIterator.
func NewBinaryIterator(sx crypto.IKeySequencer, tabs ...ciphers.IGTabulaRecta[byte]) *BinaryIterator {
	tabulas := make([]ciphers.IGTabulaRecta[byte], 0)
	for _, tab := range tabs {
		if !isNil(tab) {
			tabulas = append(tabulas, tab)
		}
	}

	return &BinaryIterator{
		tabulas:   tabulas,
		sequencer: sx,
		pos:       -1,
		dataPtr:   nil,
		max:       -1,
	}
}

// implements fmt.Stringer
func (a *byteRune) String() string {
	return fmt.Sprintf("%02xh - %d", a.Rune, a.Shift)
}

// Start initializes the binary iterator's buffer prior to the
// EncodeNext or DecodeNext operations.
func (t *BinaryIterator) Start(buf []byte) {
	t.dataPtr = &buf
	t.pos = 0
	t.max = len(buf)
}

// getByteAt returns the said byte from the binary iterator's buffer
// and an indication whether the said return byte is valid or not.
func (t *BinaryIterator) getByteAt(buf []byte, nr int) (byte, bool) {
	if nr < 0 {
		return 0, false // @audit I guess this should panic
	}

	count := 0
	for _, r := range buf {
		if count == nr {
			return r, true
		}
		count++
	}
	return 0, false
}

// EncodeNext encodes the next byte in the binary iterator's buffer.
func (t *BinaryIterator) EncodeNext() bool {
	var result byte

	if t.pos >= t.max {
		return true
	}

	// Tabula index, Rune to en/decode, Current key, Shift within alphabet
	// currTabula, targetChar, currKey, targetShift
	if currTab, target, currKey := t.next(); currTab != -1 {
		if err := t.sequencer.Feedback(rune(target.Rune)); err != nil {
			mlog.FatalT(caesarx.ERR_SEQUENCER,
				"feedback panic",
				mlog.Int("Pos", t.pos),
				mlog.At(),
			)
		}

		if currTab == 0 { // Primary alphabet
			result = t.tabulas[0].EncodeRune(target.Rune, currKey.Rune)
		} else {
			_, keyNew := t.tabulas[currTab].TransposeKey(currKey.Shift)
			result = t.tabulas[currTab].EncodeRune(target.Rune, keyNew)
		}

		mlog.PrintCatheter("Encode",
			mlog.Int("Pos", t.pos-1),
			mlog.Rune("Rune", rune(target.Rune)),
			mlog.Rune("WithKey", rune(currKey.Rune)),
			mlog.Rune("Ciphered", rune(result)))

		t.sb.WriteByte(result)
	} else {
		mlog.PrintCatheter("Skipping", mlog.Int("Pos", t.pos-1), mlog.Rune("Rune", rune(target.Rune)))
		t.sequencer.Skip()
		t.sb.WriteByte(target.Rune)
	}

	return t.pos == t.max
}

// DecodeNext decodes the next byte in the binary iterator's buffer.
func (t *BinaryIterator) DecodeNext() bool {
	var result byte

	if t.pos >= t.max {
		return true
	}

	// Tabula index, Rune to en/decode, Current key, Shift within alphabet
	if currTab, target, keyCurr := t.next(); currTab != -1 {
		if currTab == 0 { // Primary alphabet
			result = t.tabulas[0].DecodeRune(target.Rune, keyCurr.Rune)
		} else {
			_, newKey := t.tabulas[currTab].TransposeKey(keyCurr.Shift)
			result = t.tabulas[currTab].DecodeRune(target.Rune, newKey)
		}

		mlog.PrintCatheter("Decode",
			mlog.Int("Pos", t.pos-1),
			mlog.Rune("Rune", rune(target.Rune)),
			mlog.Rune("WithKey", rune(keyCurr.Rune)),
			mlog.Rune("Plain", rune(result)))

		if err := t.sequencer.Feedback(rune(result)); err != nil {
			mlog.FatalT(caesarx.ERR_SEQUENCER,
				"feedback panic",
				mlog.Int("Pos", t.pos),
				mlog.At(),
			)
		}
		t.sb.WriteByte(result)
	} else {
		mlog.PrintCatheter("Skipping", mlog.Int("Pos", t.pos-1), mlog.Rune("Rune", rune(target.Rune)))
		t.sequencer.Skip()
		t.sb.WriteByte(target.Rune)
	}

	return t.pos == t.max
}

// moves on to the next byte in the buffer and perform the transform.
func (t *BinaryIterator) next() (int, *byteRune, *byteRune) {
	var currTabula int = -1 // Tabula index
	var target *byteRune = nil
	var currKey *byteRune = nil

	// Internal utility function to locate a character in a Tabula.
	// Can be used for both targetChar & key. It also returns the
	// ID (index) of the TabulaRecta where it was found.
	locatorFx := func(rx byte) (int, *byteRune) {
		var result *byteRune = nil
		var id int = -1

		for idxTab, tr := range t.tabulas {
			if !isNil(tr) {
				if ok, shift := tr.HasRune(rx); ok {
					id = idxTab
					result = newByteRune(rx, shift)
					break
				}
			}
		}

		return id, result
	}

	if targetChar, ok := t.getByteAt(*t.dataPtr, t.pos); ok {
		key := t.sequencer.GetKey(t.pos, rune(targetChar)) // always from Primary alphabet
		keyTabulaId, keyInfo := locatorFx(byte(key))
		if keyTabulaId != -1 {
			currKey = keyInfo
		} else {
			mlog.FatalT(70, "Key's byte not found in tabulae",
				mlog.String("At", "BinaryIterator.next()"),
				mlog.Byte("Key", byte(key)),
				mlog.Int("Pos", t.pos))
		}

		currTabula, target = locatorFx(targetChar)
		if target == nil {
			target = newByteRuneNotFound(targetChar)
		}
	}
	t.pos++

	return currTabula, target, currKey
}

// Result returns the iterator's final output as a slice of bytes.
func (t *BinaryIterator) Result() []byte {
	output := []byte(t.sb.String())
	t.sb.Reset()
	return output
}
