/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package caesar

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/
const (
	CAESAR          CaesarCipherMode = iota + 1 // Plain A..Z as original cipher
	CAESAR_EXTENDED                             // Includes decimal digits and some punctuation & symbols
	CAESAR_AUGUSTUS                             // Same as Extended but double index jump back & forth: N, N+K, N, N+K...
	CAESAR_TIBERIUS                             // Same as Extended but double index two disks: N N'...
	BELLASO                                     // Giovan Battista Bellaso (1553)
	VIGNERE_AUTOKEY                             // Blas de Vignère (1850)
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type CaesarCipherMode uint

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (m CaesarCipherMode) String() string {
	return [...]string{"Caesar", "Extended Caesar", "Extended Caesar Augustus", "Extended Caesar Tiberius", "Vignere"}[m-1]
}

func (m CaesarCipherMode) EnumIndex() uint {
	return uint(m)
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/
