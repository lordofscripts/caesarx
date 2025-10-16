package caesar

import (
	"fmt"
	"lordofscripts/caesarx"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/ciphers"
	"lordofscripts/caesarx/internal/crypto"
	"reflect"
	"strings"
	"unicode/utf8"
)

type alphaRune struct {
	Rune  rune
	Shift int
}

type TextIterator struct {
	tabulas   []ciphers.ITabulaRecta
	sequencer crypto.IKeySequencer
	sb        strings.Builder
	pos       int
	dataPtr   *string
	max       int
}

/**
 * (ctor) regular constructor when the rune has a defined position
 */
func newAlphaRune(r rune, at int) *alphaRune {
	return &alphaRune{Rune: r, Shift: at}
}

/**
 * (ctor) extraordinary constructor when the rune isn't present
 */
func newAlphaRuneNotFound(r rune) *alphaRune {
	return &alphaRune{Rune: r, Shift: -1}
}

func NewTextIterator(sx crypto.IKeySequencer, tabs ...ciphers.ITabulaRecta) *TextIterator {
	tabulas := make([]ciphers.ITabulaRecta, 0)
	for _, tab := range tabs {
		if !isNil(tab) {
			tabulas = append(tabulas, tab)
		}
	}

	return &TextIterator{
		tabulas:   tabulas,
		sequencer: sx,
		pos:       -1,
		dataPtr:   nil,
		max:       -1,
	}
}

func (a *alphaRune) String() string {
	return fmt.Sprintf("%c - %d", a.Rune, a.Shift)
}

func (t *TextIterator) Start(s string) {
	t.dataPtr = &s
	t.pos = 0
	t.max = utf8.RuneCountInString(s)
}

func (t *TextIterator) getRuneAt(s string, nr int) (rune, bool) {
	if nr < 0 {
		return 0, false // @audit I guess this should panic
	}

	count := 0
	for _, r := range s {
		if count == nr {
			return r, true
		}
		count++
	}
	return 0, false
}

func (t *TextIterator) EncodeNext() bool {
	var result rune

	if t.pos >= t.max {
		return true
	}

	// Tabula index, Rune to en/decode, Current key, Shift within alphabet
	// currTabula, targetChar, currKey, targetShift
	if currTab, target, currKey := t.next(); currTab != -1 {
		if err := t.sequencer.Feedback(target.Rune); err != nil {
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
			//colIdx, _ := t.tabulas[trIdx].TransposeKey(shift)
			//result = t.tabulas[trIdx].EncodeRuneRaw(char, shift, colIdx)
		}

		t.sb.WriteRune(result)
	} else {
		t.sequencer.Skip()
		t.sb.WriteRune(target.Rune)
	}

	return t.pos == t.max
}

func (t *TextIterator) DecodeNext() bool {
	var result rune

	if t.pos >= t.max {
		return true
	}

	// Tabula index, Rune to en/decode, Current key, Shift within alphabet
	if currTab, target, keyCurr := t.next(); currTab != -1 {
		if currTab == 0 { // Primary alphabet
			result = t.tabulas[0].DecodeRune(target.Rune, keyCurr.Rune)
		} else {
			//_, keyNew := t.tabulas[trIdx].TransposeKey(shift) @audit problem here
			//result = t.tabulas[trIdx].DecodeRune(char, keyNew)
			//rowIdx, _ := t.tabulas[trIdx].TransposeKey(shift)
			//result = t.tabulas[trIdx].DecodeRuneRaw(char, rowIdx)
			_, newKey := t.tabulas[currTab].TransposeKey(keyCurr.Shift)
			result = t.tabulas[currTab].DecodeRune(target.Rune, newKey)
		}

		if err := t.sequencer.Feedback(result); err != nil {
			mlog.FatalT(caesarx.ERR_SEQUENCER,
				"feedback panic",
				mlog.Int("Pos", t.pos),
				mlog.At(),
			)
		}
		t.sb.WriteRune(result)
	} else {
		t.sequencer.Skip()
		t.sb.WriteRune(target.Rune)
	}

	return t.pos == t.max
}

func (t *TextIterator) next() (int, *alphaRune, *alphaRune) {
	var currTabula int = -1 // Tabula index
	var target *alphaRune = nil
	var currKey *alphaRune = nil

	// Internal utility function to locate a character in a Tabula.
	// Can be used for both targetChar & key. It also returns the
	// ID (index) of the TabulaRecta where it was found.
	locatorFx := func(rx rune) (int, *alphaRune) {
		var result *alphaRune = nil
		var id int = -1

		for idxTab, tr := range t.tabulas {
			if !isNil(tr) {
				if ok, shift := tr.HasRune(rx); ok {
					id = idxTab
					result = newAlphaRune(rx, shift)
					break
				}
			}
		}

		return id, result
	}

	if targetChar, ok := t.getRuneAt(*t.dataPtr, t.pos); ok {
		key := t.sequencer.GetKey(t.pos, targetChar) // always from Primary alphabet
		keyTabulaId, keyInfo := locatorFx(key)
		if keyTabulaId != -1 {
			currKey = keyInfo
		} else {
			mlog.FatalT(70, "Key's rune not found in tabulae",
				mlog.String("At", "TextIterator.next()"),
				mlog.Rune("Key", key),
				mlog.Int("Pos", t.pos))
		}
		/*
			if _, keyAt, err := t.tabulas[0].FindRune(key); err != nil {
				mlog.FatalT(70, "Key's rune not found in primary",
					mlog.String("Alpha", t.tabulas[0].GetName()),
					mlog.String("At", "TextIterator.next()"),
					mlog.Rune("Key", key),
					mlog.Int("Pos", t.pos))
			} else {
				currKey = newAlphaRune(key, keyAt)
			}
		*/
		currTabula, target = locatorFx(targetChar)
		/*
			for idxTab, tr := range t.tabulas {
				if !isNil(tr) {
					if ok, shift := tr.HasRune(targetChar); ok {
						currTabula = idxTab
						target = newAlphaRune(targetChar, shift)
						break
					}
				}
			}
		*/
		if target == nil {
			target = newAlphaRuneNotFound(targetChar)
		}
	}
	t.pos++

	return currTabula, target, currKey
}

func (t *TextIterator) Result() string {
	output := t.sb.String()
	t.sb.Reset()
	return output
}

/**
 * to correctly check that a variable of type Interface is nil. Else a
 * normal a == nil would fail as in this case:
 *		var a, b IAnyInterface
 *		a = new(TypeWhichImplementsInterface)
 *		b = nil
 *		if b == nil { fmt.Print("This will print!") }
 * That's because an interface variable has both a type and a value,
 * the value is nil but the type is not and the normal == will see
 * that its type is not nil and therefore it will print!
 */
func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
