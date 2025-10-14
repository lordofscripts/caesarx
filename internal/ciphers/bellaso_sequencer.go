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
	"strings"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const ALG_NAME_BELLASO = "Bellaso"
const ALG_CODE_BELLASO = "BELA"

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ IKeySequencer = (*BellasoSequencer)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type BellasoSequencer struct {
	secret      []rune
	subKeyCount int
	skipped     int
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewBellasoSequencer(secret string, alpha *cmn.Alphabet) *BellasoSequencer {
	bellaso := &BellasoSequencer{nil, 0, 0}
	secret = alpha.ToUpperString(secret)
	if newSecret, err := bellaso.VerifySecret(secret, alpha); err != nil {
		mlog.Errorf("invalid secret for Bellaso sequencer: %v", "Error", err.Error())
		return nil
	} else {
		bellaso.secret = []rune(newSecret)
		bellaso.subKeyCount = len(bellaso.secret)
	}

	mlog.Debug("BellasoSequencer", mlog.Int("SecretSize", bellaso.subKeyCount)) // @audit remove for release!
	return bellaso
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

/**
 * @returns The sequencer's friendly name.
 */
func (cs *BellasoSequencer) Name() string {
	return ALG_NAME_BELLASO
}

/**
 * To instruct the sequencer whether it is being used for encipherment
 * or decipherment. Not relevant with Bellaso variant.
 */
func (cs *BellasoSequencer) SetDecryptionMode(isDecrypting bool) {
}

/**
 * N.A.
 */
func (cs *BellasoSequencer) Feedback(rune) error {
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
func (cs *BellasoSequencer) Skip() int {
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
func (cs *BellasoSequencer) GetKey(pos int, target rune) rune {
	at := pos - cs.skipped // use key only on convertable positions
	keyPos := at % cs.subKeyCount

	return cs.secret[keyPos]
}

func (cs *BellasoSequencer) String() string {
	return ALG_NAME_BELLASO
}

func (cs *BellasoSequencer) GetKeyInfo() string {
	return fmt.Sprintf("%cƒ𝓍 ('%s')", UC_MATH_BOLD_B, string(cs.secret))
}

func (cs *BellasoSequencer) Verify(callback func(rune) error) error {
	chars := []rune(cs.secret)
	for _, char := range chars {
		if err1 := callback(char); err1 != nil {
			return err1
		}
	}

	return nil
}

/**
 * All the letters in the Bellaso secret must be known in the primary/master alphabet.
 * The only exception is the SPACE/TAB characters which are eliminated from the secret.
 */
func (cs *BellasoSequencer) VerifySecret(s string, alpha *cmn.Alphabet) (string, error) {
	s = strings.Trim(s, " \t")
	for _, char := range s {
		if !alpha.Contains(char, true) {
			return "", fmt.Errorf("subkey %c not present in %s alphabet for Bellaso", char, alpha.Name)
		}
	}

	return s, nil
}

/**
 * Resets the sequencer. It should be done after every Encode or Decode
 */
func (cs *BellasoSequencer) Reset() {
	cs.skipped = 0
}
