/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Plain Caesar cipher using Tabula Recta implementation. Is
 * case-insensitive but preserves case.
 *-----------------------------------------------------------------*/
package caesar

import (
	"fmt"
	"lordofscripts/caesarx/ciphers"
	"lordofscripts/caesarx/cmn"
	iciphers "lordofscripts/caesarx/internal/ciphers"
	"strings"
	"unicode/utf8"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/
var (
	Info = ciphers.NewCipherInfo(iciphers.ALG_CODE_CAESAR, "1.0",
		"Julius Caesar",
		iciphers.ALG_NAME_CAESAR,
		"Caesar cipher")
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/
func init() {
	ciphers.RegisterCipher(Info)
}

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ ciphers.ICipher = (*CaesarTabulaRecta)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type CaesarTabulaRecta struct {
	alpha     *cmn.Alphabet
	slave     *ciphers.TabulaRecta // implements cmn.IRuneLocalizer
	sequencer iciphers.IKeySequencer
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/**
 * (Ctor) Caesar Cipher using a Tabula Recta that supports ASCII and
 * foreign (UTF8) alphabets.
 * · Always follow it with a call to VerifyKey() or VerifySecret() prior to
 *	 begining encoding/decoding.
 * · follow with WithChain() to chain with supplemental alphabets.
 * · follow with WithAlphabet() to specify a different alphabet prior to encoding.
 * · It does case-folding by default, so it handles & preserves upper/lowercase
 */
func NewCaesarTabulaRecta(alphabet *cmn.Alphabet, key rune) *CaesarTabulaRecta {
	return &CaesarTabulaRecta{
		alpha:     alphabet,
		slave:     nil,
		sequencer: iciphers.NewCaesarSequencer(key),
	}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (cx *CaesarTabulaRecta) WithChain(extra *ciphers.TabulaRecta) ciphers.ICipher {
	cx.slave = extra
	return cx
}

func (cx *CaesarTabulaRecta) WithAlphabet(alphabet *cmn.Alphabet) ciphers.ICipher {
	cx.alpha = alphabet
	return cx
}

func (cx *CaesarTabulaRecta) WithSequencer(keygen iciphers.IKeySequencer) ciphers.ICipher {
	cx.sequencer = keygen
	return cx
}

/**
 * Verify key(s). If none given it checks the key given in the constructor,
 * else it checks all the given keys. The key (single character) must be
 * present in the encoding alphabet.
 */
func (cx *CaesarTabulaRecta) VerifyKey(keys ...rune) error {
	verify := func(k rune) error {
		if !cx.alpha.Contains(k, cmn.CaseInsensitive) { // @audit what if TR is not case folded!
			return fmt.Errorf("key '%c' is not part of the alphabet", k)
		}
		return nil
	}

	if len(keys) == 0 {
		//return verify(cx.sequencer.GetKey(DUMMY_POS, DUMMY_RUNE))
		return cx.sequencer.Verify(verify)
	} else {
		for i, key := range keys {
			err := verify(key)
			if err != nil {
				return fmt.Errorf("key #%d (%c): %w", i+1, key, err)
			}
		}
	}

	return nil
}

/**
 * Does the same as VerifyKey() except it checks that all the given keys
 * (if any) have exactly ONE character (could be multi-byte Unicode character).
 */
func (cx *CaesarTabulaRecta) VerifySecret(secrets ...string) error {
	if len(secrets) == 0 {
		return cx.VerifyKey()
	}

	allKeys := make([]rune, len(secrets))
	for i, keyStr := range secrets {
		if utf8.RuneCountInString(keyStr) != 1 {
			return fmt.Errorf("key #%d '%s' contains more than one char", i+1, keyStr)
		}

		allKeys[i] = []rune(keyStr)[0]
	}

	return cx.VerifyKey(allKeys...)
}

/**
 * (IRuneLocalizer) Find a rune in the object's alphabet catalog.
 * Rune not found: error set, other return values nil or -1.
 * Rune found: error nil, pointer to alphabet and position within.
 * This method seeks in both the Primary/Master & Secondary/Slave.
 */
func (t *CaesarTabulaRecta) FindRune(r rune) (alpha string, at int, err error) {
	err = nil
	//at = cmn.RuneIndex(t.alpha.Chars, r)
	at = cmn.RuneIndexFold(t.alpha.Chars, r, t.alpha.BorrowSpecialCase())
	if at == -1 {
		// not present in the Primary/Master alphabet, let's try the Slave
		if t.slave == nil {
			// bummer! there is no Secondary/Slave, ran out of options
			err = fmt.Errorf("info: '%c' absent in %s", r, t.alpha.Name)
			alpha = ""
		} else {
			// perhaps in the Slave (or not) and that's final
			alpha, at, err = t.slave.FindRune(r)
		}
	} else {
		alpha = t.alpha.Chars
	}

	return
}

func (cx *CaesarTabulaRecta) Encode(plain string) string {
	master := ciphers.NewTabulaRecta(cx.alpha, cmn.CaseInsensitive)
	cx.sequencer.SetDecryptionMode(false) // only matters with Vigenere
	iter := NewTextIterator(cx.sequencer, master, cx.slave)
	iter.Start(plain)
	for !iter.EncodeNext() {
		//fmt.Print("E")
	}
	//fmt.Println()

	cx.sequencer.Reset()
	return iter.Result()
}

func (cx *CaesarTabulaRecta) EncodeOld2(plain string) string {
	var master *ciphers.TabulaRecta = nil
	var sb strings.Builder

	// Plain Caesar uses ONE key; therefore, we only need one Tabula Recta
	var cnv ciphers.ITabulaRecta
	var masterRules bool
	var columnIdx int

	cmdletPassThrough := func(char rune, sx iciphers.IKeySequencer) {
		sb.WriteRune(char)
		sx.Skip()
	}

	master = ciphers.NewTabulaRecta(cx.alpha, cmn.CaseInsensitive)
	plainR := []rune(plain)
	// @note go through every RUNE in the input stream (not bytes!). In GO if
	// you range through a string, the 'pos' will be a BYTE position and not
	// a RUNE position. Else it would go wrong whenever 'plainR' has multi-byte runes.
	for pos, char := range plainR {
		var exists bool
		// input rune exists on Master alphabet?
		exists, columnIdx = master.HasRune(char)
		if !exists {
			if cx.slave != nil {
				// perhaps the rune exists on the Slave alphabet?
				exists, columnIdx = cx.slave.HasRune(char)
				if !exists {
					// nah! pass it through un-encoded
					cmdletPassThrough(char, cx.sequencer)
					continue
				}

				cnv = cx.slave // slave 'disk' will translate
				masterRules = false
			} else {
				cmdletPassThrough(char, cx.sequencer)
				continue
			}
		} else {
			cnv = master // master 'disk' will translate
			masterRules = true
		}

		var dummyFunc = func(c int) {

		}
		dummyFunc(columnIdx)

		// let the cipher algorithm sequencer determine the current Key
		// The sequencer determines whether either, both or none of the
		// parameters are needed to derive the current Key.
		/* @audit remove
		fmt.Println("CaesarTabulaRecta.Encode() Pos#", pos)
		if pos == 7 {
			fmt.Println("Bad Thing will happen")
		}
		*/

		// GetKey() is alphabet agnostic. Returns the key to use for encoding.
		// But before using it a decision must be made as follows:
		//	in Master : use EncodeRune() with key as-is
		//	in Slave  :
		//	in neither: pass-through unchanged
		key := cx.sequencer.GetKey(pos, char)
		var encR rune

		if !masterRules {
			/*
				// we need to transpose the Master's key onto the Slave's alphabet
				keyShift := master.TransposeKey(key) // get the numeric shift from the master
				if keyShift == -1 {
					keyShift = cx.slave.TransposeKey(key)
					if keyShift == -1 {
						cmdletPassThrough(char, cx.sequencer)
						continue
					}
				}
				rowIdx := cx.slave.TransposeKey(keyShift)
				//fmt.Printf("Encode %c at %d\n", char, pos)
				encR = cnv.EncodeRuneRaw(char, rowIdx, columnIdx)
			*/
			encR = cx.slave.EncodeRune(char, key)
		} else {
			encR = cnv.EncodeRune(char, key)
		}
		sb.WriteRune(encR)
	}

	return sb.String()
}

func (cx *CaesarTabulaRecta) Decode(plain string) string {
	master := ciphers.NewTabulaRecta(cx.alpha, cmn.CaseInsensitive)
	cx.sequencer.SetDecryptionMode(true) // only matters with Vigenere

	iter := NewTextIterator(cx.sequencer, master, cx.slave)
	iter.Start(plain)
	for !iter.DecodeNext() {
		//fmt.Print("D")
	}
	//fmt.Println()

	cx.sequencer.Reset()
	return iter.Result()
}

func (cx *CaesarTabulaRecta) DecodeOld2(cipher string) string {
	var sb strings.Builder

	var cnv ciphers.ITabulaRecta
	var master *ciphers.TabulaRecta = nil
	var masterRules bool

	cmdletPassThrough := func(char rune, sx iciphers.IKeySequencer) {
		sb.WriteRune(char)
		sx.Skip()
	}

	master = ciphers.NewTabulaRecta(cx.alpha, cmn.CaseInsensitive)
	cipherR := []rune(cipher)

	for pos, char := range cipherR {
		var exists bool
		// does the encoded rune exist in the master alphabet?
		exists, _ = master.HasRune(char)
		if !exists {
			if cx.slave != nil {
				// perhaps the slave alphabet disk?
				exists, _ = cx.slave.HasRune(char)
				if !exists {
					// nah, pass it undecoded
					cmdletPassThrough(char, cx.sequencer)
					continue
				}

				cnv = cx.slave
				masterRules = false
			} else {
				cmdletPassThrough(char, cx.sequencer)
				continue
			}
		} else {
			cnv = master
			masterRules = true
		}

		var decR rune
		key := cx.sequencer.GetKey(pos, char)

		if !masterRules {
			// we need to transpose the Master's key onto the Slave's alphabet
			keyShift, _ := master.TransposeKey(key) // shifted amount on Master
			if keyShift == -1 {
				keyShift, _ = cx.slave.TransposeKey(key)
				if keyShift == -1 {
					cmdletPassThrough(char, cx.sequencer)
					continue
				}
			}

			rowIdx, _ := cx.slave.TransposeKey(keyShift) // transpose to Slave
			decR = cnv.DecodeRuneRaw(char, rowIdx)
		} else {
			decR = cnv.DecodeRune(char, key)
		}
		sb.WriteRune(decR)
	}

	return sb.String()
}

/*
func (cx *CaesarTabulaRecta) Decode(cipher string) string {
	var cnv *ciphers.TabulaRecta = nil
	var sb strings.Builder

	cnv = ciphers.NewTabulaRecta(cx.alpha, cmn.CaseInsensitive)
	cipherR := []rune(cipher)
	for pos, char := range cipherR {
		if char == ' ' {
			mlog.Error("Bad thing about to happen")
		}
		if exists, _ := cnv.HasRune(char); !exists {
			sb.WriteRune(char)
			cx.sequencer.Skip()
		} else {
			//const DUMMY_POS = -1
			//const DUMMY_RUNE = ' '
			key := cx.sequencer.GetKey(pos, char)
			decR := cnv.DecodeRune(char, key)
			sb.WriteRune(decR)
		}
	}

	return sb.String()
}
*/

func (cx *CaesarTabulaRecta) GetAlphabet() string {
	return cx.alpha.Chars
}

func (cx *CaesarTabulaRecta) GetLanguage() string {
	return cx.alpha.Name
}

func (cx *CaesarTabulaRecta) String() string {
	return cx.sequencer.GetKeyInfo()
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/

func DemoCaesarPlain(alpha *cmn.Alphabet, message string) {
	var keyNoOp = alpha.GetRuneAt(0)
	var keyMidl = alpha.GetRuneAt(int(alpha.Size() / 2))
	var keyLast = alpha.GetRuneAt(-1)
	var multiByte = utf8.RuneCountInString(alpha.Chars) > len(alpha.Chars)

	doRoundTrip := func(key rune) {
		cnv := NewCaesarTabulaRecta(alpha, key)
		cipher := cnv.Encode(message)
		plain := cnv.Decode(cipher)
		passed := plain == message

		fmt.Printf("\tKey       : %c\n", key)
		fmt.Printf("\tMessage   : %s\n", message)
		fmt.Printf("\tCiphered  : %s\n", cipher)
		fmt.Printf("\tDeciphered: %s\n", plain)
		fmt.Printf("\tPassed    : %t\n\n", passed)
	}

	fmt.Println("Plain Caesar Demo with Tabula Recta")
	fmt.Printf("\tAlphabet  : %s\n", alpha.Name)
	fmt.Printf("\tMulti-byte: %t\n", multiByte)

	doRoundTrip(keyNoOp)
	doRoundTrip(keyMidl)
	doRoundTrip(keyLast)
}
