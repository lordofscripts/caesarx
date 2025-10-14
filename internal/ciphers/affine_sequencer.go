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
	"lordofscripts/caesarx/cmn"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const ALG_NAME_AFFINE = "Affine"
const ALG_CODE_AFFINE = "AFIN"

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ IKeySequencer = (*AffineSequencer)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type AffineSequencer struct {
	coefA   int
	coefAp  int
	coefB   int
	modulo  int
	skipped int
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/**
 * Instantiates an Affine Sequencer. Caller must ensure that the
 * parameters are valid using affine.AffineHelper()
 * @audit circular dependency if we include AffineHelper here.
 */
func NewAffineSequencer(a, ap, b int, alpha *cmn.Alphabet) *AffineSequencer {
	return &AffineSequencer{
		coefA:   a,
		coefB:   b,
		coefAp:  ap,
		modulo:  int(alpha.Size()),
		skipped: 0,
	}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// The sequencer's friendly name.
func (cs *AffineSequencer) Name() string {
	return ALG_NAME_AFFINE
}

/**
 * N.A.
 */
func (cs *AffineSequencer) SetDecryptionMode(isDecrypting bool) {
}

// (N.A.) Affine does not need stream feedback.
func (cs *AffineSequencer) Feedback(r rune) error {
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
func (cs *AffineSequencer) Skip() int {
	cs.skipped++
	return cs.skipped
}

// (N.A.) Affine does not use a "key" but coefficients to a formula.
func (cs *AffineSequencer) GetKey(pos int, target rune) rune {
	return target
}

func (cs *AffineSequencer) String() string {
	return ALG_NAME_AFFINE
}

func (cs *AffineSequencer) GetKeyInfo() string {
	return fmt.Sprintf("%cƒ𝓍 (A:%d,B:%d,N:%d)", UC_MATH_BOLD_A, cs.coefA, cs.coefB, cs.modulo)
}

// N.A.
func (cs *AffineSequencer) Verify(callback func(rune) error) error {
	return nil
}

// (N.A.) Affine does not use a "secret" but coefficients to a formula.
func (cs *AffineSequencer) VerifySecret(s string, alpha *cmn.Alphabet) (string, error) {
	return s, nil
}

// Resets the sequencer. It should be done after every Encode or Decode
func (cs *AffineSequencer) Reset() {
	cs.skipped = 0
}
