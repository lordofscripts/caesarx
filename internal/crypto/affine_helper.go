/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   CaesarX
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Affine is a cipher similar to Caesar, in fact Caesar is a variant
 * of Affine. For this cipher we have two coefficients A & B plus
 * the size of the alphabet N. For encryption the result 'y' is given:
 *			y = (A * x + B) % N
 * where 'x' is the (zero-based) "shift" of the letter in the alphabet.
 * The restriction is that A and N must be coprime (their GCD = 1).
 * For an alphabet with 26 letters (N=26) the values of A are
 * constrained to: 1, 3, 5, 7, 9, 11, 15, 17, 19, 21, 23 & 25.
 * For decryption the reciprocal formula becomes:
 *			y' = A' * x + B
 * where A' is the Modular Inverse of A modulo N.
 *-----------------------------------------------------------------*/
package crypto

import (
	"fmt"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/cmn"
	"math"
	"math/big"
	"slices"
	"unsafe"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

// AffineHelper (internal version) contains no state information.
// It simply provides validation and calculation methods needed for
// propper support of Affine encoding/decoding operations.
type AffineHelper struct {
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// (Ctor) internal AffineHelper which only provides calculation and
// validation methods. It holds no state information.
func NewAffineHelper() *AffineHelper {
	return &AffineHelper{}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

/**
 * Checks if the two integers are coprime. That is, their Greatest
 * Common Divisor is 1.
 */
func (h *AffineHelper) AreCoprime(a, b int) bool {
	A := big.NewInt(int64(a))
	B := big.NewInt(int64(b))
	gcd := new(big.Int).GCD(nil, nil, A, B)
	// Check that their GCD is 1
	return gcd.Cmp(big.NewInt(1)) == 0
}

/**
 * For a given alphabet length N, get the slice of
 * coprimes of N between (0..N]
 */
func (h *AffineHelper) ValidCoprimesUpTo(n uint) []int {
	coprimes := make([]int, 0)
	for v := range n {
		if v == 0 {
			continue
		}

		if h.AreCoprime(int(v), int(n)) {
			coprimes = append(coprimes, int(v))
		}
	}

	return coprimes
}

/**
 * Calculate the modular inverse A' of A
 */
func (h *AffineHelper) ModularInverse(a, m int) (int, error) {
	A := big.NewInt(int64(a))
	M := big.NewInt(int64(m))

	if !h.AreCoprime(a, m) {
		return -1, fmt.Errorf("values a=%d and n=%d are not coprime", a, m)
	}

	inverse := new(big.Int).ModInverse(A, M)
	inverse64 := inverse.Int64()
	if inverse64 > int64(maxIntValue()) {
		return -1, fmt.Errorf("this overflow shouldn't occur in this application")
	}

	return int(inverse64), nil
}

// When using a Slave Reference Alphabet with the Affine cipher,
// we cannot (and must not) use the same parameters because most
// certainly the Slave's length N isn't the same as the Master's.
// we resort to some (see documents) trickery to use the Master's
// A coefficient modulo the Slave's N to choose from the Slave's
// VALID list of coprimes.
func (h *AffineHelper) CalculateSlaveCoprime(masterA, masterN int, slaveB, slaveN int) int {
	// get the list of allowable A coefficients of the slave
	slaveCoprimes := h.ValidCoprimesUpTo(uint(slaveN))
	// calculate a modulo index for that using the current (master's)
	// A coefficient
	selected := masterA % len(slaveCoprimes)
	M := slaveCoprimes[selected]
	// pick an adjusted, and thus valid, A coefficient for the slave
	mlog.DebugT("adjusted A coefficient",
		mlog.Int("Master-A", masterA),
		mlog.Int("Master-N", masterN),
		mlog.Int("Slave-A (adj)", M),
		mlog.Int("Slave-N", slaveN))
	return M
}

// Check that A is a common coprime for alphabets of lengths N1 & N2
func (h *AffineHelper) IsCommonCoprime(a, n1, n2 int) bool {
	set1 := h.ValidCoprimesUpTo(uint(n1))
	set2 := h.ValidCoprimesUpTo(uint(n2))

	return slices.Contains(set1, a) && slices.Contains(set2, a)
}

// Get the list of all (if any) common coprimes for alphabets
// of length N1 & N2
func (h *AffineHelper) GetCommonCoprimes(n1, n2 int) (bool, []int) {
	set1 := h.ValidCoprimesUpTo(uint(n1))
	set2 := h.ValidCoprimesUpTo(uint(n2))

	common := cmn.IntersectInt(set1, set2)

	return len(common) > 0, common
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

// VerifySettings validates the given Affine parameters and returns
// nil error and the A' (needed for decoding) if all is good.
func (h *AffineHelper) VerifySettings(a, b, n int) (int, error) {
	const INVALID int = -1
	var inverse int
	var err error = nil

	if a < 0 || b < 0 || n <= 0 {
		return INVALID, fmt.Errorf("a, b, n should be positive")
	}

	if inverse, err = h.ModularInverse(a, n); err != nil {
		return INVALID, err
	}

	contains := func(slice []int, value int) bool {
		for _, v := range slice {
			if v == value {
				return true
			}
		}
		return false
	}

	validCoprimes := h.ValidCoprimesUpTo(uint(n))
	if !contains(validCoprimes, a) {
		return INVALID, fmt.Errorf("coeficcient A=%d is not a valid coprime of N=%d", a, n)
	}

	return inverse, nil
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

/**
 * The maximum value of an 'int' type depends on the architecture.
 * I.e. whether we are executing on a 32 or 64-bit CPU.
 */
func maxIntValue() int {
	if unsafe.Sizeof(int(0)) == 8 {
		return math.MaxInt64
	}

	return math.MaxInt32
}
