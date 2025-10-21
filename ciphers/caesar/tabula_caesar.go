/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Plain Caesar cipher using Tabula Recta implementation. Is
 * case-insensitive but preserves case.
 *-----------------------------------------------------------------*/
package caesar

import (
	"bufio"
	"fmt"
	"io"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/ciphers"
	"lordofscripts/caesarx/cmn"
	"lordofscripts/caesarx/internal/crypto"
	"os"
	"sync"
	"unicode/utf8"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/
var (
	Info = ciphers.NewCipherInfo(crypto.ALG_CODE_CAESAR, "1.0",
		"Julius Caesar",
		crypto.ALG_NAME_CAESAR,
		"Caesar cipher")
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/
func init() {
	ciphers.RegisterCipher(Info)
}

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ ciphers.ICipher = (*CaesarTabulaRecta)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type CaesarTabulaRecta struct {
	alpha     *cmn.Alphabet
	slave     *ciphers.TabulaRecta // implements cmn.IRuneLocalizer
	sequencer crypto.IKeySequencer
	mu        *sync.Mutex
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/**
 * (Ctor) Caesar Cipher using a Tabula Recta that supports ASCII and
 * foreign (UTF8) alphabets.
 * · Always follow it with a call to VerifyKey() or VerifySecret() prior to
 *	 beginning encoding/decoding.
 * · follow with WithChain() to chain with supplemental alphabets.
 * · follow with WithAlphabet() to specify a different alphabet prior to encoding.
 * · It does case-folding by default, so it handles & preserves upper/lowercase
 */
func NewCaesarTabulaRecta(alphabet *cmn.Alphabet, key rune) *CaesarTabulaRecta {
	return &CaesarTabulaRecta{
		alpha:     alphabet,
		slave:     nil,
		sequencer: crypto.NewCaesarSequencer(key),
		mu:        new(sync.Mutex),
	}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer
func (cx *CaesarTabulaRecta) String() string {
	return cx.sequencer.GetKeyInfo()
}

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *					G e n e r a l   P u r p o s e
 *- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -*/

// WithChain() reconfigures the current instance to attach a chained (secondary/slave)
// alphabet disk with supplementary characters not present in the main (primary/master).
func (cx *CaesarTabulaRecta) WithChain(extra *ciphers.TabulaRecta) ciphers.ICipher {
	cx.mu.Lock()
	defer cx.mu.Unlock()

	if !cx.alpha.IsBinary() {
		cx.slave = extra // v1.1 When primary is Binary no slaves are allowed
	} else {
		mlog.WarnT("ignored chained slave because primary alphabet is Binary", mlog.At())
	}
	return cx
}

// WithAlphabet() reconfigures the current instance to replace the MAIN
// (primary) alphabet.
func (cx *CaesarTabulaRecta) WithAlphabet(alphabet *cmn.Alphabet) ciphers.ICipher {
	cx.mu.Lock()
	defer cx.mu.Unlock()

	cx.alpha = alphabet
	return cx
}

// WithSequencer() specifies a Key Sequencer for the current instance.
func (cx *CaesarTabulaRecta) WithSequencer(keygen crypto.IKeySequencer) ciphers.ICipher {
	cx.mu.Lock()
	defer cx.mu.Unlock()

	cx.sequencer = keygen
	return cx
}

/**
 * Verify key(s). If none given it checks the key given in the constructor,
 * else it checks all the given keys. The key (single character) must be
 * present in the encoding alphabet.
 */
func (cx *CaesarTabulaRecta) VerifyKey(keys ...rune) error {
	verify := func(k rune) error {
		if !cx.alpha.Contains(k, cmn.CaseInsensitive) { // @audit what if TR is not case folded!
			return fmt.Errorf("key '%c' is not part of the alphabet", k)
		}
		return nil
	}

	if len(keys) == 0 {
		//return verify(cx.sequencer.GetKey(DUMMY_POS, DUMMY_RUNE))
		return cx.sequencer.Verify(verify)
	} else {
		for i, key := range keys {
			err := verify(key)
			if err != nil {
				return fmt.Errorf("key #%d (%c): %w", i+1, key, err)
			}
		}
	}

	return nil
}

/**
 * Does the same as VerifyKey() except it checks that all the given keys
 * (if any) have exactly ONE character (could be multi-byte Unicode character).
 */
func (cx *CaesarTabulaRecta) VerifySecret(secrets ...string) error {
	if len(secrets) == 0 {
		return cx.VerifyKey()
	}

	allKeys := make([]rune, len(secrets))
	for i, keyStr := range secrets {
		if utf8.RuneCountInString(keyStr) != 1 {
			return fmt.Errorf("key #%d '%s' contains more than one char", i+1, keyStr)
		}

		allKeys[i] = []rune(keyStr)[0]
	}

	return cx.VerifyKey(allKeys...)
}

/**
 * (IRuneLocalizer) Find a rune in the object's alphabet catalog.
 * Rune not found: error set, other return values nil or -1.
 * Rune found: error nil, pointer to alphabet and position within.
 * This method seeks in both the Primary/Master & Secondary/Slave.
 */
func (t *CaesarTabulaRecta) FindRune(r rune) (alpha string, at int, err error) {
	err = nil
	//at = cmn.RuneIndex(t.alpha.Chars, r)
	at = cmn.RuneIndexFold(t.alpha.Chars, r, t.alpha.BorrowSpecialCase())
	if at == -1 {
		// not present in the Primary/Master alphabet, let's try the Slave
		if t.slave == nil {
			// bummer! there is no Secondary/Slave, ran out of options
			err = fmt.Errorf("info: '%c' absent in %s", r, t.alpha.Name)
			alpha = ""
		} else {
			// perhaps in the Slave (or not) and that's final
			alpha, at, err = t.slave.FindRune(r)
		}
	} else {
		alpha = t.alpha.Chars
	}

	return
}

// GetAlphabet returns the contents of the alphabet.
func (cx *CaesarTabulaRecta) GetAlphabet() string {
	return cx.alpha.Chars
}

// GetLanguage returns the two-letter ISO code of the
// current alphabet's language. CaesarX supports several built-in
// language sets such as English, Spanish, German, Greek & Cyrillic.
func (cx *CaesarTabulaRecta) GetLanguage() string {
	return cx.alpha.Name
}

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *					E n c r y p t i o n
 *- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -*/

// Encode uses the appropriate substitution sequencer to encode
// a string using a modern version of the Caesar-class algorithms.
func (cx *CaesarTabulaRecta) Encode(plain string) string {
	defer cx.sequencer.Reset()

	master := ciphers.NewTabulaRecta(cx.alpha, cmn.CaseInsensitive)
	cx.sequencer.SetDecryptionMode(false) // only matters with Vigenere
	iter := NewTextIterator(cx.sequencer, master, cx.slave)
	iter.Start(plain)
	for !iter.EncodeNext() {
		//fmt.Print("E")
	}
	//fmt.Println()

	return iter.Result()
}

// EncodeBytes achieves the same as Encode except it operates on a
// binary buffer when a Binary alphabet is chosen. Added in v1.1
// for binary file encryption.
func (cx *CaesarTabulaRecta) EncodeBytes(plain []byte) []byte {
	defer cx.sequencer.Reset()

	master := ciphers.NewBinaryTabulaRecta()
	cx.sequencer.SetDecryptionMode(false) // only matters with Vigenere

	iter := NewBinaryIterator(cx.sequencer, master) // @note no slaves with Binary!
	iter.Start(plain)
	for !iter.EncodeNext() {
	}

	return iter.Result()
}

// Encrypts the input TEXT file using the selected Caesar variant and
// produces the output filename with the encrypted contents.
func (cx *CaesarTabulaRecta) EncryptTextFile(input, output string) error {
	cx.mu.Lock()
	defer cx.mu.Unlock()

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

	master := ciphers.NewTabulaRecta(cx.alpha, cmn.CaseInsensitive)
	cx.sequencer.SetDecryptionMode(false) // only matters with Vigenere
	iter := NewTextIterator(cx.sequencer, master, cx.slave)
	defer cx.sequencer.Reset()

	var lineIn string
	err = nil
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lineIn = scanner.Text()
		iter.Start(lineIn)
		for !iter.EncodeNext() {
		}
		if _, err = fmt.Fprintln(fdOut, iter.Result()); err != nil {
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
func (cx *CaesarTabulaRecta) EncryptBinaryFile(input, output string) error {
	cx.mu.Lock()
	defer cx.mu.Unlock()

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
	master := ciphers.NewBinaryTabulaRecta()
	cx.sequencer.SetDecryptionMode(false) // only matters with Vigenere
	iter := NewBinaryIterator(cx.sequencer, master)
	defer cx.sequencer.Reset()

	// -- Process cryptostream
	const BUFFER_SIZE int = 4096
	buffer := make([]byte, BUFFER_SIZE)
	firstCall := true

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
		if firstCall {
			iter.Start(buffer[:n])
			firstCall = false
		} else {
			iter.Update(buffer[:n])
		}

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

// Decode uses the appropriate substitution sequencer to decode
// a string using a modern version of the Caesar-class algorithms.
func (cx *CaesarTabulaRecta) Decode(ciphered string) string {
	defer cx.sequencer.Reset()

	master := ciphers.NewTabulaRecta(cx.alpha, cmn.CaseInsensitive)
	cx.sequencer.SetDecryptionMode(true) // only matters with Vigenere

	iter := NewTextIterator(cx.sequencer, master, cx.slave)
	iter.Start(ciphered)
	for !iter.DecodeNext() {
		//fmt.Print("D")
	}
	//fmt.Println()

	return iter.Result()
}

// DecodeBytes achieves the same as Decode except it operates on a
// binary buffer when a Binary alphabet is chosen. Added in v1.1
// for binary file decryption.
func (cx *CaesarTabulaRecta) DecodeBytes(ciphered []byte) []byte {
	defer cx.sequencer.Reset()

	master := ciphers.NewBinaryTabulaRecta()
	cx.sequencer.SetDecryptionMode(true) // only matters with Vigenere

	iter := NewBinaryIterator(cx.sequencer, master) // @note no slaves with Binary!
	iter.Start(ciphered)
	for !iter.DecodeNext() {
	}

	return iter.Result()
}

// Decrypts the input TEXT file using the selected Caesar variant and
// produces the output filename with the decrypted contents.
func (cx *CaesarTabulaRecta) DecryptTextFile(input, output string) error {
	cx.mu.Lock()
	defer cx.mu.Unlock()

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

	master := ciphers.NewTabulaRecta(cx.alpha, cmn.CaseInsensitive)
	cx.sequencer.SetDecryptionMode(true) // only matters with Vigenere
	iter := NewTextIterator(cx.sequencer, master, cx.slave)
	defer cx.sequencer.Reset()

	var lineIn string
	err = nil
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lineIn = scanner.Text()
		iter.Start(lineIn)
		for !iter.DecodeNext() {
		}
		if _, err = fmt.Fprintln(fdOut, iter.Result()); err != nil {
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

// Decrypts a binary file and reports any error. If there was an error of
// any kind, the unfinished output file is deleted from the filesystem. (v1.1+)
func (cx *CaesarTabulaRecta) DecryptBinaryFile(input, output string) error {
	cx.mu.Lock()
	defer cx.mu.Unlock()

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
	master := ciphers.NewBinaryTabulaRecta()
	cx.sequencer.SetDecryptionMode(true) // only matters with Vigenere
	iter := NewBinaryIterator(cx.sequencer, master)
	defer cx.sequencer.Reset()

	// -- Process cryptostream
	const BUFFER_SIZE int = 4096
	buffer := make([]byte, BUFFER_SIZE)
	firstCall := true

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
		if firstCall {
			iter.Start(buffer[:n])
			firstCall = false
		} else {
			iter.Update(buffer[:n])
		}

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
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/

func DemoCaesarPlain(alpha *cmn.Alphabet, message string) {
	var keyNoOp = alpha.GetRuneAt(0)
	var keyMidl = alpha.GetRuneAt(int(alpha.Size() / 2))
	var keyLast = alpha.GetRuneAt(-1)
	var multiByte = utf8.RuneCountInString(alpha.Chars) > len(alpha.Chars)

	doRoundTrip := func(key rune) {
		cnv := NewCaesarTabulaRecta(alpha, key)
		cipher := cnv.Encode(message)
		plain := cnv.Decode(cipher)
		passed := plain == message

		fmt.Printf("\tKey       : %c\n", key)
		fmt.Printf("\tMessage   : %s\n", message)
		fmt.Printf("\tCiphered  : %s\n", cipher)
		fmt.Printf("\tDeciphered: %s\n", plain)
		fmt.Printf("\tPassed    : %t\n\n", passed)
	}

	fmt.Println("Plain Caesar Demo with Tabula Recta")
	fmt.Printf("\tAlphabet  : %s\n", alpha.Name)
	fmt.Printf("\tMulti-byte: %t\n", multiByte)

	doRoundTrip(keyNoOp)
	doRoundTrip(keyMidl)
	doRoundTrip(keyLast)
}
