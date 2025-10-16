/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Caesar Encoding using a Tabula Recta where the key is based on
 * the Fibonacci series terms that correspond to a legal key within
 * the realm of the chosen alphabet.
 *-----------------------------------------------------------------*/
package caesar

import (
	"lordofscripts/caesarx/ciphers"
	"lordofscripts/caesarx/cmn"
	"lordofscripts/caesarx/internal/crypto"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/
var (
	InfoFibonacci = ciphers.NewCipherInfo(crypto.ALG_CODE_FIBONACCI, "1.0",
		"Didimo Grimaldo",
		crypto.ALG_NAME_FIBONACCI,
		"Fibonacci polyalphabetic cipher")
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/
func init() {
	ciphers.RegisterCipher(InfoFibonacci)
}

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ ciphers.ICipher = (*FibonacciTabulaRecta)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type FibonacciTabulaRecta struct {
	CaesarTabulaRecta
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
func NewFibonacciTabulaRecta(alphabet *cmn.Alphabet, primeKey rune) *FibonacciTabulaRecta {
	base := NewCaesarTabulaRecta(alphabet, alphabet.GetRuneAt(0))
	base.sequencer = crypto.NewFibonacciSequencer(alphabet, primeKey)
	fibo := &FibonacciTabulaRecta{*base}
	fibo.WithChain(ciphers.NewTabulaRecta(cmn.NUMBERS_DISK_EXT, true))
	return fibo
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (cx *FibonacciTabulaRecta) String() string {
	return cx.sequencer.GetKeyInfo()
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/
