/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * (Command Pattern - See "Design Patterns")
 * The same Plain Caesar cipher based on the TabulaRecta implementation,
 * but wrapped as a Command Pattern.
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
	// Filename extension for files encrypted with plain Caesar
	FILE_EXT_CAESAR string = ".cae"
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ ciphers.IPipe = (*CaesarCommand)(nil)
var _ ciphers.ICipherCommand = (*CaesarCommand)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type CaesarCommand struct {
	ciphers.Pipe
	//core *caesar.CaesarTabulaRecta
	core ciphers.ICipher
	opts *caesar.CaesarOptions
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewCaesarCommand(alpha *cmn.Alphabet, key rune) *CaesarCommand {
	return &CaesarCommand{
		Pipe: ciphers.NewEmptyPipe(),
		core: caesar.NewCaesarTabulaRecta(alpha, key),
		opts: caesar.NewCaesarOpts(key),
	}
}

func NewCaesarCommandWithOptions(alpha *cmn.Alphabet, opts *caesar.CaesarOptions) *CaesarCommand {
	if opts.Variant != caesar.CAESAR { // @audit expand as I add more TESTED Caesar variant implementations
		mlog.Error("invalid Caesar variant for CaesarCommand")
		return nil
	}

	return &CaesarCommand{
		core: caesar.NewCaesarTabulaRecta(alpha, opts.Initial),
		opts: opts,
	}
}

/* ----------------------------------------------------------------
 *							M e t h o d s (ICipherCommand)
 *-----------------------------------------------------------------*/

/**
 * Same as Rebuild() for this simple cipher.
 */
func (c *CaesarCommand) WithAlphabet(alphabet *cmn.Alphabet) ciphers.ICipherCommand {
	c.Rebuild(alphabet)
	return c
}

/**
 * Chain a slave disk/tabula that would use the same (or corrected) alphabet shift
 * as the main alphabet. A slave disk/tabula is usually a Numbers and/or Symbols
 * to IMPROVE the ancient cipher against attacks.
 */
func (c *CaesarCommand) WithChain(slave *cmn.Alphabet) ciphers.ICipherCommand {
	if slave != nil {
		slaveTR := ciphers.NewTabulaRecta(slave, true)
		c.core.WithChain(slaveTR)
	} else {
		c.core.WithChain(nil)
	}

	return c
}

func (c *CaesarCommand) Encode(plain string) (string, error) {
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

func (c *CaesarCommand) Decode(ciphered string) (string, error) {
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
// The output file has the FILE_EXT_CAESAR file extension. Please note that
// this method is only for text files.
func (c *CaesarCommand) EncryptTextFile(src string) error {
	var err error = nil
	if err = c.core.VerifyKey(); err == nil {
		fileOut := cmn.NewNameExtOnly(src, FILE_EXT_CAESAR, true)
		err = c.core.EncryptTextFile(src, fileOut) // error already logged by core
	}

	return err
}

// DecryptTextFile decrypts the filename src using the standard Caesar cipher.
// The output file target must be explicitely given. Please note that
// this method is only for text files.
func (c *CaesarCommand) DecryptTextFile(src, target string) error {
	var err error = nil
	if err = c.core.VerifyKey(); err == nil {
		err = c.core.DecryptTextFile(src, target) // error already logged by core
	}

	return err
}

// Encodes a binary file and produces a binary encoded file
func (c *CaesarCommand) EncryptBinFile(filenameIn string) error {
	var err error = nil
	if err = c.core.VerifyKey(); err == nil {
		fileOut := cmn.NewNameExtOnly(filenameIn, FILE_EXT_CAESAR, true)
		err = c.core.EncryptBinaryFile(filenameIn, fileOut) // error already logged by core
	}

	return err
}

// Decodes a binary file and produces a plain binary file
func (c *CaesarCommand) DecryptBinFile(filenameIn, filenameOut string) error {
	var err error = nil
	if err = c.core.VerifyKey(); err == nil {
		err = c.core.DecryptBinaryFile(filenameIn, filenameOut) // error already logged by core
	}

	return err
}

func (c *CaesarCommand) Alphabet() string {
	return c.core.GetAlphabet()
}

/**
 * Checks the alphabet, if OK it is applied to the underlying cipher machine.
 * Else it logs an error and exits with ERR_BAD_ALPHABET.
 */
func (c *CaesarCommand) Rebuild(alphabet *cmn.Alphabet, opts ...any) {
	if alphabet.Check() {
		c.core.WithAlphabet(alphabet)
	} else {
		err := fmt.Errorf("invalid alphabet '%s' size:%d", alphabet.Name, alphabet.Size())
		mlog.ErrorE(err)
		app.DieWithError(err, caesarx.ERR_BAD_ALPHABET)
	}
}

func (c *CaesarCommand) String() string {
	return fmt.Sprintf("%s %s", c.core.GetLanguage(), c.core.String())
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/

func DemoCaesarCommand(alpha, numbers *cmn.Alphabet, phrase string) bool {
	fmt.Println("Caesar Encryption (Command-pattern version)")
	fmt.Println("Master : ", alpha.Name)
	if numbers != nil {
		fmt.Println("Slave  : ", numbers.Name)
	}
	var ok bool = true
	for _, key := range []rune{alpha.GetRuneAt(10), alpha.GetRuneAt(-1)} {
		var encTxt, encTxt2, decTxt string
		var err error

		ngram := cmn.NewNgramFormatter(5, '·')
		cnv1 := NewCaesarCommand(alpha, key).WithChain(numbers)
		cnv2 := NewCaesarCommand(alpha, key).WithChain(numbers)
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
		fmt.Println()

		if decTxt != phrase {
			ok = false
		}
	}

	return ok
}
