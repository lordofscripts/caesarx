/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   CaesarX
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Affine support added on 2025W39 out of that programmer's syndrome
 * of always wanting to add new features...
 *-----------------------------------------------------------------*/
package affine

import (
	"bufio"
	"fmt"
	"io"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/ciphers"
	"lordofscripts/caesarx/ciphers/caesar"
	"lordofscripts/caesarx/cmn"
	"lordofscripts/caesarx/internal/crypto"
	"os"
	"strings"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

//var _ ciphers.ICipher = (*AffineCrypto)(nil) // @todo pending...

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type AffineCrypto struct {
	langCode   string
	master     *affineContext
	slave      *affineContext
	sequencerE *crypto.AffineSequencer
	sequencerD *crypto.AffineSequencer
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

// Due to the nature of the Affine algorithm, this object holds two master
// RuneTranslator instances, one for Encoding, the other for Decoding. They
// are built before-hand so that we can cache the lookup tables rather than
// recalculating them for every item. This happens for both the master
// and slave.
type affineContext struct {
	E      *cmn.RuneTranslator // for Encoding
	D      *cmn.RuneTranslator // for Decoding
	params *AffineParams
}

/* ----------------------------------------------------------------
 *				I n i t i a l i z e r
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// (Ctor) creates a new instance of an Affine encryption/decryption
// engine. It assumes the given Affine coefficients are valid as this
// only verifies the N value matches the size of the alpha. If these
// are the result of NewAffineParams() then it should be okay.
//
//	It returns the new instance or nil if there was an error.
func NewAffineCrypto(alpha *cmn.Alphabet, params *AffineParams) *AffineCrypto {
	main := &affineContext{
		E:      nil,
		D:      nil,
		params: params,
	}

	const FOR_ENCODING bool = true
	const FOR_DECODING bool = false
	var rtE, rtD *cmn.RuneTranslator
	var err error

	rtE, err = buildTabula(alpha, params, true, FOR_ENCODING)
	if err != nil {
		return nil
	}
	main.E = rtE

	rtD, err = buildTabula(alpha, params, true, FOR_DECODING)
	if err != nil {
		return nil
	}
	main.D = rtD

	return &AffineCrypto{
		langCode:   alpha.LangCodeISO(),
		master:     main,
		slave:      nil,
		sequencerE: crypto.NewAffineSequencer(params.A, params.B, alpha),
		sequencerD: crypto.NewAffineSequencer(params.A, params.B, alpha),
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// String() implements the fmt.Stringer interface
func (c *AffineCrypto) String() string {
	return fmt.Sprintf("%s %s", crypto.ALG_CODE_AFFINE, c.langCode)
}

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *					G e n e r a l   P u r p o s e
 *- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -*/

// gets the pure alphabet as a string.
// Implements the ciphers.ICipher interface
func (c *AffineCrypto) GetAlphabet() string {
	return c.master.E.GetSource() // same as c.master.D.GetSource()
}

// Get a copy of the master & slave (if any) parameters. There is no purpose
// in modifying them because the returned values are cloned.
func (c *AffineCrypto) GetParams() (masterP *AffineParams, slaveP *AffineParams) {
	masterP = c.master.params.Clone()
	slaveP = nil
	if c.slave != nil {
		slaveP = c.slave.params.Clone()
	}

	return
}

// Attach a secondary (slave) alphabet chained to the master alphabet
// Implements the ciphers.ICipher interface
func (c *AffineCrypto) WithChain(alphaSlave *cmn.Alphabet) error {
	// request to remove chained alphabet
	if alphaSlave == nil {
		c.slave = nil
		return nil
	}

	// use the helper to (re)build coefficients for the chain/slave
	helper := NewAffineHelper()
	if err := helper.SetParams(c.master.params); err != nil {
		return err
	}

	// calculate the slave coefficients conditioned by the master's
	var slaveParams *AffineParams = nil
	var err error

	slaveN := int(alphaSlave.Size())
	if slaveN == c.master.params.N {
		// Case A
		// alphabet lengths are equal, we can use the same parameters
		slaveParams = c.master.params.Clone()
		mlog.Info("Affine slave A1=A2, N1=N2")
	} else { // differing alphabet lengths
		// Case B
		if helper.IsCommonCoprime(c.master.params.A, c.master.params.N, slaveN) {
			// common in both sets, no change either on A but on N
			slaveParams = c.master.params.Clone()
			slaveParams.N = slaveN
			mlog.Info("Affine slave A1=A2, N1!=N2")
		} else {
			// Case C
			// based on the master parameters but restrained to the slave's
			// condition, recalculate the A coefficient to apply to the SLAVE
			// the master remains with its own A coefficient.
			A := helper.CalculateSlaveCoprime(c.master.params, c.master.params.B, slaveN)
			if slaveParams, err = NewAffineParams(A, c.master.params.B, slaveN); err != nil {
				mlog.Error("could not set Affine slave due to error", err)
				return err
			}

			// not an error but needs to be logged
			mlog.InfoT("recalculated Affine A coefficient",
				mlog.Int("Master-A", c.master.params.A),
				mlog.Int("Master-N", c.master.params.N),
				mlog.Int("Slave-A", A),
				mlog.Int("Slave-N", slaveN))
		}
	}

	// now build the Tabulae Rectae for the Slave/Chain/Secondary
	const FOR_ENCODING bool = true
	const FOR_DECODING bool = false
	var rtSlaveE, rtSlaveD *cmn.RuneTranslator
	rtSlaveE, err = buildTabula(alphaSlave, slaveParams, false, FOR_ENCODING) // don't check N
	if err != nil {
		return err
	}
	rtSlaveD, err = buildTabula(alphaSlave, slaveParams, false, FOR_DECODING) // don't check N
	if err != nil {
		return err
	}

	// And finally store it in this instance as a properly configured chained alpha
	c.slave = &affineContext{
		E:      rtSlaveE,
		D:      rtSlaveD,
		params: slaveParams, // the slave coefficients we calculated above
	}

	return nil
}

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *					E n c r y p t i o n
 *- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -*/

// Encode takes a plain text (short) message and encrypts it using
// the current Affine coefficients using the current master & slave
// alphabets.
func (c *AffineCrypto) Encode(plain string) (string, error) {
	const EMPTY string = ""
	var cipher strings.Builder
	var encR rune
	var err error = nil

	for _, charP := range plain {
		encR = charP
		if !c.master.E.Exists(charP) {
			if c.slave != nil && c.slave.E.Exists(charP) {
				encR, err = c.slave.E.Lookup(charP)
			}
		} else {
			encR, err = c.master.E.Lookup(charP)
		}

		if err != nil {
			mlog.ErrorE(err)
			return EMPTY, err
		}
		if _, err = cipher.WriteRune(encR); err != nil {
			mlog.ErrorE(err)
			return EMPTY, err
		}
	}

	return cipher.String(), nil
}

// Encrypts a TEXT file using the current Affine coefficients and
// alphabets.
func (c *AffineCrypto) EncryptTextFile(input, output string) error {
	fdIn, err := os.Open(input)
	if err != nil {
		mlog.ErrorE(err)
	}
	defer fdIn.Close()
	reader := bufio.NewReader(fdIn)

	fdOut, err := os.Create(output)
	if err != nil {
		mlog.ErrorE(err)
	}
	defer fdOut.Close()

	destroyOpenFile := func(fd *os.File) {
		fd.Close() // in Windows a file must be closed prior to Remove...
		os.Remove(fd.Name())
	}

	var lineIn, lineOut string
	err = nil
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lineIn = scanner.Text()
		lineOut, err = c.Encode(lineIn)
		if err == nil {
			_, err = fmt.Fprintln(fdOut, lineOut)
		}

		if err != nil {
			mlog.ErrorE(err)
			destroyOpenFile(fdOut)
			return err
		}
	}

	if err = scanner.Err(); err != nil {
		mlog.ErrorE(err)
		destroyOpenFile(fdOut)
	}

	return err
}

// Encrypts a binary file and reports any error. If there was an error of
// any kind, the unfinished output file is deleted from the filesystem. (v1.1+)
func (c *AffineCrypto) EncryptBinaryFile(input, output string) error {
	// -- Preamble
	fdIn, err := os.Open(input)
	if err != nil {
		mlog.ErrorE(err)
	}
	defer fdIn.Close()

	fdOut, err := os.Create(output)
	if err != nil {
		mlog.ErrorE(err)
	}
	defer fdOut.Close()

	destroyOpenFile := func(fd *os.File) {
		fd.Close() // in Windows a file must be closed prior to Remove...
		os.Remove(fd.Name())
	}

	// -- Setup Cryptostream
	// @todo implement AffineSequencer
	master := ciphers.NewBinaryTabulaRecta()
	c.sequencerE.SetDecryptionMode(false) // only matters with Vigenere
	iter := caesar.NewBinaryIterator(c.sequencerE, master)
	defer c.sequencerE.Reset()

	// -- Process cryptostream
	const BUFFER_SIZE int = 4096
	buffer := make([]byte, BUFFER_SIZE)

	for {
		// (a) read bytes from input stream
		n, errR := fdIn.Read(buffer)
		if errR != nil {
			if errR == io.EOF {
				err = nil // successful termination of file
				break
			}

			// bad yu-yu
			err = fmt.Errorf("error reading binary file: %w", errR)
			break
		}

		// (b) GH-002 encode byte(s)
		iter.Start(buffer[:n])
		for !iter.EncodeNext() {
		}

		// (c) GH-002 write byte(s) to binary output file
		if writeCount, errW := fdOut.Write(iter.Result()); errW != nil {
			// oops! something happened with the filesystem
			err = errW
			break
		} else if writeCount != n {
			// mismatch between data buffer content size and written count
			err = fmt.Errorf("write count mismatch for binary file %d != %d", writeCount, n)
			break
		}
	}

	// -- Epilogue
	if err != nil {
		mlog.ErrorE(err)
		destroyOpenFile(fdOut)
	}

	return err
}

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *					D e c r y p t i o n
 *- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -*/

// Decode is a short message level Affine decryptor. It decrypts
// using the currently set Affine coefficients and alphabets.
// If decryption is successful it returns a nil error, else the
// returned data is empty and err contains the error.
func (c *AffineCrypto) Decode(cipher string) (string, error) {
	const EMPTY string = ""
	var plain strings.Builder
	var decR rune
	var err error = nil

	for _, charP := range cipher {
		decR = charP
		if !c.master.D.Exists(charP) {
			if c.slave != nil && c.slave.D.Exists(charP) {
				decR, err = c.slave.D.ReverseLookup(charP)
			}
		} else {
			decR, err = c.master.D.ReverseLookup(charP)
		}

		if err != nil {
			return EMPTY, err
		}
		if _, err = plain.WriteRune(decR); err != nil {
			return EMPTY, err
		}
	}

	return plain.String(), nil
}

// Decrypts a TEXT file using the Affine coefficient configuration
// and current alphabet(s). It does so by repeatedly calling the
// Decode() method for every line of the text file.
func (c *AffineCrypto) DecryptTextFile(input, output string) error {
	fdIn, err := os.Open(input)
	if err != nil {
		mlog.ErrorE(err)
	}
	defer fdIn.Close()
	reader := bufio.NewReader(fdIn)

	fdOut, err := os.Create(output)
	if err != nil {
		mlog.ErrorE(err)
	}
	defer fdOut.Close()

	destroyOpenFile := func(fd *os.File) {
		fd.Close() // in Windows a file must be closed prior to Remove...
		os.Remove(fd.Name())
	}

	var lineIn, lineOut string
	err = nil
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lineIn = scanner.Text()
		lineOut, err = c.Decode(lineIn)
		if err == nil {
			_, err = fmt.Fprintln(fdOut, lineOut)
		}

		if err != nil {
			mlog.ErrorE(err)
			destroyOpenFile(fdOut)
			return err
		}
	}

	if err = scanner.Err(); err != nil {
		mlog.ErrorE(err)
		destroyOpenFile(fdOut)
	}

	return err
}

