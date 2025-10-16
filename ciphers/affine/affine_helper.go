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
	"lordofscripts/caesarx/internal/crypto"
	"strings"
	"unicode/utf8"
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
	helper *crypto.AffineHelper
	param  *AffineParams
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
	return &AffineHelper{crypto.NewAffineHelper(), nil}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s ( Inherited )
 *-----------------------------------------------------------------*/

// Checks if the two integers are coprime. That is, their Greatest
// Common Divisor is 1. (inmutable)
func (h *AffineHelper) AreCoprime(a, b int) bool {
	return h.helper.AreCoprime(a, b)
}

// For a given alphabet length N, get the slice of
// coprimes of N between (0..N] (inmutable)
func (h *AffineHelper) ValidCoprimesUpTo(n uint) []int {
	return h.helper.ValidCoprimesUpTo(n)
}

// Calculate the modular inverse A' of A. (inmutable)
func (h *AffineHelper) ModularInverse(a, m int) (int, error) {
	return h.helper.ModularInverse(a, m)
}

// When using a Slave Reference Alphabet with the Affine cipher,
// we cannot (and must not) use the same parameters because most
// certainly the Slave's length N isn't the same as the Master's.
// we resort to some (see documents) trickery to use the Master's
// A coefficient modulo the Slave's N to choose from the Slave's
// VALID list of coprimes.
// NOTE: From the master we use A & N, from the slave B & N and
// calculate the slave's A coefficient.
func (h *AffineHelper) CalculateSlaveCoprime(master *AffineParams, slaveB, slaveN int) int {
	return h.helper.CalculateSlaveCoprime(master.A, master.N, slaveB, slaveN)
}

// Check that A is a common coprime for alphabets of lengths N1 & N2
func (h *AffineHelper) IsCommonCoprime(a, n1, n2 int) bool {
	return h.helper.IsCommonCoprime(a, n1, n2)
}

// Get the list of all (if any) common coprimes for alphabets
// of length N1 & N2
func (h *AffineHelper) GetCommonCoprimes(n1, n2 int) (bool, []int) {
	return h.helper.GetCommonCoprimes(n1, n2)
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s ( Extra )
 *-----------------------------------------------------------------*/

// Sets the Affine parameters and calculates the A' for
// the given parameters. A' is used for decoding instead of A.
func (h *AffineHelper) SetParams(p *AffineParams) error {
	var err error = nil
	var inverse int
	if inverse, err = h.helper.VerifySettings(p.A, p.B, p.N); err == nil {
		h.param = &AffineParams{
			A:  p.A,
			B:  p.B,
			Ap: inverse,
			N:  p.N,
		}
	}

	return err
}

// GetParams retrieves the current Affine parameters of this instance
func (h *AffineHelper) GetParams() *AffineParams {
	if h.param == nil || h.param.N <= 0 {
		mlog.Error("no Affine parameters have been set!", mlog.At())
		return nil
	}

	return h.param.Clone()
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

	sb.WriteString(pre + HEAD_1[:h.param.N] + "\n")
	sb.WriteString(pre + HEAD_2[:h.param.N] + "\n")
	sb.WriteString(pre + strings.Repeat("-", h.param.N) + "\n")
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

func (h *AffineHelper) EncodeRuneFrom(r rune, alphabet string) (rune, error) {
	x := cmn.RuneIndex(alphabet, r)
	y, err := h.Encode(x)

	return cmn.RuneAt(alphabet, y), err
}

func (h *AffineHelper) Encode(x int) (int, error) {
	if h.param == nil || h.param.N <= 0 {
		return -1, fmt.Errorf("not initialiazed, call SetParameters")
	}

	y := (h.param.A*x + h.param.B) % h.param.N
	return y, nil
}

func (h *AffineHelper) DecodeRuneFrom(r rune, alphabet string) (rune, error) {
	y := cmn.RuneIndex(alphabet, r)
	x, err := h.Decode(y)

	return cmn.RuneAt(alphabet, x), err
}

func (h *AffineHelper) Decode(y int) (int, error) {
	if h.param == nil || h.param.N <= 0 {
		return -1, fmt.Errorf("not initialiazed, call SetParameters")
	}

	// (A' * ( y - B)) % N
	x := (h.param.Ap * (y - h.param.B)) % h.param.N
	return x, nil
}

func (h *AffineHelper) String() string {
	return h.param.String()
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

// VerifyParams checks the correctness of the provided Affine
// coefficients and returns an error if they are incorrect beyond repair.
//
//	If on the other hand they are almost correct but the A'
//
// (used for decoding) needs adjustment, then A' is corrected/calculated
// but no error is returned, only a log Warning entry is produced.
//
//	If the parameters are correct and the set parameter is true,
//
// the provided (and possibly corrected) parameters are stored in
// this instance (by reference, not by value!).
func (h *AffineHelper) VerifyParams(params *AffineParams, set bool) error {
	aP, err := h.helper.VerifySettings(params.A, params.B, params.N)
	if err != nil {
		return err
	}

	if aP != params.Ap {
		mlog.Warn("repaired Affine.A' after verification",
			mlog.Int("Was", params.Ap),
			mlog.Int("Becomes", aP))
		params.Ap = aP
	}

	// set the reference in this instance
	if set {
		h.param = params
	}

	return nil
}

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *					M A I N    |     D E M O
 *-----------------------------------------------------------------*/

func DemoAffine() bool {
	// these combinations MUST BE VALID
	const A int = 15
	const B int = 7
	const ALPHABET string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const SAMPLE rune = 'M'

	//toASCII := func(offset int) rune {
	//	return rune(offset + 65) // A is ASCII 65
	//}

	var aparams *AffineParams
	var err error

	// to ensure it works with multi-byte characters, we use the utf8
	// sizer instead of the standard len()
	aparams, _ = NewAffineParams(A, B, utf8.RuneCountInString(ALPHABET))
	fmt.Println("\t", aparams)

	cipher := NewAffineHelper()
	if err = cipher.SetParams(aparams); err != nil {
		fmt.Println(err)
		return false
	}

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
