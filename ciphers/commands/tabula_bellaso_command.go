/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * (Command Pattern - See "Design Patterns")
 * The cipher created by Giovanni Bellaso based on the Caesar cipher.
 * However, Bellaso is polialphabetic because it does not rely on a
 * single-character key, instead it uses a "secret" word or phrase
 * that is repeated over the input message, but ONLY over the characters
 * that are present in the primary/master alphabet.
 *  This cipher is often misattributed to Blas de Vigenère. Vigenère
 * created the auto-key variant based on Bellaso's work.
 *-----------------------------------------------------------------*/
package commands

import (
	"fmt"
	"lordofscripts/caesarx"
	"lordofscripts/caesarx/app"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/ciphers"
	"lordofscripts/caesarx/ciphers/bellaso"
	"lordofscripts/caesarx/cmn"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	// Filename extension for files encrypted with Bellaso
	FILE_EXT_BELLASO string = ".bel"
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ ciphers.IPipe = (*BellasoCommand)(nil)
var _ ciphers.ICipherCommand = (*BellasoCommand)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type BellasoCommand struct {
	ciphers.Pipe
	core        *bellaso.BellasoTabulaRecta
	outFilename string
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewBellasoCommand(alpha *cmn.Alphabet, secret string) *BellasoCommand {
	return &BellasoCommand{
		Pipe:        ciphers.NewEmptyPipe(),
		core:        bellaso.NewBellasoTabulaRecta(alpha, secret),
		outFilename: "",
	}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer
func (c *BellasoCommand) String() string {
	return c.core.String()
}

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *					G e n e r a l   P u r p o s e
 *- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -*/

/**
 * Same as Rebuild() for this simple cipher.
 */
func (c *BellasoCommand) WithAlphabet(alphabet *cmn.Alphabet) ciphers.ICipherCommand {
	c.Rebuild(alphabet)
	return c
}

/**
 * Chain a slave disk/tabula that would use the same (or corrected) alphabet shift
 * as the main alphabet. A slave disk/tabula is usually a Numbers and/or Symbols
 * to IMPROVE the ancient cipher against attacks.
 */
func (c *BellasoCommand) WithChain(slave *cmn.Alphabet) ciphers.ICipherCommand {
	if slave != nil {
		slaveTR := ciphers.NewTabulaRecta(slave, true)
		c.core.WithChain(slaveTR)
	} else {
		c.core.WithChain(nil)
	}

	return c
}

// this result is only meaningful after EncryptBinFile() or EncryptTextFile()
// where the output filename is not explicitely given but generated.
func (c *BellasoCommand) GetOutputFilename() string {
	return c.outFilename
}

// get the current alphabet's string
func (c *BellasoCommand) Alphabet() string {
	return c.core.GetAlphabet()
}

// Checks the alphabet, if OK it is applied to the underlying cipher machine.
// Else it logs an error and exits with ERR_BAD_ALPHABET.
func (c *BellasoCommand) Rebuild(alphabet *cmn.Alphabet, opts ...any) {
	if alphabet.Check() {
		c.core.WithAlphabet(alphabet)
	} else {
		err := fmt.Errorf("invalid alphabet '%s' size:%d", alphabet.Name, alphabet.Size())
		mlog.ErrorE(err)
		app.DieWithError(err, caesarx.ERR_BAD_ALPHABET)
	}
}

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *					E n c r y p t i o n
 *- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -*/

func (c *BellasoCommand) Encode(plain string) (string, error) {
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

// EncryptTextFile encrypts the filename src using the standard Caesar cipher.
// The output file has the FILE_EXT_BELLASO file extension. Please note that
// this method is only for text files.
func (c *BellasoCommand) EncryptTextFile(src string) error {
	var err error = nil
	if err = c.core.VerifyKey(); err == nil {
		fileOut := cmn.NewNameExtOnly(src, FILE_EXT_BELLASO, true)
		err = c.core.EncryptTextFile(src, fileOut) // error already logged by core
		if err == nil {
			c.outFilename = fileOut
		}
	}

	return err
}

// Encodes a binary file and produces a binary encoded file (v1.1+)
func (c *BellasoCommand) EncryptBinFile(filenameIn string) error {
	var err error = nil
	if err = c.core.VerifyKey(); err == nil {
		fileOut := cmn.NewNameExtOnly(filenameIn, FILE_EXT_BELLASO, true)
		err = c.core.EncryptBinaryFile(filenameIn, fileOut) // error already logged by core
		if err == nil {
			c.outFilename = fileOut
		}
	}

	return err
}

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *					D e c r y p t i o n
 *- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -*/

// decode a short message encrypted with Bellaso cipher
func (c *BellasoCommand) Decode(ciphered string) (string, error) {
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

// DecryptTextFile decrypts the filename src using the standard Caesar cipher.
// The output file target must be explicitely given. Please note that
// this method is only for text files.
func (c *BellasoCommand) DecryptTextFile(src, target string) error {
	var err error = nil
	if err = c.core.VerifyKey(); err == nil {
		err = c.core.DecryptTextFile(src, target) // error already logged by core
	}

	return err
}

// Decodes a binary file and produces a plain binary file (v1.1+)
func (c *BellasoCommand) DecryptBinFile(filenameIn, filenameOut string) error {
	var err error = nil
	if err = c.core.VerifyKey(); err == nil {
		err = c.core.DecryptBinaryFile(filenameIn, filenameOut) // error already logged by core
	}

	return err
}

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/

func DemoBellasoCommand(alpha, numeric *cmn.Alphabet, phrase string) bool {
	fmt.Println("Bellaso Encryption (Command-pattern version)")
	fmt.Println("( A polyalphabetic Caesar variant )")

	var Secret1 string = fmt.Sprintf("%c%c%c", alpha.GetRuneAt(10),
		alpha.GetRuneAt(5),
		alpha.GetRuneAt(-1))
	var Secret2 string = fmt.Sprintf("%c%c%c%c", alpha.GetRuneAt(2),
		alpha.GetRuneAt(8),
		alpha.GetRuneAt(20),
		alpha.GetRuneAt(-1))
	var ok bool = true
	for _, secret := range []string{Secret1, Secret2} {
		var encTxt, encTxt2, decTxt string
		var err error

		ngram := cmn.NewNgramFormatter(5, '·')
		cnv1 := NewBellasoCommand(alpha, secret)
		cnv2 := NewBellasoCommand(alpha, secret)
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

		fmt.Printf("Secret : %s\n", secret)
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
