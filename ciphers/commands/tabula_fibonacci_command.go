/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * (Command Pattern - See "Design Patterns")
 * A variation of Caesar that encodes even positions with a key that
 * varies according to a Fibonacci series.
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

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ ciphers.IPipe = (*FibonacciCommand)(nil)
var _ ciphers.ICipherCommand = (*FibonacciCommand)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type FibonacciCommand struct {
	ciphers.Pipe
	core *caesar.FibonacciTabulaRecta
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewFibonacciCommand(alpha *cmn.Alphabet, primeKey rune) *FibonacciCommand {
	return &FibonacciCommand{
		Pipe: ciphers.NewEmptyPipe(),
		core: caesar.NewFibonacciTabulaRecta(alpha, primeKey),
	}
}

/* ----------------------------------------------------------------
 *							M e t h o d s (ICipherCommand)
 *-----------------------------------------------------------------*/

/**
 * Same as Rebuild() for this simple cipher.
 */
func (c *FibonacciCommand) WithAlphabet(alphabet *cmn.Alphabet) ciphers.ICipherCommand {
	c.Rebuild(alphabet)
	return c
}

/**
 * Chain a slave disk/tabula that would use the same (or corrected) alphabet shift
 * as the main alphabet. A slave disk/tabula is usually a Numbers and/or Symbols
 * to IMPROVE the ancient cipher against attacks.
 */
func (c *FibonacciCommand) WithChain(slave *cmn.Alphabet) ciphers.ICipherCommand {
	if slave != nil {
		slaveTR := ciphers.NewTabulaRecta(slave, true)
		c.core.WithChain(slaveTR)
	} else {
		c.core.WithChain(nil)
	}

	return c
}

func (c *FibonacciCommand) Encode(plain string) (string, error) {
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

func (c *FibonacciCommand) Decode(ciphered string) (string, error) {
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

func (c *FibonacciCommand) Alphabet() string {
	return c.core.GetAlphabet()
}

/**
 * Checks the alphabet, if OK it is applied to the underlying cipher machine.
 * Else it logs an error and exits with ERR_BAD_ALPHABET.
 */
func (c *FibonacciCommand) Rebuild(alphabet *cmn.Alphabet, opts ...any) {
	if alphabet.Check() {
		c.core.WithAlphabet(alphabet)
	} else {
		err := fmt.Errorf("invalid alphabet '%s' size:%d", alphabet.Name, alphabet.Size())
		mlog.ErrorE(err)
		app.DieWithError(err, caesarx.ERR_BAD_ALPHABET)
	}
}

func (c *FibonacciCommand) String() string {
	return c.core.String()
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/

func DemoFibonacciCommand(alpha, numeric *cmn.Alphabet, phrase string) bool {
	fmt.Println("Fibonacci Encryption (Command-pattern version)")
	fmt.Println("( A variant of the Caesar family )")
	var ok bool = true
	for _, key := range []rune{alpha.GetRuneAt(10), alpha.GetRuneAt(-1)} {
		var encTxt, encTxt2, decTxt string
		var err error

		ngram := cmn.NewNgramFormatter(5, '·')
		cnv1 := NewFibonacciCommand(alpha, key)
		cnv2 := NewFibonacciCommand(alpha, key)
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
