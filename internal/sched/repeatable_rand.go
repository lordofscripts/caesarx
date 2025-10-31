/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A pseudo-random generator that produces a repeatable sequence for
 * any given input. It is NOT meant as a true random generator.
 *-----------------------------------------------------------------*/
package sched

import (
	"lordofscripts/caesarx/cmn"
	"math/rand"
	"sync"
	"time"
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

// Defines an interface for all Random generator providers used in
// this application so that we can interchange them.
type IRandomizer interface {
	Intn() int
	Runen(alphabet string, size int) string
}

var _ IRandomizer = (*RepeatableRand)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// A pseudo-random generator of a repeatable list of integers for
// any given input (constructor) values. It repeats values between
// the range [min, max).
type RepeatableRand struct {
	min int
	max int
	rnd *rand.Rand
	mu  *sync.Mutex
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// (ctor) New instance of a Repeatable Random integer list to generate
// in the closed range [min,max].
func NewRepeatableRand(date time.Time, userSeed int64, min, max int) *RepeatableRand {
	dateSeed := date.UTC().Unix()
	return &RepeatableRand{
		min,
		max,
		rand.New(rand.NewSource(userSeed + dateSeed)),
		new(sync.Mutex),
	}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// returns a pseudo-random integer from the [min,max] range. Depending
// on the size of the selected range, the value MAY be repeated.
func (r *RepeatableRand) Intn() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.rnd.Intn(r.max-r.min+1) + r.min
}

// returns a pseudo-random string of runes composed of the characters
// found in the alphabet.
func (r *RepeatableRand) Runen(alphabet string, size int) string {
	b := make([]rune, size)
	for range b {
		letter := cmn.RuneAt(alphabet, r.Intn())
		b = append(b, letter)
	}

	return string(b)
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// A simple hash function to convert the date string to an int
func Hash(s string) int {
	hash := 0
	for _, c := range s {
		hash = hash*31 + int(c)
	}

	return hash
}

// A simple hash function to convert the date string to an int64
func Hash64(s string) int64 {
	var hash int64 = 0
	for _, c := range s {
		hash = hash*31 + int64(c)
	}

	return hash
}
