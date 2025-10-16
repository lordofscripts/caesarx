/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * The Key Sequencer is a flexible generator of the current key
 * during Encoding/Decoding.
 *-----------------------------------------------------------------*/
package ciphers

import "fmt"

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const ALG_NAME_CAESAR = "Caesar"
const ALG_CODE_CAESAR = "CAES"

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ IKeySequencer = (*CaesarSequencer)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type CaesarSequencer struct {
	prime   rune
	skipped int
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewCaesarSequencer(key rune) *CaesarSequencer {
	return &CaesarSequencer{key, 0}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

/**
 * @returns The sequencer's friendly name.
 */
func (cs *CaesarSequencer) Name() string {
	return ALG_NAME_CAESAR
}

/**
 * To instruct the sequencer whether it is being used for encipherment
 * or decipherment. Not relevant with standard Caesar.
 */
func (cs *CaesarSequencer) SetDecryptionMode(isDecrypting bool) {
}

/**
 * N.A.
 */
func (cs *CaesarSequencer) Feedback(rune) error {
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
func (cs *CaesarSequencer) Skip() int {
	cs.skipped++
	return cs.skipped
}

/**
 * Get the key to be used for encoding target rune at this position.
 * Should only be called if target is part of the encoding alphabet!
 * NOTE: Caesar algorithm always uses the same single key for all
 *	valid characters in the input stream.
 *
 * @param pos (int) ignored in this algorithm
 * @param target (rune) ignored in this algorithm
 * @returns the basic key to use for encoding/decoding at this position.
 */
func (cs *CaesarSequencer) GetKey(pos int, target rune) rune {
	return cs.prime
}

func (cs *CaesarSequencer) GetKeyInfo() string {
	return fmt.Sprintf("ƒ𝓍%c (%s)", UC_MATH_BOLD_C, string(cs.prime))
}

func (cs *CaesarSequencer) Verify(callback func(rune) error) error {
	return callback(cs.prime)
}

/**
 * Resets the sequencer. It should be done after every Encode or Decode
 */
func (cs *CaesarSequencer) Reset() {}
