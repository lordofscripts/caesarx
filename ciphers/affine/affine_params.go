/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   CaesarX
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Creates an instance of validated Affine parameters.
 *-----------------------------------------------------------------*/
package affine

import (
	"fmt"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/ciphers"
	iciphers "lordofscripts/caesarx/internal/ciphers"
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/

func init() {
	ciphers.RegisterCipher(Info)
}

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/
var (
	Info = ciphers.NewCipherInfo(iciphers.ALG_CODE_AFFINE, "1.0",
		"Unknown",
		iciphers.ALG_NAME_AFFINE,
		"Affine linear")
)

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

// The parameters that define an Affine encoding/decoding transform.
// where the encoding affine is f(x) = (A * x + B) % N
// the program should NEVER manipulate these parameters directly.
type AffineParams struct {
	A  int // coefficient A
	B  int // coefficient B
	Ap int // coefficient A' (for decoding only)
	N  int // module (alphabet length)
}

/* ----------------------------------------------------------------
 *				C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// Sets the Affine parameters, calculates the A' coefficient, checks
// that A and N are coprimes and only then returns a non-nil instance
// of verified parameters suitable for both encription and decription.
func NewAffineParams(a, b, n int) (*AffineParams, error) {
	helper := NewAffineHelper()
	if err := helper.SetParameters(a, b, n); err != nil {
		mlog.ErrorE(err)
		return nil, err
	}

	return helper.GetParams(), nil
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

func (a *AffineParams) Clone() *AffineParams {
	clone, _ := NewAffineParams(a.A, a.B, a.N)
	clone.Ap = a.Ap
	return clone
}

/**
 * DESCR
 * @params a (type):
 * @returns
 */
func (a *AffineParams) String() string {
	return fmt.Sprintf("Affine ::= A=%d B=%d N=%d A'=%d", a.A, a.B, a.N, a.Ap)
}
