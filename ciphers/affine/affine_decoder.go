/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   CaesarX
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Affine support added on 2025W39 out of that programmer's syndrome
 * of always wanting to add new features...
 *-----------------------------------------------------------------*/
package affine

import (
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/cmn"
	"strings"
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

type AffineDecoder struct {
	affineCryptoBase
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

// (Ctor)  Affine Encoder instance. It will return nil if the parameters
// are unsuitable.
func NewAffineDecoder(alpha *cmn.Alphabet, params *AffineParams) *AffineDecoder {
	if alpha.Size() != uint(params.N) {
		mlog.Errorf("bad params for '%s' AffineDecoder, mismatching N", alpha.Name, params.N)
		return nil
	}

	helper := NewAffineHelper()
	if err := helper.SetParameters(params.A, params.B, params.N); err != nil {
		//if err := helper.VerifyParams(params); err != nil {
		return nil
	}

	// build transliterated primary alphabet based on chosen parameters
	// this way we don't have to repeat these relative expensive calculation
	// as we decode, sort of caching.
	var decipheredAlphabet strings.Builder
	for _, charP := range alpha.Chars {
		charD, err := helper.DecodeRuneFrom(charP, alpha.Chars)
		if err != nil {
			mlog.ErrorE(err)
			return nil
		}
		decipheredAlphabet.WriteRune(charD)
	}

	// ensure there is at least a default case handler
	caser := alpha.BorrowSpecialCase()
	if caser == nil {
		caser = cmn.DefaultCaseHandler
	}

	// create the translator
	rt := cmn.NewSimpleRuneTranslator(alpha.Name, decipheredAlphabet.String(), alpha.Chars, caser)

	base := affineCryptoBase{
		langCode: alpha.LangCodeISO(),
		paramsM:  params,
		paramsS:  nil,
		master:   rt,
		slave:    nil,
	}

	return &AffineDecoder{
		affineCryptoBase: base,
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// Attach a secondary (slave) alphabet chained to the master alphabet
func (a *AffineDecoder) WithChain(alphaSlave *cmn.Alphabet) error {
	if alphaSlave == nil {
		a.slave = nil
		return nil
	}

	helper := NewAffineHelper()
	if err := helper.SetParams(a.paramsM); err != nil {
		return err
	}

	var rtSlave *cmn.RuneTranslator
	var slaveParams *AffineParams = nil
	var err error
	slaveN := int(alphaSlave.Size())
	if slaveN == a.paramsM.N {
		// alphabet lengths are equal, we can use the same parameters
		slaveParams = a.paramsM.Clone()
		mlog.Info("Affine slave A1=A2, N1=N2")
	} else { // differing alphabet lengths
		if helper.IsCommonCoprime(a.paramsM.A, a.paramsM.N, slaveN) {
			// common in both sets, no change either on A but on N
			slaveParams = a.paramsM.Clone()
			slaveParams.N = slaveN
			mlog.Info("Affine slave A1=A2, N1!=N2")
		} else {
			// based on the master parameters but restrained to the slave's
			// condition, recalculate the A coefficient to apply to the SLAVE
			// the master remains with its own A coefficient.
			A := helper.CalculateSlaveCoprime(a.paramsM.B, slaveN)
			if slaveParams, err = NewAffineParams(A, a.paramsM.B, slaveN); err != nil {
				mlog.Error("could not set Affine slave due to error", err)
				return err
			}

			mlog.InfoT("recalculated Affine A coefficient",
				mlog.Int("Master-A", a.paramsM.A),
				mlog.Int("Master-N", a.paramsM.N),
				mlog.Int("Slave-A", A),
				mlog.Int("Slave-N", slaveN))
		}
	}

	rtSlave, _ = buildTabula(alphaSlave, slaveParams, false) // don't check N
	a.paramsS = slaveParams
	a.slave = rtSlave

	return nil
}

/**
 * DESCR
 * @params a (type):
 * @returns
 */
func (a *AffineDecoder) Decode(cipher string) (string, error) {
	var plain strings.Builder
	var decR rune
	var err error = nil

	for _, charP := range cipher {
		decR = charP
		if !a.master.Exists(charP) {
			if a.slave != nil && a.slave.Exists(charP) {
				decR, err = a.slave.ReverseLookup(charP)
			}
		} else {
			decR, err = a.master.ReverseLookup(charP)
		}

		if err != nil {
			return "", err
		}
		if _, err = plain.WriteRune(decR); err != nil {
			return "", err
		}
	}

	return plain.String(), nil
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/
