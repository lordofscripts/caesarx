/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * (Command Pattern - See "Design Patterns")
 * The Affine is a simple linear cipher f(x) = (A*x+B)%N. In fact the
 * plain Caesar cipher is a simplification of Affine where A=1 and
 * B is key shift (Caesar key).
 *-----------------------------------------------------------------*/
package commands

import (
	"fmt"
	"lordofscripts/caesarx"
	"lordofscripts/caesarx/app"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/ciphers"
	"lordofscripts/caesarx/ciphers/affine"
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

var _ ciphers.IPipe = (*AffineCommand)(nil)
var _ ciphers.ICipherCommand = (*AffineCommand)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type AffineCommand struct {
	ciphers.Pipe
	crypto *affine.AffineCrypto
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewAffineCommand(alpha *cmn.Alphabet, a, b int) *AffineCommand {
	params, err := affine.NewAffineParams(a, b, int(alpha.Size()))
	if err != nil {
		return nil // already logged by AffineParams
	}

	return NewAffineCommandExt(alpha, params)
}

func NewAffineCommandExt(alpha *cmn.Alphabet, params *affine.AffineParams) *AffineCommand {
	eng := affine.NewAffineCrypto(alpha, params)

	return &AffineCommand{
		Pipe:   ciphers.NewEmptyPipe(),
		crypto: eng,
	}
}

/* ----------------------------------------------------------------
 *							M e t h o d s (ICipherCommand)
 *-----------------------------------------------------------------*/

/**
 * Same as Rebuild() for this simple cipher.
 */
func (c *AffineCommand) WithAlphabet(alphabet *cmn.Alphabet) ciphers.ICipherCommand {
	c.Rebuild(alphabet) //@audit TODO
	return c
}

/**
 * Chain a slave disk/tabula that would use the same (or corrected) alphabet shift
 * as the main alphabet. A slave disk/tabula is usually a Numbers and/or Symbols
 * to IMPROVE the ancient cipher against attacks.
 */
func (c *AffineCommand) WithChain(slave *cmn.Alphabet) ciphers.ICipherCommand {
	if slave != nil {
		if err := c.crypto.AffineEncoder.WithChain(slave); err != nil {
			mlog.Fatal(caesarx.ERR_INTERNAL, err)
		}
		if err := c.crypto.AffineDecoder.WithChain(slave); err != nil {
			mlog.Fatal(caesarx.ERR_INTERNAL, err)
		}
	} else {
		c.crypto.AffineEncoder.WithChain(nil)
		c.crypto.AffineDecoder.WithChain(nil)
	}

	return c
}

func (c *AffineCommand) Encode(plain string) (string, error) {
	ciphered, err := c.crypto.Encode(plain)
	if err != nil {
		return "", err
	}

	if c.IsPipeOpen() {
		return c.PipeOutput(ciphers.PipeEncode, ciphered)
	} else {
		return ciphered, nil
	}
}

func (c *AffineCommand) Decode(ciphered string) (string, error) {
	plain, err := c.crypto.Decode(ciphered)
	if err != nil {
		return "", err
	}

	if c.IsPipeOpen() {
		return c.PipeOutput(ciphers.PipeDecode, plain)
	} else {
		return plain, nil
	}
}

func (c *AffineCommand) Alphabet() string {
	return c.crypto.Alphabet()
}

func (c *AffineCommand) GetParams() (masterP *affine.AffineParams, slaveP *affine.AffineParams) {
	return c.crypto.GetParams()
}

/**
 * N.A.
 */
func (c *AffineCommand) Rebuild(alphabet *cmn.Alphabet, opts ...any) { //@audit TODO rewrite this hack
	/*
		var newModulo int = int(alphabet.Size())
		paramOld := c.crypto.AffineEncoder.
		cA, _, cB, oldModulo := c.helper.GetParameters()

		if oldModulo != newModulo {
			if err := c.helper.SetParameters(cA, cB, newModulo); err != nil {
				mlog.FatalT(caesarx.ERR_BAD_ALPHABET,
					"new alphabet incompatible with current Affine coeficcients",
					mlog.String("OldAlpha", c.alpha.Name),
					mlog.String("NewAlpha", alphabet.Name),
					mlog.Int("Old-N", oldModulo),
					mlog.Int("New-N", newModulo),
					mlog.String("Error", err.Error()))
			}
		}

		c.alpha = alphabet
	*/
}

func (c *AffineCommand) String() string {
	return c.crypto.String()
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/

func DemoAffineCommand(alpha, numbers *cmn.Alphabet, phrase string) bool {
	fmt.Println("Affine Encryption (Command-pattern version)")
	fmt.Println("Master : ", alpha.Name)
	if numbers != nil {
		fmt.Println("Slave  : ", numbers.Name)
	}

	const A int = 7 // 7 is a coprime that is common to all my built-in alphabets
	const B int = 12
	var ok bool = false
	if params, err := affine.NewAffineParams(A, B, int(alpha.Size())); err != nil {
		mlog.Fatal(caesarx.ERR_DEMO_ERROR, err.Error())
	} else {
		var err error
		var encTxt, decTxt string

		alg := NewAffineCommandExt(alpha, params)
		encTxt, err = alg.Encode(phrase)
		if err != nil {
			app.DieWithError(err, caesarx.ERR_DEMO_ERROR)
		}
		decTxt, err = alg.Decode(encTxt)
		if err != nil {
			app.DieWithError(err, caesarx.ERR_DEMO_ERROR)
		}

		ok = decTxt == phrase

		fmt.Printf("Params : %s\n", params.String())
		fmt.Println("Plain  : ", phrase)
		fmt.Println("Encoded: ", encTxt)
		fmt.Println("Decoded: ", decTxt)
		fmt.Println()
	}

	return ok
}
