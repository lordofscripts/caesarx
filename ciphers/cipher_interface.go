/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package ciphers

import (
	"lordofscripts/caesarx/cmn"
	iciphers "lordofscripts/caesarx/internal/ciphers"
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

type ICipher interface {
	cmn.IRuneLocalizer
	WithChain(*TabulaRecta) ICipher
	WithAlphabet(alphabet *cmn.Alphabet) ICipher
	WithSequencer(iciphers.IKeySequencer) ICipher
	VerifySecret(secret ...string) error

	Encode(plain string) string
	Decode(cipher string) string
}
