/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * (Command Pattern - See "Design Patterns")
 * A variation of Caesar that encodes even positions with the Prime Key
 * and odd positions with the Alternate Key.
 * 	The Alternate Key is the offset of the Prime Key plus the given Offset
 * modulo the length of the alphabet. If after the modulo operation the
 * Alternate Key offset is zero, it is bumped to the next letter because in
 * the Caesar algorithm the first letter represents no conversion.
 * 	The calculated position is based exclusively on convertable characters
 * (those within the selected alphabet). Therefore the actual position is
 * calculated as the position within the stream minus the number of skipped
 * letters at that moment in time.
 *	I created this variant just for fun, although it is almost equivalent
 * to a Bellaso cipher with key length = 2.
 *-----------------------------------------------------------------*/
package commands

import (
	"fmt"
	"lordofscripts/caesarx"
	"lordofscripts/caesarx/app"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/ciphers"
	"lordofscripts/caesarx/ciphers/caesar"
	"lordofscripts/caesarx/cmn"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	// Filename extension for files encrypted with Didimus
	FILE_EXT_DIDIMUS string = ".did"
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ ciphers.IPipe = (*DidimusCommand)(nil)
var _ ciphers.ICipherCommand = (*DidimusCommand)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type DidimusCommand struct {
	ciphers.Pipe
	core *caesar.DidimusTabulaRecta
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewDidimusCommand(alpha *cmn.Alphabet, key rune, offset uint8) *DidimusCommand {
	return &DidimusCommand{
		Pipe: ciphers.NewEmptyPipe(),
		core: caesar.NewDidimusTabulaRecta(alpha, key, offset),
	}
}

/* ----------------------------------------------------------------
 *							M e t h o d s (ICipherCommand)
 *-----------------------------------------------------------------*/

/**
 * Same as Rebuild() for this simple cipher.
 */
func (c *DidimusCommand) WithAlphabet(alphabet *cmn.Alphabet) ciphers.ICipherCommand {
	c.Rebuild(alphabet)
	return c
}

/**
 * Chain a slave disk/tabula that would use the same (or corrected) alphabet shift
 * as the main alphabet. A slave disk/tabula is usually a Numbers and/or Symbols
 * to IMPROVE the ancient cipher against attacks.
 */
func (c *DidimusCommand) WithChain(slave *cmn.Alphabet) ciphers.ICipherCommand {
	if slave != nil {
		slaveTR := ciphers.NewTabulaRecta(slave, true)
		c.core.WithChain(slaveTR)
	} else {
		c.core.WithChain(nil)
	}

	return c
}

func (c *DidimusCommand) Encode(plain string) (string, error) {
	err := c.core.VerifyKey()
	if err != nil {
		return "", err
	}

	ciphered := c.core.Encode(plain)
	if c.IsPipeOpen() {
		return c.PipeOutput(ciphers.PipeEncode, ciphered)
	} else {
		return ciphered, nil
	}
}

func (c *DidimusCommand) Decode(ciphered string) (string, error) {
	err := c.core.VerifyKey()
	if err != nil {
		return "", err
	}

	plain := c.core.Decode(ciphered)
	if c.IsPipeOpen() {
		return c.PipeOutput(ciphers.PipeDecode, plain)
	} else {
		return plain, nil
	}
}

// EncryptTextFile encrypts the filename src using the standard Caesar cipher.
// The output file has the FILE_EXT_DIDIMUS file extension. Please note that
// this method is only for text files.
func (c *DidimusCommand) EncryptTextFile(src string) error {
	var err error = nil
	if err = c.core.VerifyKey(); err == nil {
		fileOut := cmn.NewNameExtOnly(src, FILE_EXT_DIDIMUS, true)
		err = c.core.EncryptTextFile(src, fileOut) // error already logged by core
	}

	return err
}

// DecryptTextFile decrypts the filename src using the standard Caesar cipher.
// The output file target must be explicitely given. Please note that
// this method is only for text files.
func (c *DidimusCommand) DecryptTextFile(src, target string) error {
	var err error = nil
	if err = c.core.VerifyKey(); err == nil {
		err = c.core.DecryptTextFile(src, target) // error already logged by core
	}

	return err
}

// Encodes a binary file and produces a binary encoded file
func (c *DidimusCommand) EncryptBinFile(filenameIn string) error {
	var err error = nil
	if err = c.core.VerifyKey(); err == nil {
		fileOut := cmn.NewNameExtOnly(filenameIn, FILE_EXT_DIDIMUS, true)
		err = c.core.EncryptBinaryFile(filenameIn, fileOut) // error already logged by core
	}

	return err
}

// Decodes a binary file and produces a plain binary file
func (c *DidimusCommand) DecryptBinFile(filenameIn, filenameOut string) error {
	var err error = nil
	if err = c.core.VerifyKey(); err == nil {
		err = c.core.DecryptBinaryFile(filenameIn, filenameOut) // error already logged by core
	}

	return err
}

func (c *DidimusCommand) Alphabet() string {
	return c.core.GetAlphabet()
}

/**
 * Checks the alphabet, if OK it is applied to the underlying cipher machine.
 * Else it logs an error and exits with ERR_BAD_ALPHABET.
 */
func (c *DidimusCommand) Rebuild(alphabet *cmn.Alphabet, opts ...any) {
	if alphabet.Check() {
		c.core.WithAlphabet(alphabet)
	} else {
		err := fmt.Errorf("invalid alphabet '%s' size:%d", alphabet.Name, alphabet.Size())
		mlog.ErrorE(err)
		app.DieWithError(err, caesarx.ERR_BAD_ALPHABET)
	}
}

func (c *DidimusCommand) String() string {
	return fmt.Sprintf("%s %s", c.core.GetLanguage(), c.core.String())
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/
func DemoDidimusCommand(alpha, numbers *cmn.Alphabet, phrase string) bool {
	fmt.Println("Didimus Encryption (Command-pattern version)")
	fmt.Println("( A variant of the Caesar family )")
	fmt.Println("Master : ", alpha.Name)
	if numbers != nil {
		fmt.Println("Slave  : ", numbers.Name)
	}
	var ok bool = true
	for _, key := range []rune{alpha.GetRuneAt(10), alpha.GetRuneAt(-1)} {
		const OFFSET uint8 = 5
		var encTxt, encTxt2, decTxt string
		var err error

		ngram := cmn.NewNgramFormatter(5, '·')
		cnv1 := NewDidimusCommand(alpha, key, OFFSET)
		cnv2 := NewDidimusCommand(alpha, key, OFFSET)
		cnv2.WithPipe(ngram)

		encTxt, err = cnv1.Encode(phrase)
		if err != nil {
			app.DieWithError(err, caesarx.ERR_DEMO_ERROR)
		}
		encTxt2, err = cnv2.Encode(phrase)
		if err != nil {
			app.DieWithError(err, caesarx.ERR_DEMO_ERROR)
		}

		decTxt, err = cnv1.Decode(encTxt)
		if err != nil {
			app.DieWithError(err, caesarx.ERR_DEMO_ERROR)
		}

		fmt.Printf("Key    : %c (shift=%d)\n", key, alpha.PositionOf(key))
		fmt.Println("Plain  : ", phrase)
		fmt.Println("Encoded: ", encTxt)
		fmt.Println("Format : ", encTxt2)
		fmt.Println("Decoded: ", decTxt)

		if decTxt != phrase {
			ok = false
		}

		fmt.Println()
	}

	return ok
}
