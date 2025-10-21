/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A Tabula Recta implementation for simple ciphers. This Tabula Recta
 * folds the alphabet to uppercase but supports mixed input of
 * upper and lowercase characters. It supports Unicode and thus not
 * just plain ASCII making it suitable for foreign alphabets.
 * Status: Works
 *-----------------------------------------------------------------*/
package ciphers

import (
	"fmt"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/cmn"
	"slices"
	"strings"
	"unicode/utf8"
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

type IGTabulaRecta[E rune | byte] interface {
	fmt.Stringer
	cmn.IRuneLocalizer
	GetName() string
	HasRune(r E) (bool, int)
	EncodeRune(r, key E) E
	DecodeRune(r, key E) E
	EncodeRuneRaw(E, int, int) E
	DecodeRuneRaw(E, int) E
	IsCaseInsensitive() bool
	TransposeKey(k any) (int, E)
}

var _ IGTabulaRecta[byte] = (*BinaryTabulaRecta)(nil)
var _ cmn.IRuneLocalizer = (*BinaryTabulaRecta)(nil)

var ErrNotSimpleRune error = fmt.Errorf("rune is not single-byte rune")

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type BinaryTabulaRecta struct {
	Name   string
	tabula [][]byte
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/**
 * (Ctor) a Binary Tabula Recta uses a plain Caesar cipher to encode
 * or decode BINARY data. Its source alphabet is BINARY_DISK which is
 * nothing but a Binary Alphabet encompassing ASCII values 00...255
 */
func NewBinaryTabulaRecta() *BinaryTabulaRecta {
	btr := &BinaryTabulaRecta{
		Name: "Binary",
	}
	btr.tabula = btr.generateTabulaRecta(0xFF) // 256x256 (+1 is added for [0])
	return btr
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer by rendering the Tabula Recta
func (t *BinaryTabulaRecta) String() string {
	return t.renderTabulaRecta(false)
}

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *					G e n e r a l   P u r p o s e
 *- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -*/

func (t *BinaryTabulaRecta) GetName() string {
	return t.Name
}

func (t *BinaryTabulaRecta) HasRune(r byte) (bool, int) {
	var exists bool = false
	var where int = -1

	if slices.Contains(t.tabula[0], r) {
		exists = true
		where = slices.Index(t.tabula[0], r)
	}

	return exists, where
}

/**
 * (IRuneLocalizer) Find a rune in the object's alphabet catalog.
 * Rune not found: error set, other return values nil or -1.
 * Rune found: error nil, pointer to alphabet and position within.
 */
func (t *BinaryTabulaRecta) FindRune(r rune) (string, int, error) {
	var alphaStr string = ""
	var at int
	var err error = nil

	rl := utf8.RuneLen(r)
	if rl <= 2 {
		var exists bool
		exists, at = t.rowContains(0, byte(r))
		if !exists || at == -1 {
			err = fmt.Errorf("info: '%c' absent in %s", r, t.Name)
		} else {
			alphaStr = string(t.tabula[0])
		}
	} else {
		err = ErrNotSimpleRune
	}

	return alphaStr, at, err
}

func (t *BinaryTabulaRecta) TransposeKey(k any) (int, byte) {
	var transposed int = -1
	var resultingKey byte = 0

	switch cv := k.(type) {
	case rune: // called within Master
		rl := utf8.RuneLen(cv)
		if rl > 2 {
			msg := "Cannot transpose key that is not (extended)ASCII"
			mlog.ErrorT(msg, mlog.Rune("Key", cv))
		} else if exists, at := t.rowContains(0, byte(cv)); exists {
			transposed = at
			resultingKey = byte(cv)
		} else {
			msg := "Cannot transpose key we don't have"
			mlog.ErrorT(msg, mlog.Rune("Key", cv))
		}

	case byte:
		if exists, at := t.rowContains(0, cv); exists {
			transposed = at
			resultingKey = cv
		} else {
			msg := "Cannot transpose key we don't have"
			mlog.ErrorT(msg, mlog.Byte("Key", cv))
		}

	case int: // called within Slave
		max := len(t.tabula[0])
		transposed = cv % max // limit it to the length of OUR alphabet
		resultingKey = t.tabula[0][transposed]

	default:
		mlog.Error("Panicking with invalid type to transpose")
		panic("not a valid transpose parameter")
	}

	return transposed, resultingKey
}

func (t *BinaryTabulaRecta) IsCaseInsensitive() bool {
	return false
}

// Prints a BinaryTabulaRecta which is a square matrix of all
// 0xFF values for all 0xFF keys.
func (t *BinaryTabulaRecta) PrintTabulaRecta(center bool) {
	fmt.Println(t.renderTabulaRecta(true))
}

func (t *BinaryTabulaRecta) PrintTape(key byte) {
	if ok, pos := t.rowContains(0, key); ok {
		fmt.Println(t.renderTape(pos, true))
	}
}

// Gets the substitution tabula (0..255) for the given shift key value.
func (t *BinaryTabulaRecta) GetTabulaForKey(shift byte) []byte {
	return t.tabula[shift]
}

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *					E n c r y p t i o n
 *- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -*/

/**
 * Encode a rune with the key using the current Tabula Recta.
 * @param r (rune) character to encode
 * @param key (rune) encoding key
 * @returns (rune) encoded rune, or r if not found.
 */
func (t *BinaryTabulaRecta) EncodeRune(r, key byte) byte {
	var result byte = r // pass-through if not found

	if exists, atColumn := t.rowContains(0, r); exists {
		keyIndex := byte(key)
		result = t.tabula[keyIndex][atColumn]
	}

	return result
}

func (t *BinaryTabulaRecta) EncodeRuneRaw(r byte, rowIdx, colIdx int) byte {
	var result byte = r

	// Pre
	if rowIdx >= len(t.tabula) {
		mlog.Error("out-of-range row", mlog.String("At", "DecodeRuneRaw"), mlog.Int("Value", rowIdx))
		panic("Bad thing happened")
	}
	if colIdx >= len(t.tabula[0]) {
		mlog.ErrorT("out-of-range column", mlog.String("At", "DecodeRuneRaw"), mlog.Int("Value", colIdx))
		panic("Bad thing happened")
	}

	// Middle
	result = t.tabula[rowIdx][colIdx]

	return result
}

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *					D e c r y p t i o n
 *- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -*/

// Decode a rune using the key
func (t *BinaryTabulaRecta) DecodeRune(r, key byte) byte {
	var result byte = r

	if exists, keyIndex := t.HasRune(key); exists {
		exists, column := t.rowContains(keyIndex, r)
		if exists {
			result = t.tabula[0][column]
		}
	} else {
		mlog.WarnT("Key absent in alphabet", mlog.String("Alpha", t.Name), mlog.Byte("Byte", key))
	}

	return result
}

func (t *BinaryTabulaRecta) DecodeRuneRaw(r byte, rowIdx int) byte {
	var result byte = r

	if rowIdx >= len(t.tabula[0]) { // the Tabula Recta is a square matrix NxN
		mlog.ErrorT("out-of-range column", mlog.String("At", "DecodeRuneRaw"), mlog.Int("Value", rowIdx))
		panic("Bad thing happened")
	}

	// Middle
	if exists, colIdx := t.rowContains(rowIdx, r); exists {
		result = t.tabula[0][colIdx]
	} else {
		// This would never happen UNLESS someone edits TabulaCaesar(Command) and didn't
		// check for rune's presence in the slave alphabet. But this check is here
		// as a safeguard.
		mlog.ErrorT("internal error: shouldn't be looking for that rune in this slave", mlog.Byte("Byte", r), mlog.String("Alpha", t.Name))
		panic("Internal Error")
	}

	return result
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

// generates a square matrix of 256x256 with values 0..255
func (t *BinaryTabulaRecta) generateTabulaRecta(size int) [][]byte {
	return cmn.MakeSquareNumericTabula[byte](size + 1) // 0xFF becomes 256x256
}

func (t *BinaryTabulaRecta) rowContains(rowNum int, target byte) (bool, int) {
	exists := false
	column := -1

	for pos, r := range t.tabula[rowNum] {
		if r == target {
			exists = true
			column = pos
			break
		}
	}

	return exists, column
}

func (t *BinaryTabulaRecta) renderTabulaRecta(boxDrawing bool) string {
	var sb strings.Builder
	// Prints a Row of Runes
	rowPrinterFunc := func(row []byte) {
		for _, char := range row {
			sb.WriteString(fmt.Sprintf("%02x ", char))
		}

		sb.WriteRune('\n')
	}

	// Print Heading
	var bC, bH, bV rune
	if boxDrawing {
		bC = '\u250c' // ┌
		bH = '\u2500' // ─
		bV = '\u2502' // │
	} else {
		bC = '+'
		bH = '-'
		bV = '|'
	}

	title := fmt.Sprintf("%c %s %c\n", 0x00ab, t.Name, 0x00bb)
	sb.WriteString(title)

	sb.WriteString("   ")
	rowPrinterFunc(t.tabula[0])
	sb.WriteString(fmt.Sprintf("  %c%s\n", bC, strings.Repeat(string(bH), 3*len(t.tabula[0])-1)))

	for _, row := range t.tabula {
		sb.WriteString(fmt.Sprintf("%02x%c", row[0], bV))
		rowPrinterFunc(row)
	}

	return sb.String()
}

func (t *BinaryTabulaRecta) renderTape(keyShift int, boxDrawing bool) string {
	var sb strings.Builder
	// Prints a Row of Runes
	rowPrinterFunc := func(row []byte) {
		for _, char := range row {
			sb.WriteString(fmt.Sprintf("%02x ", char))
		}

		sb.WriteRune('\n')
	}

	// Print Heading
	var bC, bH, bV rune
	if boxDrawing {
		bC = '\u250c' // ┌
		bH = '\u2500' // ─
		bV = '\u2502' // │
	} else {
		bC = '+'
		bH = '-'
		bV = '|'
	}

	title := fmt.Sprintf("%c %s %c\n", 0x00ab, t.Name, 0x00bb)
	sb.WriteString(title)

	sb.WriteString("   ")
	rowPrinterFunc(t.tabula[0])
	sb.WriteString(fmt.Sprintf("  %c%s\n", bC, strings.Repeat(string(bH), 3*len(t.tabula[0])-1)))

	sb.WriteString(fmt.Sprintf("%02x%c", t.tabula[keyShift][0], bV))
	rowPrinterFunc(t.tabula[keyShift])

	return sb.String()
}

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/

func DemoBinaryTabulaRecta() {
	trB := NewBinaryTabulaRecta()

	r := cmn.BINARY_DISK.GetRuneAt(2)
	trB.PrintTape(byte(r))
	trB.PrintTabulaRecta(false)
}
