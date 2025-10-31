/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A pseudo-random generator that produces a repeatable sequence of
 * unique integers for any given input. It is NOT meant as a true
 * random generator.
 *-----------------------------------------------------------------*/
package sched

import (
	"fmt"
	"time"
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

// Defines an interface for all Unique Random generator providers used in
// this application so that we can interchange them.
type IUniqueRandomizer interface {
	Intn() (int, error)
	Runen(alphabet string, size int) string
}

var _ IUniqueRandomizer = (*RepeatableUniqueRand)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// A pseudo-random generator of a repeatable list of UNIQUE integers for
// any given input (constructor) values. It repeats values between
// the range [min, max).
type RepeatableUniqueRand struct {
	RepeatableRand
	uniq map[int]struct{}
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// (ctor) New instance of a Repeatable Random integer list
func NewRepeatableUniqueRand(date time.Time, userSeed int64, min, max int) *RepeatableUniqueRand {
	return &RepeatableUniqueRand{
		*NewRepeatableRand(date, userSeed, min, max),
		make(map[int]struct{}),
	}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// returns a pseudo-random integer from the [(]min,max) range. Depending
// on the size of the selected range, the value MAY be repeated.
func (r *RepeatableUniqueRand) Intn() (int, error) {
	const MAX_TRY int = 100
	tried := 0
	var result int = -1
	var err error = nil

	for {
		nr := (&r.RepeatableRand).Intn()
		if _, ok := r.uniq[nr]; !ok {
			r.uniq[nr] = struct{}{}
			result = nr
			break
		} else {
			tried += 1
			if tried >= MAX_TRY {
				err = fmt.Errorf("too many attempts (%d) to get unique number [%d,%d)]", tried, r.min, r.max)
			}
		}
	}

	return result, err
}

func (r *RepeatableUniqueRand) Capacity() int {
	return r.max - r.min - len(r.uniq)
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/
