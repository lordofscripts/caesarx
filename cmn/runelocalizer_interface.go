/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package cmn

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

type IRuneLocalizer interface {
	/**
	 * Find a rune in the object's alphabet catalog.
	 * Rune not found: error set, other return values nil/empty or -1.
	 * Rune found: error nil, pointer to alphabet and position within.
	 */
	FindRune(rune) (string, int, error)
}
