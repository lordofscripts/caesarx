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
	// Pipe this command's output to the input of another command
	IPipe
	// Set the primary/master alphabet to substitute the one given in
	// the constructor.
	WithAlphabet(alphabet *cmn.Alphabet) ICipherCommand
	// Set a secondary/slave alphabet to use in case target is not
	// found in primary/master alphabet.
	WithChain(*cmn.Alphabet) ICipherCommand

	// Encodes/Encrypts a string
	Encode(plain string) (string, error)
	// Decodes/Decrypts a string
	Decode(ciphered string) (string, error)

	// Encodes a text file and produces an output file.
	EncryptTextFile(filenameIn string) error
	// Decodes a text file to produce a plain-text file.
	DecryptTextFile(filenameIn, filenameOut string) error
	// Encodes a binary file and produces a binary encoded file
	EncryptBinFile(filenameIn string) error
	// Decodes a binary file and produces a plain binary file
	DecryptBinFile(filenameIn, filenameOut string) error

	// get the output filename when it was inferred
	GetOutputFilename() string
	// Get the alphabet string (don't use it for binary alphabets)
	Alphabet() string
	// deprecated
	Rebuild(alphabet *cmn.Alphabet, opts ...any)
	// String representation of the command
	fmt.Stringer
}
