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
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

type ICipherCommand interface {
	IPipe
	WithAlphabet(alphabet *cmn.Alphabet) ICipherCommand
	WithChain(*cmn.Alphabet) ICipherCommand
	Encode(plain string) (string, error)
	Decode(ciphered string) (string, error)
	Alphabet() string
	//Rebuild(alphabet string, opts ...any)
	Rebuild(alphabet *cmn.Alphabet, opts ...any)
	fmt.Stringer
}
