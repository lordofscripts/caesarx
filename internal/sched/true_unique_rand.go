/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A true random generator a sequence of unique integers for any given
 * input. It is based on crypto/rand operations.
 *-----------------------------------------------------------------*/
package sched

import (
	"fmt"
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ IUniqueRandomizer = (*TrueUniqueRand[int])(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// A pseudo-random generator of a repeatable list of UNIQUE integers for
// any given input (constructor) values. It repeats values between
// the range [min, max).
type TrueUniqueRand[T int | int64] struct {
	TrueRand[T]
	uniq map[T]struct{}
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// (ctor) New instance of a Repeatable Random integer list
func NewTrueUniqueRand[T int | int64](min, max T) *TrueUniqueRand[T] {
	return &TrueUniqueRand[T]{
		*NewTrueRand(min, max, false, false),
		make(map[T]struct{}),
	}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// returns a pseudo-random integer from the [(]min,max) range. Depending
// on the size of the selected range, the value MAY be repeated.
func (r *TrueUniqueRand[T]) Intn() (T, error) {
	const MAX_TRY int = 100
	tried := 0
	var result T = -1
	var err error = nil

	for {
		nr := (&r.TrueRand).Intn()
		if _, ok := r.uniq[nr]; !ok {
			r.uniq[nr] = struct{}{}
			result = nr
			break
		} else {
			tried += 1
			if tried >= MAX_TRY {
				err = fmt.Errorf("too many attempts (%d) to get unique number [%d,%d] have %d", tried, r.min, r.max, len(r.uniq))
				break
			}
		}
	}

	return result, err
}

func (r *TrueUniqueRand[T]) Capacity() T {
	return r.max - r.min - T(len(r.uniq))
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/
