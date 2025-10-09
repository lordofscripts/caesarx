/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package ciphers

import (
	"fmt"
	"lordofscripts/caesarx/cmn"
	iciphers "lordofscripts/caesarx/internal/ciphers"
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

type ICipher interface {
	// Setup methods
	WithChain(*TabulaRecta) ICipher
	WithAlphabet(alphabet *cmn.Alphabet) ICipher
	WithSequencer(iciphers.IKeySequencer) ICipher

	// Queries
	cmn.IRuneLocalizer
	VerifySecret(secret ...string) error
	VerifyKey(keys ...rune) error
	GetAlphabet() string
	GetLanguage() string

	// Execution methods
	Encode(plain string) string
	Decode(cipher string) string
	EncryptTextFile(fileIn, fileOut string) error
	DecryptTextFile(fileIn, fileOut string) error

	fmt.Stringer
}
