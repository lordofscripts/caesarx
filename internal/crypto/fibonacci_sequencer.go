/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * The Key Sequencer is a flexible generator of the current key
 * during Encoding/Decoding.
 * Fibonacci sequence: 0 1 1 2 3 5 8 13 21 34
 *-----------------------------------------------------------------*/
package crypto

import (
	"fmt"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/cmn"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const ALG_NAME_FIBONACCI = "Fibonacci"
const ALG_CODE_FIBONACCI = "FIBO"

var (
	fibonacci []int = []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ IKeySequencer = (*FibonacciSequencer)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type FibonacciSequencer struct {
	skipped int
	current int
	fkeys   []rune
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/**
 * (Ctor) A Fibonacci key sequencer based on the number of terms
 * suitable for the given alphabet. This object keeps no reference
 * to the passed alphabet, it just picks the necessary keys.
 */
func NewFibonacciSequencer(alpha *cmn.Alphabet, primeKey rune) *FibonacciSequencer {
	fkeys := make([]rune, 0)
	maxA := alpha.Size()
	// The prime key serves as a pivot to introduce variability
	pivot := alpha.PositionOf(primeKey)
	for _, offset := range fibonacci {
		value := (pivot + offset) % int(maxA)
		if value == 0 {
			value = 1
		}
		fkeys = append(fkeys, alpha.GetRuneAt(value))
	}

	mlog.TraceT("FibonacciSequencer", mlog.Rune("Prime", primeKey), mlog.Int("Factors", len(fibonacci)))
	return &FibonacciSequencer{0, 0, fkeys}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

/**
 * @returns The sequencer's friendly name.
 */
func (fs *FibonacciSequencer) Name() string {
	return ALG_NAME_FIBONACCI
}

/**
 * To instruct the sequencer whether it is being used for encipherment
 * or decipherment. Not relevant with Fibonacci variant.
 */
func (cs *FibonacciSequencer) SetDecryptionMode(isDecrypting bool) {
}

/**
 * N.A.
 */
func (cs *FibonacciSequencer) Feedback(rune) error {
	return nil
}

/**
 * Skip the current position. Must be called when character in the
 * input stream is not part of the encoding alphabet. This is important
 * to maintain integrity if non-encoding runes are removed from the
 * message at some point.
 *
 * @returns (int) number of skipped runes so far.
 */
func (fs *FibonacciSequencer) Skip() int {
	fs.skipped++
	return fs.skipped
}

/**
 * Get the key to be used for encoding target rune at this position.
 * Should only be called if target is part of the encoding alphabet!
 * NOTE: Fibonacci variant uses a derived Caesar key that repeats
 *		 over the input. It uses a 10-term Fibonacci series as shift
 *		 over the Prime Key.
 *
 * @param pos (int) ignored in this algorithm
 * @param target (rune) ignored in this algorithm
 * @returns the basic key to use for encoding/decoding at this position.
 */
func (fs *FibonacciSequencer) GetKey(pos int, target rune) rune {
	useKey := fs.fkeys[fs.current]
	fs.current++
	if fs.current >= len(fs.fkeys) {
		fs.current = 0
	}

	return useKey
}

func (fs *FibonacciSequencer) GetKeyInfo() string {
	return fmt.Sprintf("%cƒ𝓍 (%d)", UC_MATH_BOLD_F, len(fs.fkeys))
}

func (cs *FibonacciSequencer) Verify(callback func(rune) error) error {
	for _, k := range cs.fkeys {
		if err := callback(k); err != nil {
			return err
		}
	}
	return nil
}

/**
 * Resets the sequencer. It should be done after every Encode or Decode
 */
func (cs *FibonacciSequencer) Reset() {
	cs.skipped = 0
	cs.current = 0
}
