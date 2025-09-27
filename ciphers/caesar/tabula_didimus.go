/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package caesar

import (
	"lordofscripts/caesarx/ciphers"
	"lordofscripts/caesarx/cmn"
	iciphers "lordofscripts/caesarx/internal/ciphers"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/
var (
	InfoDidimus = ciphers.NewCipherInfo(iciphers.ALG_CODE_DIDIMUS, "1.0",
		"Didimo Grimaldo",
		iciphers.ALG_NAME_DIDIMUS,
		"Didimus polyalphabetic cipher")
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/
func init() {
	ciphers.RegisterCipher(InfoDidimus)
}

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ ciphers.ICipher = (*DidimusTabulaRecta)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type DidimusTabulaRecta struct {
	CaesarTabulaRecta
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/**
 * (Ctor) Caesar Cipher using a Tabula Recta that supports ASCII and
 * foreign (UTF8) alphabets. Uses a Prime Key and an Alternate Key
 * for odd/even and an Extended Numeric Alphabet.
 * · Always follow it with a call to VerifyKey() or VerifySecret() prior to
 *	 begining encoding/decoding.
 * · follow with WithChain() to chain with supplemental alphabets.
 * · follow with WithAlphabet() to specify a different alphabet prior to encoding.
 * · It does case-folding by default, so it handles & preserves upper/lowercase
 */
func NewDidimusTabulaRecta(alphabet *cmn.Alphabet, primeKey rune, offset uint8) *DidimusTabulaRecta {
	base := NewCaesarTabulaRecta(alphabet, primeKey)
	base.sequencer = iciphers.NewDidimusSequencer(primeKey, offset, alphabet)
	didimus := &DidimusTabulaRecta{*base}
	didimus.WithChain(ciphers.NewTabulaRecta(cmn.NUMBERS_DISK_EXT, true))

	return didimus
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (cx *DidimusTabulaRecta) String() string {
	return cx.sequencer.GetKeyInfo()
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/