// Encrypts a binary file and reports any error. If there was an error of
// any kind, the unfinished output file is deleted from the filesystem. (v1.1+)
func (c *AffineCrypto) DecryptBinaryFile(input, output string) error {
	// -- Preamble
	fdIn, err := os.Open(input)
	if err != nil {
		mlog.ErrorE(err)
	}
	defer fdIn.Close()

	fdOut, err := os.Create(output)
	if err != nil {
		mlog.ErrorE(err)
	}
	defer fdOut.Close()

	destroyOpenFile := func(fd *os.File) {
		fd.Close() // in Windows a file must be closed prior to Remove...
		os.Remove(fd.Name())
	}

	// -- Setup Cryptostream
	// @todo implement AffineSequencer
	master := ciphers.NewBinaryTabulaRecta()
	c.sequencerD.SetDecryptionMode(true) // only matters with Vigenere
	iter := caesar.NewBinaryIterator(c.sequencerD, master)
	defer c.sequencerD.Reset()

	// -- Process cryptostream
	const BUFFER_SIZE int = 4096
	buffer := make([]byte, BUFFER_SIZE)

	for {
		// (a) read bytes from input stream
		n, errR := fdIn.Read(buffer)
		if errR != nil {
			if errR == io.EOF {
				err = nil // successful termination of file
				break
			}

			// bad yu-yu
			err = fmt.Errorf("error reading binary file: %w", errR)
			break
		}

		// (b) GH-002 encode byte(s)
		iter.Start(buffer[:n])
		for !iter.DecodeNext() {
		}

		// (c) GH-002 write byte(s) to binary output file
		if writeCount, errW := fdOut.Write(iter.Result()); errW != nil {
			// oops! something happened with the filesystem
			err = errW
			break
		} else if writeCount != n {
			// mismatch between data buffer content size and written count
			err = fmt.Errorf("write count mismatch for binary file %d != %d", writeCount, n)
			break
		}
	}

	// -- Epilogue
	if err != nil {
		mlog.ErrorE(err)
		destroyOpenFile(fdOut)
	}

	return err
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s  (Encoder)
 *-----------------------------------------------------------------*/

// EncodeBytes achieves the same as Encode except it operates on a
// binary buffer when a Binary alphabet is chosen. Added in v1.1
// for binary file encryption.
/*
func (a *AffineEncoder) EncodeBytes(plain []byte) []byte {
	master := ciphers.NewBinaryTabulaRecta()
	a.sequencer.SetDecryptionMode(false) // only matters with Vigenere

	iter := NewBinaryIterator(a.sequencer, master) // @note no slaves with Binary!
	iter.Start(plain)
	for !iter.EncodeNext() {
	}

	a.sequencer.Reset()
	return iter.Result()
}
@TODO Binary for Affine */

// Encrypts a binary file and reports any error. If there was an error of
// any kind, the unfinished output file is deleted from the filesystem. (v1.1+)
/*
func (a *AffineEncoder) EncryptBinaryFile(input, output string) error {
	// -- Preamble
	fdIn, err := os.Open(input)
	if err != nil {
		mlog.ErrorE(err)
	}
	defer fdIn.Close()

	fdOut, err := os.Create(output)
	if err != nil {
		mlog.ErrorE(err)
	}
	defer fdOut.Close()

	destroyOpenFile := func(fd *os.File) {
		fd.Close() // in Windows a file must be closed prior to Remove...
		os.Remove(fd.Name())
	}

	// -- Setup Cryptostream
	// @todo implement AffineSequencer
	master := ciphers.NewBinaryTabulaRecta()
	a.sequencer.SetDecryptionMode(false) // only matters with Vigenere
	iter := NewBinaryIterator(a.sequencer, master)
	defer a.sequencer.Reset()

	// -- Process cryptostream
	const BUFFER_SIZE int = 4096
	buffer := make([]byte, BUFFER_SIZE)

	for {
		// (a) read bytes from input stream
		n, errR := fdIn.Read(buffer)
		if errR != nil {
			if errR == io.EOF {
				err = nil // successful termination of file
				break
			}

			// bad yu-yu
			err = fmt.Errorf("error reading binary file: %w", errR)
			break
		}

		// (b) GH-002 encode byte(s)
		iter.Start(buffer[:n])
		for !iter.EncodeNext() {
		}

		// (c) GH-002 write byte(s) to binary output file
		if writeCount, errW := fdOut.Write(iter.Result()); errW != nil {
			// oops! something happened with the filesystem
			err = errW
			break
		} else if writeCount != n {
			// mismatch between data buffer content size and written count
			err = fmt.Errorf("write count mismatch for binary file %d != %d", writeCount, n)
			break
		}
	}

	// -- Epilogue
	if err != nil {
		mlog.ErrorE(err)
		destroyOpenFile(fdOut)
	}

	return err
}
@TODO Binary for Affine */

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

// Builds a Tabula Recta that reflects the current Affine coefficients.
// If check is true, it checks that the size of the given alpha corresponds
// to the N value in params.
// The forEncoding parameter specifies whether the tabula must be built
// for encoding (uses params.A) or decoding (uses params.Ap).
func buildTabula(alpha *cmn.Alphabet, params *AffineParams, check bool, forEncoding bool) (*cmn.RuneTranslator, error) {
	var err error = nil
	if check && alpha.Size() != uint(params.N) {
		err = fmt.Errorf("bad params for '%s' AffineEncoder, mismatching N=%d", alpha.Name, params.N)
		mlog.ErrorE(err)
		return nil, err
	}

	helper := NewAffineHelper()
	if err = helper.SetParams(params); err != nil {
		return nil, err
	}

	// build transliterated primary alphabet based on chosen parameters
	// this way we don't have to repeat these relative expensive calculation
	// as we decode, sort of caching.
	var cipheredAlphabet strings.Builder
	for _, charP := range alpha.Chars {
		var altChar rune
		if forEncoding {
			altChar, err = helper.EncodeRuneFrom(charP, alpha.Chars)
		} else {
			altChar, err = helper.DecodeRuneFrom(charP, alpha.Chars)
		}

		if err != nil {
			mlog.ErrorE(err)
			return nil, err
		}
		cipheredAlphabet.WriteRune(altChar)
	}

	// ensure there is at least a default case handler
	caser := alpha.BorrowSpecialCase()
	if caser == nil {
		caser = cmn.DefaultCaseHandler
	}

	// create the translator
	var rt *cmn.RuneTranslator
	if forEncoding {
		// al derecho
		rt = cmn.NewSimpleRuneTranslator(alpha.Name, alpha.Chars, cipheredAlphabet.String(), caser)
	} else {
		// y al revéz...
		rt = cmn.NewSimpleRuneTranslator(alpha.Name, cipheredAlphabet.String(), alpha.Chars, caser)
	}

	return rt, nil
}
