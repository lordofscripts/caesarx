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
package affine

import (
	"fmt"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/cmn"
	"math"
	"math/big"
	"slices"
	"strings"
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

type AffineHelper struct {
	cA   int
	cB   int
	cAp  int
	modN int
}

/* ----------------------------------------------------------------
 *				P r i v a t e	T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				I n i t i a l i z e r
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

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
		return -1, fmt.Errorf("this overflow shouldn't ocurr in this application")
	}

	return int(inverse64), nil
}

func (h *AffineHelper) SetParams(p *AffineParams) error {
	var err error
	var inverse int
	if inverse, err = h.VerifySettings(p.A, p.B, p.N); err == nil {
		h.cA = p.A
		h.cB = p.B
		h.cAp = inverse
		h.modN = p.N
	}

	return err
}

func (h *AffineHelper) SetParameters(a, b, n int) error {
	if inverse, err := h.VerifySettings(a, b, n); err != nil {
		return err
	} else {
		h.cAp = inverse
		h.cA = a
		h.cB = b
		h.modN = n
	}

	return nil
}

func (h *AffineHelper) GetParams() *AffineParams {
	if h.modN == 0 {
		mlog.Error("AffineHelper.GetParams but there are none")
		return nil
	}

	return &AffineParams{
		A:  h.cA,
		B:  h.cB,
		Ap: h.cAp,
		N:  h.modN,
	}
}

func (h *AffineHelper) GetParameters() (a, ap, b, n int) {
	return h.cA, h.cAp, h.cB, h.modN
}

// A multi-line string block with the Tabula Recta for the current
// Affine parameters.
// @param alphabet (string) plain alphabet string by shift order
// @returns (string) the corresponding Tabula Recta
// @returns (error) nil on success, else encountered error
func (h *AffineHelper) GetTabulaString(alphabet string, leader ...string) (string, error) {
	const HEAD_1 string = "0.........1.........2.........3.........4"
	const HEAD_2 string = "01234567890123456789012345678901234567890"
	var sb strings.Builder

	var pre string = ""
	if leader != nil {
		pre = leader[0]
	}

	sb.WriteString(pre + HEAD_1[:h.modN] + "\n")
	sb.WriteString(pre + HEAD_2[:h.modN] + "\n")
	sb.WriteString(pre + strings.Repeat("-", h.modN) + "\n")
	sb.WriteString(pre + alphabet + "\n") // plain alphabet

	// now the conversion duplet
	sb.WriteString(pre)
	for _, char := range alphabet {
		if r, err := h.EncodeRuneFrom(char, alphabet); err != nil {
			return "", err
		} else {
			sb.WriteRune(r)
		}
	}

	return sb.String(), nil
}

// When using a Slave Reference Alphabet with the Affine cipher,
// we cannot (and must not) use the same parameters because most
// certainly the Slave's length N isn't the same as the Master's.
// we resort to some (see documents) trickery to use the Master's
// A coefficient modulo the Slave's N to choose from the Slave's
// VALID list of coprimes.
func (h *AffineHelper) CalculateSlaveCoprime(slaveB, slaveN int) int {
	// get the list of allowable A coefficients of the slave
	slaveCoprimes := h.ValidCoprimesUpTo(uint(slaveN))
	// calculate a modulo index for that using the current (master's)
	// A coefficient
	selected := h.cA % len(slaveCoprimes)
	M := slaveCoprimes[selected]
	// pick an adjusted, and thus valid, A coefficient for the slave
	mlog.DebugT("adjusted A coefficient", mlog.Int("Master-A", h.cA),
		mlog.Int("Master-N", h.modN),
		mlog.Int("Slave-A", M),
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

func (h *AffineHelper) EncodeRuneFrom(r rune, alphabet string) (rune, error) {
	x := cmn.RuneIndex(alphabet, r)
	y, err := h.Encode(x)

	return cmn.RuneAt(alphabet, y), err
}

func (h *AffineHelper) Encode(x int) (int, error) {
	if h.modN <= 0 {
		return -1, fmt.Errorf("not initialiazed, call SetParameters")
	}

	y := (h.cA*x + h.cB) % h.modN
	return y, nil
}

func (h *AffineHelper) DecodeRuneFrom(r rune, alphabet string) (rune, error) {
	y := cmn.RuneIndex(alphabet, r)
	x, err := h.Decode(y)

	return cmn.RuneAt(alphabet, x), err
}

func (h *AffineHelper) Decode(y int) (int, error) {
	if h.modN <= 0 {
		return -1, fmt.Errorf("not initialiazed, call SetParameters")
	}

	// (A' * ( y - B)) % N
	x := (h.cAp * (y - h.cB)) % h.modN
	return x, nil
}

func (h *AffineHelper) String() string {
	return fmt.Sprintf("Affine{A:%d,B:%d,A':%d,N:%d}", h.cA, h.cB, h.cAp, h.modN)
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

func (h *AffineHelper) VerifyParams(params *AffineParams) error {
	aP, err := h.VerifySettings(params.A, params.B, params.N)
	if err != nil {
		return err
	}

	if aP != params.Ap {
		mlog.Warn("repaired Affine.A' after verification",
			mlog.Int("Was", params.Ap),
			mlog.Int("Becomes", aP))
		params.Ap = aP // @audit WTF this is not reflected upon return
		panic("unexplained... not reflected on caller despite pointers")
	}

	return nil
}

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

/* ----------------------------------------------------------------
 *					M A I N    |     D E M O
 *-----------------------------------------------------------------*/

func DemoAffine() bool {
	const A int = 15
	const B int = 7
	const ALPHABET string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const SAMPLE rune = 'M'

	//toASCII := func(offset int) rune {
	//	return rune(offset + 65) // A is ASCII 65
	//}

	cipher := NewAffineHelper()
	// if Alphabet contains multi-byte chars use utf8.RuneCountInString()
	cipher.SetParameters(A, B, len(ALPHABET))
	fmt.Println("\t", cipher)

	y, err := cipher.EncodeRuneFrom(SAMPLE, ALPHABET)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\t·Encode('%c') = '%c'\n", SAMPLE, y)

	x, err := cipher.DecodeRuneFrom(y, ALPHABET)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\t·Decode('%c') = '%c'\n", y, x)

	return x == SAMPLE
}
