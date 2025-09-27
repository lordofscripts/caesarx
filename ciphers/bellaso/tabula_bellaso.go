/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Caesar Encoding using a Tabula Recta where the key is based on
 * the Fibonacci series terms that correspond to a legal key within
 * the realm of the chosen alphabet.
 *-----------------------------------------------------------------*/
package bellaso

import (
	"lordofscripts/caesarx/ciphers"
	"lordofscripts/caesarx/ciphers/caesar"
	"lordofscripts/caesarx/cmn"
	iciphers "lordofscripts/caesarx/internal/ciphers"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/
var (
	Info = ciphers.NewCipherInfo(iciphers.ALG_CODE_BELLASO, "0.9",
		"Giovan Battista Bellaso",
		iciphers.ALG_NAME_BELLASO,
		"Bellaso polyalphabetic cipher")
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

var _ ciphers.ICipher = (*BellasoTabulaRecta)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type BellasoTabulaRecta struct {
	caesar.CaesarTabulaRecta
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/**
 * (Ctor) Caesar Cipher using a Tabula Recta that supports ASCII and
 * foreign (UTF8) alphabets. Uses a 10-term Fibonacci series as offset
 * to the Prime key. Has an Extended Numeric alphabet.
 * · Always follow it with a call to VerifyKey() or VerifySecret() prior to
 *	 begining encoding/decoding.
 * · follow with WithChain() to chain with supplemental alphabets.
 * · follow with WithAlphabet() to specify a different alphabet prior to encoding.
 * · It does case-folding by default, so it handles & preserves upper/lowercase
 */
func NewBellasoTabulaRecta(alphabet *cmn.Alphabet, secret string) *BellasoTabulaRecta {
	base := caesar.NewCaesarTabulaRecta(alphabet, alphabet.GetRuneAt(0))
	base.WithSequencer(iciphers.NewBellasoSequencer(secret, alphabet))

	giovanni := &BellasoTabulaRecta{*base}
	giovanni.WithChain(ciphers.NewTabulaRecta(cmn.NUMBERS_DISK_EXT, true))

	return giovanni
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (cx *BellasoTabulaRecta) String() string {
	return iciphers.ALG_NAME_BELLASO
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/
