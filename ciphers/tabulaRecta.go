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
	"strings"
	"unicode"
	"unicode/utf8"

	"lordofscripts/caesarx/cmn"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

type ITabulaRecta interface {
	fmt.Stringer
	cmn.IRuneLocalizer
	GetName() string
	HasRune(r rune) (bool, int)
	EncodeRune(r, key rune) rune
	DecodeRune(r, key rune) rune
	EncodeRuneRaw(rune, int, int) rune
	DecodeRuneRaw(rune, int) rune
	IsCaseInsensitive() bool
	TransposeKey(k any) (int, rune)
}

var _ IGTabulaRecta[rune] = (*TabulaRecta)(nil)
var _ cmn.IRuneLocalizer = (*TabulaRecta)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// A Tabula Recta is the basis of simple ciphers like Caesar,
// Bellaso and Vigenère.
// · Julius Caesar uses a single TabulaRecta because it has a key consisting//
//
//	of a single rune from the encoding alphabet. (Anno 58 B.C.)
//
// · Giovanni Battista Bellaso's cipher uses multiple TabulaRectas, one for//
//
//		 each character in the Secret. The Secret is repeated over the length of
//		 the text so that each character in the Secret becomes the key for the
//		 corresponding character position in the input text. The Bellaso cipher
//	  is thus based on Caesar but is polialphabetic. (Anno 1553).
//		 Bellasos' cipher had half the amount of keys because it used key pairs.
//
// · Blaise de Vigenère's Auto-Key cipher is based on the Bellaso cipher
type TabulaRecta struct {
	Name        string
	caseFolding bool
	alphabet    string
	tabula      [][]rune
	specialCase *cmn.SpecialCaseHandler
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/**
 * Create a new TabulaRecta. If foldCase is true, then the table is
 * converted to uppercase. If the alphabet has special casing rules,
 * they are taken in consideration as well.
 * Case-folding handling is only guaranteed when this constructor is
 * used! Therefore, do not construct your own TabulaRecta{}.
 * @param alphabet (*cmn.Alphabet) The encoding alphabet
 * @param foldCase (bool) true to preserve (upper/lower)
 * @param encodeSpace (bool) Add SPACE at end of alphabet to improve cipher.
 */
func NewTabulaRecta(alphabetO *cmn.Alphabet, foldCase bool) *TabulaRecta {
	var specialCase *cmn.SpecialCaseHandler = nil

	// Let's not alter the original
	alphabet := alphabetO.Clone()

	var letters string
	if foldCase {
		if specialCase = alphabet.BorrowSpecialCase(); specialCase != nil {
			specialCase.Assert()
			letters = specialCase.ToUpperString(alphabet.Chars)
		} else {
			letters = strings.ToUpper(alphabet.Chars)
		}
	} else {
		letters = alphabet.Chars
	}

	tr := &TabulaRecta{
		Name:        alphabet.Name,
		caseFolding: foldCase,
		alphabet:    letters,
		tabula:      nil,
		specialCase: specialCase,
	}

	tr.generateTabulaRecta()
	return tr
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// GenerateTabulaRecta creates a Tabula Recta for the given rune alphabet.
func (t *TabulaRecta) generateTabulaRecta() {
	alphabet := []rune(t.alphabet)
	size := len(alphabet)
	t.tabula = make([][]rune, size)

	for i := range size {
		row := make([]rune, size)
		for j := 0; j < size; j++ {
			var item rune = alphabet[(i+j)%size] // as-is
			/*
				Case-folding already done by the Ctor.
				if foldCase {
					item = unicode.ToUpper(item)
				}
			*/
			row[j] = item
		}
		t.tabula[i] = row
	}
}

func (t *TabulaRecta) renderTabulaRecta(center, boxDrawing bool) string {
	var sb strings.Builder
	// Prints a Row of Runes
	rowPrinterFunc := func(row []rune) {
		for _, char := range row {
			sb.WriteString(fmt.Sprintf("%c ", char))
		}

		sb.WriteRune('\n')
	}

	// Generates a Space Leader to Center a string
	const MAX_WIDTH = 80
	centerLeaderFunc := func(length int) string {
		leaderLength := int((MAX_WIDTH - length) / 2)
		return strings.Repeat(" ", leaderLength)
	}

	var leader string = ""
	if center {
		leader = centerLeaderFunc(len(t.tabula[0])*2 + 2)
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
	sb.WriteString(centerLeaderFunc(len(title)))
	sb.WriteString(title)

	sb.WriteString(leader)
	sb.WriteString("  ")
	rowPrinterFunc(t.tabula[0])
	sb.WriteString(fmt.Sprintf("%s %c%s\n", leader, bC, strings.Repeat(string(bH), 2*len(t.tabula[0])-1)))

	for _, row := range t.tabula {
		sb.WriteString(fmt.Sprintf("%s%c%c", leader, row[0], bV))
		rowPrinterFunc(row)
	}

	return sb.String()
}

func (t *TabulaRecta) renderTape(keyShift int, center, boxDrawing bool) string {
	var sb strings.Builder
	// Prints a Row of Runes
	rowPrinterFunc := func(row []rune) {
		for _, char := range row {
			sb.WriteString(fmt.Sprintf("%c ", char))
		}

		sb.WriteRune('\n')
	}

	// Generates a Space Leader to Center a string
	const MAX_WIDTH = 80
	centerLeaderFunc := func(length int) string {
		leaderLength := int((MAX_WIDTH - length) / 2)
		return strings.Repeat(" ", leaderLength)
	}

	var leader string = ""
	if center {
		leader = centerLeaderFunc(len(t.tabula[0])*2 + 2)
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
	sb.WriteString(centerLeaderFunc(len(title)))
	sb.WriteString(title)

	sb.WriteString(leader)
	sb.WriteString("  ")
	rowPrinterFunc(t.tabula[0])
	sb.WriteString(fmt.Sprintf("%s %c%s\n", leader, bC, strings.Repeat(string(bH), 2*len(t.tabula[0])-1)))

	sb.WriteString(fmt.Sprintf("%s%c%c", leader, t.tabula[keyShift][0], bV))
	rowPrinterFunc(t.tabula[keyShift])

	return sb.String()
}

func (t *TabulaRecta) rowContains(rowNum int, target rune) (bool, int) {
	rowRunes := t.tabula[rowNum]
	//log.Printf("**** Row %d %c %s", rowNum, target, string(rowRunes))
	for pos, r := range rowRunes {

		if r == target {
			//log.Printf("rowContains Pos %d/%d is %c", pos, len(rowRunes), r)
			return true, pos
		}
	}
	return false, -1
}

/**
 * Checks if the target & key runes must be case-folded either through the
 * default library functions, or the special case rules of the alphabet.
 * @param r (rune) target rune that MAY be converted
 * @param key (rune) key's rune that MAY be converted
 * @returns (rune) (converted, if case-folded) target rune
 * @returns (rune) (converted, if case-folded) key's rune
 * @returns (bool) true if conversion took place, false if returned unchanged.
 */
func (t *TabulaRecta) normalizeRuneCase(r, key rune) (rune, rune, bool) {
	// @note IMPORTANT The German ß rune is always reported as
	// lowercase by Unicode. Converting it to Uppercase returns
	// ß in GO but 'SS' in other languages! That is causing an
	// issue in this code because IsLower(ß) says TRUE when all
	// TabulaRecta values are Uppercase!
	var isConvertedCase bool = false
	if t.caseFolding { // because the T.R. is in uppercase
		if t.specialCase == nil {
			if unicode.IsLower(r) {
				r = unicode.ToUpper(r)
				isConvertedCase = true
			}
			key = unicode.ToUpper(key)
		} else {
			if t.specialCase.IsLowerRune(r) {
				r = t.specialCase.ToUpperRune(r)
				isConvertedCase = true
			}
			key = t.specialCase.ToUpperRune(key)
		}
	}

	return r, key, isConvertedCase
}

func (t *TabulaRecta) denormalizeRuneCase(r rune) rune {
	var result rune
	if t.specialCase == nil {
		result = unicode.ToLower(r)
	} else {
		result = t.specialCase.ToLowerRune(r)
	}

	return result
}

func (t *TabulaRecta) GetName() string {
	return t.Name
}

func (t *TabulaRecta) HasRune(r rune) (bool, int) {
	if t.caseFolding {
		// The tabula directory is uppercase; therefore, convert param
		if t.specialCase == nil {
			r = unicode.ToUpper(r)
		} else {
			r = t.specialCase.ToUpperRune(r)
		}
	}

	//exists, where := t.rowContains(0, r)
	var exists bool
	var where int = -1
	exists = strings.Contains(t.alphabet, string(r))
	if exists {
		where = cmn.RuneIndex(t.alphabet, r)
	}

	return exists, where
}

/**
 * Encode a rune with the key using the current Tabula Recta.
 * @param r (rune) character to encode
 * @param key (rune) encoding key
 * @returns (rune) encoded rune, or r if not found.
 */
func (t *TabulaRecta) EncodeRune(r, key rune) rune {
	var result rune = r // pass-through if not found

	// When case-folding is enabled, the TabulaRecta and the Key
	// are uppercase and we respect lower/uppercase in input text.
	var isConvertedCase bool = false
	if t.caseFolding { // because the T.R. is in uppercase
		if t.specialCase == nil {
			if unicode.IsLower(r) {
				r = unicode.ToUpper(r)
				isConvertedCase = true
			}
			key = unicode.ToUpper(key)
		} else {
			if t.specialCase.IsLowerRune(r) {
				r = t.specialCase.ToUpperRune(r)
				isConvertedCase = true
			}
			key = t.specialCase.ToUpperRune(key)
		}
	}

	//if exists, keyIndex := t.rowContains(0, key); exists {
	if exists, keyIndex := t.HasRune(key); exists {
		if exists, column := t.rowContains(0, r); exists {
			result = t.tabula[keyIndex][column]
		}
	} else {
		mlog.WarnT("Key absent in alphabet", mlog.String("Alpha", t.Name), mlog.Rune("Rune", key))
	}

	if isConvertedCase { // respect input text's upper/lowercase
		if t.specialCase == nil {
			result = unicode.ToLower(result)
		} else {
			result = t.specialCase.ToLowerRune(result)
		}
	}

	return result
}

func (t *TabulaRecta) EncodeRuneRaw(r rune, rowIdx, colIdx int) rune {
	var result rune = r

	// Pre
	if rowIdx >= len(t.tabula) {
		mlog.Error("out-of-range row", mlog.String("At", "DecodeRuneRaw"), mlog.Int("Value", rowIdx))
		panic("Bad thing happened")
	}
	if colIdx >= len(t.tabula[0]) {
		mlog.ErrorT("out-of-range column", mlog.String("At", "DecodeRuneRaw"), mlog.Int("Value", colIdx))
		panic("Bad thing happened")
	}

	isConvertedCase := t.caseFolding && ((t.specialCase == nil && unicode.IsLower(r)) ||
		(t.specialCase != nil && t.specialCase.IsLowerRune(r)))

	// Middle
	result = t.tabula[rowIdx][colIdx]

	// Post
	if isConvertedCase { // respect input text's upper/lowercase
		if t.specialCase == nil {
			result = unicode.ToLower(result)
		} else {
			result = t.specialCase.ToLowerRune(result)
		}
	}

	return result
}

/*
func (t *TabulaRecta) EncodeRuneByShift(r rune, shift int) rune {
	var result rune = r

	// check just-in-case
	max := utf8.RuneCountInString(t.alphabet)
	if shift >= max {
		slog.Error("Cannot shift-encode beyond alphabet", slog.Int("Shift", shift), slog.Int("Max", max))
		panic("Bad thing happened")
	}

	// Pre
	var isConvertedCase bool = false
	if t.caseFolding { // because the T.R. is in uppercase
		if t.specialCase == nil {
			if unicode.IsLower(r) {
				r = unicode.ToUpper(r)
				isConvertedCase = true
			}
		} else {
			if t.specialCase.IsLowerRune(r) {
				r = t.specialCase.ToUpperRune(r)
				isConvertedCase = true
			}
		}
	}

	// Middle
	if exists, column := t.rowContains(0, r); exists {
		result = t.tabula[shift][column]
	} else {
		slog.Error("reference Slave alphabet does not contain rune", slog.String("Alpha", t.Name), slog.String("Rune", string(r)))
	}

	// Post
	if isConvertedCase { // respect input text's upper/lowercase
		if t.specialCase == nil {
			result = unicode.ToLower(result)
		} else {
			result = t.specialCase.ToLowerRune(result)
		}
	}

	return result
}
*/

func (t *TabulaRecta) DecodeRune(r, key rune) rune {
	var result rune = r

	// @note IMPORTANT The German ß rune is always reported as
	// lowercase by Unicode. Converting it to Uppercase returns
	// ß in GO but 'SS' in other languages! That is causing an
	// issue in this code because IsLower(ß) says TRUE when all
	// TabulaRecta values are Uppercase!
	var isConvertedCase bool = false
	if t.caseFolding { // because the T.R. is in uppercase
		if t.specialCase == nil {
			if unicode.IsLower(r) {
				r = unicode.ToUpper(r)
				isConvertedCase = true
			}
			key = unicode.ToUpper(key)
		} else {
			if t.specialCase.IsLowerRune(r) {
				r = t.specialCase.ToUpperRune(r)
				isConvertedCase = true
			}
			key = t.specialCase.ToUpperRune(key)
		}
	}

	//if exists, keyIndex := t.rowContains(0, key); exists {
	if exists, keyIndex := t.HasRune(key); exists {
		exists, column := t.rowContains(keyIndex, r)
		if exists {
			result = t.tabula[0][column]
		}
	} else {
		mlog.WarnT("Key absent in alphabet", mlog.String("Alpha", t.Name), mlog.Rune("Rune", key))
	}

	if isConvertedCase { // respect input text's upper/lowercase
		if t.specialCase == nil {
			result = unicode.ToLower(result)
		} else {
			result = t.specialCase.ToLowerRune(result)
		}
	}

	return result
}

func (t *TabulaRecta) DecodeRuneRaw(r rune, rowIdx int) rune {
	var result rune = r

	if rowIdx >= len(t.tabula[0]) { // the Tabula Recta is a square matrix NxN
		mlog.ErrorT("out-of-range column", mlog.String("At", "DecodeRuneRaw"), mlog.Int("Value", rowIdx))
		panic("Bad thing happened")
	}

	// Pre
	isConvertedCase := t.caseFolding && ((t.specialCase == nil && unicode.IsLower(r)) ||
		(t.specialCase != nil && t.specialCase.IsLowerRune(r)))

	// Middle
	if exists, colIdx := t.rowContains(rowIdx, r); exists {
		result = t.tabula[0][colIdx]
	} else {
		// This would never happen UNLESS someone edits TabulaCaesar(Command) and didn't
		// check for rune's presence in the slave alphabet. But this check is here
		// as a safeguard.
		mlog.ErrorT("internal error: shouldn't be looking for that rune in this slave", mlog.Rune("Rune", r), mlog.String("Alpha", t.Name))
		panic("Internal Error")
	}

	// Post
	if isConvertedCase { // respect input text's upper/lowercase
		if t.specialCase == nil {
			result = unicode.ToLower(result)
		} else {
			result = t.specialCase.ToLowerRune(result)
		}
	}

	return result
}

func (t *TabulaRecta) DecodeRuneByShiftOld(r rune, shift int) rune { // @audit deprecate?
	var result rune = r

	// Pre
	max := utf8.RuneCountInString(t.alphabet)
	if shift >= max { // check just-in-case
		mlog.ErrorT("Cannot shift-encode beyond alphabet", mlog.Int("Shift", shift), mlog.Int("Max", max))
		panic("Bad thing happened")
	}

	var isConvertedCase bool = false
	if t.caseFolding { // because the T.R. is in uppercase
		if t.specialCase == nil {
			if unicode.IsLower(r) {
				r = unicode.ToUpper(r)
				isConvertedCase = true
			}
		} else {
			if t.specialCase.IsLowerRune(r) {
				r = t.specialCase.ToUpperRune(r)
				isConvertedCase = true
			}
		}
	}

	// Middle
	if exists, column := t.rowContains(shift, r); exists {
		result = t.tabula[0][column]
	} else {
		mlog.ErrorT("reference Slave alphabet does not contain rune", mlog.String("Alpha", t.Name), mlog.Rune("Rune", r))
	}

	// Post
	if isConvertedCase { // respect input text's upper/lowercase
		if t.specialCase == nil {
			result = unicode.ToLower(result)
		} else {
			result = t.specialCase.ToLowerRune(result)
		}
	}

	return result
}

/**
 * (IRuneLocalizer) Find a rune in the object's alphabet catalog.
 * Rune not found: error set, other return values nil or -1.
 * Rune found: error nil, pointer to alphabet and position within.
 */
func (t *TabulaRecta) FindRune(r rune) (string, int, error) {
	var alphaStr string = ""
	var at int
	var err error = nil

	//at = cmn.RuneIndex(t.alphabet, r)
	at = cmn.RuneIndexFold(t.alphabet, r, t.specialCase)
	if at == -1 {
		err = fmt.Errorf("info: '%c' absent in %s", r, t.Name)
	} else {
		alphaStr = t.alphabet
	}

	return alphaStr, at, err
}

func (t *TabulaRecta) TransposeKey(k any) (int, rune) {
	var transposed int = -1
	var resultingKey rune = 0
	switch cv := k.(type) {
	case rune: // called within Master
		if idx := cmn.RuneIndex(t.alphabet, cv); idx != -1 {
			transposed = idx
			resultingKey = cv
		} else {
			msg := "Cannot transpose key we don't have"
			mlog.ErrorT(msg, mlog.Rune("Key", cv))
			//panic(msg)
		}

	case int: // called within Slave
		max := utf8.RuneCountInString(t.alphabet)
		transposed = cv % max // limit it to the length of OUR alphabet
		resultingKey = cmn.RuneAt(t.alphabet, transposed)

	default:
		mlog.Error("Panicking with invalid type to transpose")
		panic("not a valid transpose parameter")
	}

	return transposed, resultingKey
}

func (t *TabulaRecta) IsCaseInsensitive() bool {
	return t.caseFolding
}

func (t *TabulaRecta) PrintTabulaRecta(center bool) {
	fmt.Println(t.renderTabulaRecta(center, true))
}

func (t *TabulaRecta) PrintTape(key rune) {
	if ok, pos := t.rowContains(0, key); ok {
		fmt.Println(t.renderTape(pos, true, true))
	}
}

func (t *TabulaRecta) String() string {
	return t.renderTabulaRecta(false, false)
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

/**
 * Given a valid Alphabet name, prepare a pre-configured Tabula Recta.
 */
func TabulaRectaForLanguage(lname string) *TabulaRecta {
	var alphabet *cmn.Alphabet = nil

	switch strings.ToLower(lname) {
	case cmn.ALPHA_NAME_ENGLISH:
		alphabet = cmn.ALPHA_DISK

	case cmn.ALPHA_NAME_SPANISH:
		fallthrough
	case cmn.ALPHA_NAME_LATIN:
		alphabet = cmn.ALPHA_DISK_LATIN

	case cmn.ALPHA_NAME_GREEK:
		alphabet = cmn.ALPHA_DISK_GREEK

	case cmn.ALPHA_NAME_GERMAN:
		alphabet = cmn.ALPHA_DISK_GERMAN

	case cmn.ALPHA_NAME_UKRAINIAN:
		fallthrough
	case cmn.ALPHA_NAME_RUSSIAN:
		fallthrough
	case cmn.ALPHA_NAME_CYRILLIC:
		alphabet = cmn.ALPHA_DISK_CYRILLIC

	case cmn.ALPHA_NAME_NUMBERS_ARABIC:
		alphabet = cmn.NUMBERS_DISK

	case cmn.ALPHA_NAME_NUMBERS_ARABIC_EXTENDED:
		alphabet = cmn.NUMBERS_DISK_EXT

	case cmn.ALPHA_NAME_NUMBERS_EASTERN:
		alphabet = cmn.NUMBERS_EASTERN_DISK

	case cmn.ALPHA_NAME_BINARY:
		mlog.Error("Please use NewBinaryTabulaRecta() instead of this")

	default:
		mlog.ErrorT("TabulaRecta requested for unknown alphabet", mlog.String("Alpha", lname))
	}

	tr := NewTabulaRecta(alphabet, true)
	return tr
}

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/

func DemoTabulaRecta(alphabet *cmn.Alphabet, foldCase bool, numerics *cmn.Alphabet) {
	trG := NewTabulaRecta(alphabet, foldCase)

	r := alphabet.GetRuneAt(2)
	trG.PrintTape(r)
	trG.PrintTabulaRecta(true)

	if numerics != nil {
		trN := NewTabulaRecta(numerics, false)
		trN.PrintTabulaRecta(true)
	}
}
