/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * The Key Sequencer is a flexible generator of the current key
 * during Encoding/Decoding.
 *-----------------------------------------------------------------*/
package ciphers

import (
	"fmt"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/cmn"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const ALG_NAME_DIDIMUS = "Didimus"
const ALG_CODE_DIDIMUS = "DIDI"

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ IKeySequencer = (*DidimusSequencer)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type DidimusSequencer struct {
	prime   rune // the prime key for even convertable positions (0,2,4...)
	alt     rune // alternate key for odd convertable positions (1,3,5...)
	skipped int
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewDidimusSequencer(primeKey rune, offset uint8, alpha *cmn.Alphabet) *DidimusSequencer {
	alphaLen := alpha.Size()
	keyOrd := alpha.PositionOf(primeKey)
	altOrd := (keyOrd + int(offset)) % int(alphaLen)
	if altOrd == 0 {
		altOrd++ // skip first key because it is the same as no conversion
	}

	altKey := alpha.GetRuneAt(altOrd)
	mlog.TraceT("DidimusSequencer", mlog.Rune("Prime", primeKey), mlog.Rune("Alt", altKey)) // @audit remove for release!
	return &DidimusSequencer{primeKey, altKey, 0}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

/**
 * @returns The sequencer's friendly name.
 */
func (cs *DidimusSequencer) Name() string {
	return ALG_NAME_DIDIMUS
}

/**
 * To instruct the sequencer whether it is being used for encipherment
 * or decipherment. Not relevant with Didimus variant.
 */
func (cs *DidimusSequencer) SetDecryptionMode(isDecrypting bool) {
}

/**
 * N.A.
 */
func (cs *DidimusSequencer) Feedback(rune) {
}

/**
 * Skip the current position. Must be called when character in the
 * input stream is not part of the encoding alphabet. This is important
 * to maintain integrity if non-encoding runes are removed from the
 * message at some point.
 *
 * @returns (int) number of skipped runes so far.
 */
func (cs *DidimusSequencer) Skip() int {
	cs.skipped++
	return cs.skipped
}

/**
 * Get the key to be used for encoding target rune at this position.
 * Should only be called if target is part of the encoding alphabet!
 * NOTE: Didimus variant uses a Prime Key over even character positions
 *		 and an Alternate Key for odd character positions. The odd/even
 *		 is determined based on the stream position MINUS the amount
 *		 of skipped characters.
 *
 * @param pos (int) Current character position in input stream
 * @param target (rune) ignored in this algorithm
 * @returns the basic key to use for encoding/decoding at this position.
 */
func (cs *DidimusSequencer) GetKey(pos int, target rune) rune {
	at := pos - cs.skipped // use key only on convertable positions
	if at%2 == 0 {
		return cs.prime
	}

	return cs.alt
}

func (cs *DidimusSequencer) String() string {
	return fmt.Sprintf("%s f2k(%c)/f2k+1(%c)", ALG_NAME_DIDIMUS, cs.prime, cs.alt)
}

func (cs *DidimusSequencer) GetKeyInfo() string {
	return fmt.Sprintf("%cƒ𝓍 (%s=%c,%s=%c)", UC_MATH_BOLD_D, keyEvenString, cs.prime, keyOddString, cs.alt)
}

func (cs *DidimusSequencer) Verify(callback func(rune) error) error {
	if err1 := callback(cs.prime); err1 != nil {
		return err1
	}
	if err2 := callback(cs.alt); err2 != nil {
		return err2
	}
	return nil
}

/**
 * Resets the sequencer. It should be done after every Encode or Decode
 */
func (cs *DidimusSequencer) Reset() {
	cs.skipped = 0
}
