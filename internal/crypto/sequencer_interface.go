/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * The Key Sequencer is a flexible generator of the current key
 * during Encoding/Decoding.
 *-----------------------------------------------------------------*/
package ciphers

const (
	UC_MATH_FUNC   rune = rune(0x1d453)
	UC_MATH_K      rune = rune(0x1d458)
	UC_MATH_X      rune = rune(0x1d465)
	UC_MATH_BOLD_A rune = rune(0x1d468) // for Affine
	UC_MATH_BOLD_B rune = rune(0x1d469) // for Bellaso
	UC_MATH_BOLD_C rune = rune(0x1d46a) // for Caesar
	UC_MATH_BOLD_D rune = rune(0x1d46b) // for Didimus
	UC_MATH_BOLD_F rune = rune(0x1d46d) // for Fibonacius
	UC_MATH_BOLD_V rune = rune(0x1d47d) // for Vigenère

	UC_SUBSCRIPT_0 = rune(0x2080)
	UC_SUBSCRIPT_1 = rune(0x2081)
)

var (
	keyEvenString string = string(UC_MATH_K) + string(UC_SUBSCRIPT_0)
	keyOddString  string = string(UC_MATH_K) + string(UC_SUBSCRIPT_1)
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/**
 * Each algorithm in the Caesar family of variants implements its
 * own sequencer. i.e. CaesarSequencer, DidimusSequencer, BellasoSequencer,
 * and VigenereSequencer.
 */
type IKeySequencer interface {
	/**
	 * The sequencer's friendly name.
	 */
	Name() string
	/**
	 * Skip the current position. Must be called when character in the
	 * input stream is not part of the encoding alphabet. This is important
	 * to maintain integrity if non-encoding runes are removed from the
	 * message at some point.
	 * @returns (int) number of skipped runes so far.
	 */
	Skip() int
	/**
	 * Get the key to be used for encoding target rune at this position.
	 * Should only be called if target is part of the encoding alphabet!
	 * @returns the basic key to use for encoding/decoding at this position.
	 */
	GetKey(pos int, target rune) rune

	GetKeyInfo() string

	/**
	 * The callback is used to verify that the rune can be used
	 * as a key.
	 */
	Verify(func(rune) error) error

	/**
	 * Resets the sequencer. It should be done after every Encode or Decode
	 */
	Reset()

	/**
	 * To instruct the sequencer whether it is being used for encipherment
	 * or decipherment. The Vigenere variant's sequencer works differently
	 * depending on whether it does encryption/decryption. For the other
	 * sequencers, it makes no difference.
	 */
	SetDecryptionMode(isDecrypting bool)

	/**
	 * Only for algorithms that need to feed the last decrypted rune
	 * back into the sequencer (i.e. Vigenère), else does nothing.
	 */
	Feedback(rune) error
}
