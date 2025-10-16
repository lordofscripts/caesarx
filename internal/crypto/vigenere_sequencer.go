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
	"lordofscripts/caesarx"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/cmn"
	"strings"

	roundrobin "github.com/lordofscripts/go-roundrobin"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const ALG_NAME_VIGENERE = "Vigenère"
const ALG_CODE_VIGENERE = "VIGN"

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ IKeySequencer = (*VigenereSequencer)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type VigenereSequencer struct {
	secret      []rune
	subKeyCount int
	skipped     int
	isDecoding  bool
	//buffer roundrobin.IRingQueue[rune] // only for Decryption to progressively build auto-key
	buffer *roundrobin.RuneRingQueue
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewVigenereSequencer(secret string, alpha *cmn.Alphabet) *VigenereSequencer {
	var vg *VigenereSequencer

	secret = alpha.ToUpperString(secret)
	if newSecret, err := vg.VerifySecret(secret, alpha); err != nil {
		mlog.Errorf("invalid secret for Vigenère sequencer: %v", "Error", err.Error())
		vg = nil
	} else {
		baseKeyRunes := []rune(newSecret)
		baseKeyLen := len(baseKeyRunes)

		vg = &VigenereSequencer{
			secret:      baseKeyRunes,
			subKeyCount: baseKeyLen,
			skipped:     0,
			isDecoding:  false,
			buffer:      roundrobin.NewRuneRingQueue(baseKeyLen),
		}
	}

	mlog.TraceT("VigenèreSequencer", mlog.Int("SecretSize", vg.subKeyCount)) // @audit remove for release!
	return vg
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

/**
 * @returns The sequencer's friendly name.
 */
func (cs *VigenereSequencer) Name() string {
	return ALG_NAME_VIGENERE
}

/**
 * To instruct the sequencer whether it is being used for encipherment
 * or decipherment.
 */
func (cs *VigenereSequencer) SetDecryptionMode(isDecrypting bool) {
	// the sequencer works differently during Decode
	cs.isDecoding = isDecrypting // @note Deprecate
}

/**
 * Vigenère en/decryption feeds back the last en/decoded rune into the
 * sequencer to progressively rebuild the auto-key .
 */
func (cs *VigenereSequencer) Feedback(r rune) error {
	if size, err := cs.buffer.Push(r); err != nil {
		mlog.ErrorT(
			"couldn't push decoded rune.",
			mlog.At(),
			mlog.Int("Size", size),
			mlog.Int("Skipped", cs.skipped),
			mlog.Rune("Rune", r),
			mlog.Err(err))
		return err
	}

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
func (cs *VigenereSequencer) Skip() int {
	cs.skipped++
	return cs.skipped
}

/**
 * Get the key to be used for encoding target rune at this position.
 * Should only be called if target is part of the encoding alphabet!
 * NOTE: Vigenère autokey uses the secret for the first part of the message.
 * 		 If the message is longer than the secret, the secret is NOT
 *		 repeated (as in Bellaso), instead the current message's
 *		 character is used as the basic Caesar key.
 * Pre-Condition:
 *	- This method is ONLY called IF the target rune is present in the
 *	  primary/master alphabet.
 *	- The Vigenère 'secret' is composed only of characters in Master
 *	  alphabet (plus Spaces or Tabs which are removed).
 * @param pos (int) ignored in this algorithm
 * @param target (rune) ignored in this algorithm
 * @returns the basic key to use for encoding/decoding at this position.
 */
func (cs *VigenereSequencer) GetKey(pos int, target rune) rune {
	var currentKey rune
	at := pos - cs.skipped // use key only on convertable positions
	if at >= cs.subKeyCount {
		// auto-key
		// we pop the last decoded rune as the progressive auto-key
		// to decode the target rune.
		var size int
		var err error
		currentKey, size, err = cs.buffer.Pop()
		if err != nil {
			mlog.Fatal(caesarx.ERR_SEQUENCER,
				"couldn't pop key rune.",
				//mlog.String("At", "VigenereSequencer.GetKey"),
				mlog.At(),
				mlog.Int("Size", size),
				mlog.Rune("Target", target),
				mlog.Int("Pos", pos))
		}
	} else {
		// secret ()
		keyPos := at % cs.subKeyCount
		currentKey = cs.secret[keyPos]
	}

	return currentKey
}

func (cs *VigenereSequencer) String() string {
	return ALG_NAME_VIGENERE
}

func (cs *VigenereSequencer) GetKeyInfo() string {
	mode := "encode"
	if cs.isDecoding {
		mode = "decode"
	}
	return fmt.Sprintf("%cƒ𝓍 ('%s',autokey,%s)", UC_MATH_BOLD_V, string(cs.secret), mode)
}

func (cs *VigenereSequencer) Verify(callback func(rune) error) error {
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
func (cs *VigenereSequencer) VerifySecret(s string, alpha *cmn.Alphabet) (string, error) {
	s = strings.Trim(s, " \t")
	for _, char := range s {
		if !alpha.Contains(char, true) {
			return "", fmt.Errorf("subkey %c not present in %s alphabet for Vigenère", char, alpha.Name)
		}
	}

	return s, nil
}

/**
 * Resets the sequencer. It should be done after every Encode or Decode
 */
func (cs *VigenereSequencer) Reset() {
	cs.skipped = 0
	if cs.buffer != nil {
		cs.buffer.Reset()
	}
}

/*
Key    : KEY
Plain  : MESSAGE
AutoKey: KEYMESSage

ENCRYPT:
	M E S S A G E
	k e y m e s s
	W I Q E E Y W

DECRYPT: (with M as key and cipher E we get M, with E as key and cipher E we get E and so on...)
	W I Q E E Y W
	k e y M E S
	M E S s a g
*/
