/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A pseudo-random generator that produces a repeatable sequence for
 * any given input. It is NOT meant as a true random generator.
 *-----------------------------------------------------------------*/
package sched

import (
	"crypto/rand"
	"encoding/base64"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/cmn"
	"math/big"
	"regexp"
	"sync"
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ IRandomizer = (*TrueRand[int])(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// A pseudo-random generator of a repeatable list of integers for
// any given input (constructor) values. It repeats values between
// the range [min, max).
type TrueRand[T int | int64] struct {
	min          T
	max          T
	useBase32    bool
	removeDigits bool
	mu           *sync.Mutex
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// (ctor) New instance of a Repeatable Random integer list to generate
// in the closed range [min,max].
func NewTrueRand[T int | int64](min, max T, preferBase32, noDigits bool) *TrueRand[T] {
	return &TrueRand[T]{
		min,
		max,
		preferBase32,
		noDigits,
		new(sync.Mutex),
	}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// returns a pseudo-random integer from the [min,max] range. Depending
// on the size of the selected range, the value MAY be repeated.
func (r *TrueRand[T]) Intn() T {
	r.mu.Lock()
	defer r.mu.Unlock()

	// a true random number in the closed range [min,max]
	n, err := rand.Int(rand.Reader, big.NewInt(int64(r.max-r.min+1)))
	if err != nil {
		mlog.FatalT(120, "fatal randomness", mlog.Err(err), mlog.At())
	}

	nCasted := n.Int64()
	return T(nCasted) + r.min
}

// returns a pseudo-random string of runes composed of the characters
// found in the alphabet.
func (r *TrueRand[T]) Runen(alphabet string, size int) string {
	var result string
	if r.useBase32 {
		result = GenerateToken(r.removeDigits)
	} else {
		b := make([]rune, size)
		for range b {
			letter := cmn.RuneAt(alphabet, int(r.Intn()))
			b = append(b, letter)
		}
		result = string(b)
	}

	return result
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
// For secure tokens ensure N is at least 24 (192-bits) but preferably
// 32 (256-bits).
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// Alternatively, use a more secure crypto/rand.Text()
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// Generates a random secure token using cryptographically-secure
// methods and returns a string of 26 characters in Base32. If
// the removeDigits parameter is true, all digits are removed
// and the result will be a shorter string, how much shorter it
// is not predictable, and therefore being shorter sacrifices some
// of the gains of using rand.Text(). However, in the CaesarX context
// the built-in primary alphabets contain only letters, unless the
// user makes a custom alphabet as primary.
func GenerateToken(removeDigits bool) string {
	token := rand.Text() // @note Requires GO v1.24+
	if removeDigits {
		re := regexp.MustCompile("[0-9]+")
		return re.ReplaceAllString(token, "")
	}

	return token
}
