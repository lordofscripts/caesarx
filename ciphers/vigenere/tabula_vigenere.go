/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Caesar Encoding using a Tabula Recta where the key is based on
 * the Fibonacci series terms that correspond to a legal key within
 * the realm of the chosen alphabet.
 * NOTE: For Vigenere, we must instruct the sequencer (VigenereSequencer)
 *	  via the parent (CaesarTabulaRecta), whether it is sequencing
 *	  an encryption or decryption. For (Vigenere) AutoKey the key
 *	  is progressively built, whereas during encryption, the autokey
 *	  is known beforehand (secret+plain) *-----------------------------------------------------------------*/
package vigenere

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
	Info = ciphers.NewCipherInfo(iciphers.ALG_CODE_VIGENERE, "1.0",
		"Blaise de Vigenère",
		iciphers.ALG_NAME_VIGENERE,
		"Vigenère polyalphabetic cipher")
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

var _ ciphers.ICipher = (*VigenereTabulaRecta)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type VigenereTabulaRecta struct {
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
func NewVigenereTabulaRecta(alphabet *cmn.Alphabet, secret string) *VigenereTabulaRecta {
	base := caesar.NewCaesarTabulaRecta(alphabet, alphabet.GetRuneAt(0))
	base.WithSequencer(iciphers.NewVigenereSequencer(secret, alphabet))

	vige := &VigenereTabulaRecta{*base}
	vige.WithChain(ciphers.NewTabulaRecta(cmn.NUMBERS_DISK_EXT, true))

	return vige
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (cx *VigenereTabulaRecta) String() string {
	return iciphers.ALG_NAME_VIGENERE
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/
