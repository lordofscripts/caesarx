/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Caesarium is a code book for CaesarX. It may generate truly random
 * code books, but on request it can also generate recoverable codebooks.
 * A recoverable codebook is one generated based on a user "secret",
 * and it would use a pseudo-random generator so that every time that
 * same "secret" is used, it would generate the same original codebook.
 *-----------------------------------------------------------------*/
package sched

import (
	"fmt"
	"lordofscripts/caesarx"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/cmn"
	"lordofscripts/caesarx/internal/bip39"
	"lordofscripts/caesarx/internal/crypto"
	"strings"
	"time"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	// The default Codebook secret word length for Bellaso & Vigenère
	DEFAULT_SECRET_LENGTH int = 26

	extraOffsetSeed int64 = 98254762
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type Caesarium struct {
	title      string
	alphabet   *cmn.Alphabet
	date       time.Time
	userSeed   int64
	alphaLen   int
	repeatable bool
	yearBook   []caesarx.CipherVariant
}

// A bi-parametric is used for Didimus & Fibonacci and correspond
// to the shift of the (main) key and the secondary offset for
// generating the secondary key in order to produce feed a
// bi-alphabetic cipher.
type BiParametric struct {
	A int
	B int
}

type TriParametric struct {
	A int
	B int
	C int
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// instantiate a new Caesarium.
func NewCaesarium(title string, alpha *cmn.Alphabet, date time.Time, userSeed int64) *Caesarium {
	// we need to strip out the time variant
	useDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	return &Caesarium{
		title:      title,
		alphabet:   alpha,
		date:       useDate,
		userSeed:   userSeed,
		alphaLen:   int(alpha.Size()),
		repeatable: false,
		yearBook:   make([]caesarx.CipherVariant, 0),
	}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (b BiParametric) String() string {
	return fmt.Sprintf("%d, %d", b.A, b.B)
}

// implements fmt.Stringer
func (c *Caesarium) String() string {
	return fmt.Sprintf("'%s' %s", c.title, c.yearBook)
}

// By default the Caesarium generates truly random codebooks using
// cryptographically secure libraries. But, if there is a need to
// have a code book that can be recovered via a secret (yes, there is
// a valid use case!), then use this method so that given the same
// recovery secret, the Caesarium will generate the same codebook
// as the original.
func (c *Caesarium) MakeRecoverable(recovery, passphrase string) *Caesarium {
	mnemonics := strings.Fields(recovery)

	return c.MakeRecoverableFromList(mnemonics, passphrase)
}

func (c *Caesarium) MakeRecoverableFromList(recovery []string, passphrase string) *Caesarium {
	if modeBIP, err := bip39.Bip39Words12.Convert(len(recovery)); err == nil {
		bip := bip39.NewBip39(modeBIP, ' ')
		_, reducedSeed := bip.ToSeedAlt(recovery, passphrase)
		c.userSeed = int64(reducedSeed)
		c.repeatable = true
	} else {
		mlog.Error("cannot make recoverable, not a proper BIP39 recovery length", mlog.At())
	}

	return c
}

// Use the Caesarium year book to get the cipher to use for the given
// month. The month number must be 1..12.
func (c *Caesarium) GetCipherForMonth(monthNr time.Month) caesarx.CipherVariant {
	var result caesarx.CipherVariant = caesarx.NoCipher

	if len(c.yearBook) == 0 {
		mlog.Fatal(125, "Caesarium year book has not yet been generated!")
	}

	if monthNr > 0 && monthNr <= 12 {
		result = c.yearBook[monthNr-1]
	}

	return result
}

// Get the year book that states which cipher should be used
// every month of the year for the selected year.
func (c *Caesarium) CompileYearBook() []caesarx.CipherVariant {
	// omit 0 (NoCipher) up to AffineCipher.
	var rnd IRandomizer
	if c.repeatable {
		// The selected date & user seed will always generate the same YearBook
		rnd = NewRepeatableRand(c.date, c.userSeed, 1, caesarx.MaxCipher())
	} else {
		// Every time a truly random YearBook
		rnd = NewTrueRand(1, caesarx.MaxCipher(), false, false)
	}

	// generate the year book
	c.yearBook = make([]caesarx.CipherVariant, 12)

	for monthNr := range time.December {
		cipherId := rnd.Intn() //@audit despite min=1 sometimes I get NoCipher on Convert()
		selectCipher, _, err := caesarx.NoCipher.Convert(cipherId)
		mlog.TraceT("compile yearbook", mlog.Int("Id", cipherId), mlog.String("Cipher", selectCipher.String()))
		if err != nil {
			mlog.FatalT(120, "error parsing cipher variant", mlog.Err(err), mlog.At())
		}
		c.yearBook[monthNr] = selectCipher
	}

	return c.yearBook
}

// generate a Caesar code booklet for the specific
// month of that year indicating the key (as a shift number) taking
// into consideration N which is the amount of characters in the
// primary alphabet as given in the constructor.
//
// Note: A wrapper function could render the shift number as a rune.
func (c *Caesarium) CompileCaesarBook() []int {
	// we omit the value 0 because that results in no encryption
	var rnd IRandomizer
	if c.repeatable {
		rnd = NewRepeatableRand(c.date, c.userSeed, 1, c.alphaLen-1)
	} else {
		rnd = NewTrueRand(1, c.alphaLen-1, false, false)
	}

	totalDays := DaysInMonth(c.date)
	paramBooklet := make([]int, totalDays)
	for day := range totalDays {
		paramBooklet[day] = rnd.Intn()
	}

	return paramBooklet
}

// generate a  Didimus & Fibonacci cipher code booklet for the specific
// month of that year indicating the key (as a shift number) taking
// into consideration N which is the amount of characters in the
// primary alphabet as given in the constructor.
//
// Note: A wrapper function could render the shift number as a rune.
func (c *Caesarium) CompileBiAlphabeticBook() []BiParametric {
	// for the main key, it is bound by N
	var rndM IRandomizer
	if c.repeatable {
		rndM = NewRepeatableRand(c.date, c.userSeed, 1, c.alphaLen-1)
	} else {
		rndM = NewTrueRand(1, c.alphaLen-1, false, false)
	}

	// for the offset we want a different list but predictable
	offsetSeed := c.userSeed + extraOffsetSeed
	// for the secondary key, expanded but internally the cipher applies modulo N
	// expand the range to avoid depletion of pool when N < DaysInAmonth
	var rndO IUniqueRandomizer
	if c.repeatable {
		rndO = NewRepeatableUniqueRand(c.date, offsetSeed, 1, 3*c.alphaLen)
	} else {
		rndO = NewTrueUniqueRand(1, 3*c.alphaLen)
	}

	// the last day of this month
	totalDays := DaysInMonth(c.date)
	paramBooklet := make([]BiParametric, totalDays)
	N := c.alphaLen
	for day := range totalDays {
		secOffset, err := rndO.Intn()
		if err != nil {
			mlog.FatalT(120, "error obtaining unique value", mlog.Err(err), mlog.At())
		} else {
			secOffset = secOffset % N
			if secOffset == 0 {
				secOffset++
			}
		}

		paramBooklet[day] = BiParametric{
			A: rndM.Intn(),
			B: secOffset,
		}
	}

	return paramBooklet
}

// generate a code book of secret words for use with Bellaso & Vigenère.
func (c *Caesarium) CompileWordBook(passwordLen int) []string {
	// @note if alphabet is EN, ES, DE or RU/UA we can promote the ASCII
	// to the foreign alphabet via transliteration. If it is Greek we
	// cannot because it has less runes than the source alphabet.
	canPromote := c.alphabet != cmn.ALPHA_DISK_GREEK

	var rndS IRandomizer
	if c.repeatable {
		rndS = NewRepeatableRand(c.date, c.userSeed, 0, c.alphaLen-1)
	} else {
		const NO_DIGITS = true
		rndS = NewTrueRand(0, c.alphaLen-1, canPromote, NO_DIGITS)
	}

	// the last day of this month
	totalDays := DaysInMonth(c.date)
	paramBooklet := make([]string, totalDays)
	for day := range totalDays {
		paramBooklet[day] = rndS.Runen(c.alphabet.Chars, passwordLen)
	}

	return paramBooklet
}

// generate a code book of secret A & B coefficients for use with Affine.
// A is limited by N because only coprimes up to N are valid.
// B is limited by N, values greater than get a modulo N applied, therefore
// we restrict to the basic range.
// Additionally ensure that B > 0 and A > 1 because an Affine with A=1 is
// nothing but a Caesar cipher with a key given by B.
// We do not need to specify A' because that is always derived from A & N.
func (c *Caesarium) CompileAffineBook() []TriParametric {
	// get the list of valid A coefficient coprimes for the known N
	ahelp := crypto.NewAffineHelper()
	validA := ahelp.ValidCoprimesUpTo(uint(c.alphaLen))

	// for the A coefficient, it is bound by N and A-coprimes
	var rndA IRandomizer
	if c.repeatable {
		rndA = NewRepeatableRand(c.date, c.userSeed, 0, len(validA)-1)
	} else {
		rndA = NewTrueRand(0, len(validA)-1, false, false)
	}

	// for the offset we want a different list but predictable
	offsetSeed := c.userSeed + extraOffsetSeed
	// for the B coefficient, expanded but internally the cipher applies modulo N
	var rndB IRandomizer
	if c.repeatable {
		rndB = NewRepeatableRand(c.date, offsetSeed, 1, c.alphaLen)
	} else {
		rndB = NewTrueRand(1, c.alphaLen, false, false)
	}

	// the last day of this month
	totalDays := DaysInMonth(c.date)
	paramBooklet := make([]TriParametric, totalDays)

	for day := range totalDays {
		// select an A coefficient from the valid pool for N=alphaLength
		coefA := validA[rndA.Intn()]
		// calculate its multiplicative inverse, the A' coefficient for decryption
		primeA, err := ahelp.ModularInverse(coefA, c.alphaLen)
		if err != nil {
			mlog.FatalT(caesarx.ERR_PARAMETER, "no multiplicative inverse",
				mlog.Err(err),
				mlog.At(),
			)
		}

		paramBooklet[day] = TriParametric{
			A: coefA,       // A coefficient - must be coprime of N
			B: rndB.Intn(), // B coefficient - N > B
			C: primeA,      // A' - multiplicative inverse of A
		}
	}

	return paramBooklet
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// get the last day of the month for the given time
func LastDay(t time.Time) time.Time {
	year, month, _ := t.Date()
	day := 0
	switch time.Month(month) {
	case time.April, time.June, time.September, time.November:
		day = 30
	case time.February:
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) { // leap year
			day = 29
		}
		day = 28
	default:
		day = 31
	}
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

// get the number of days in that month of the year.
func DaysInMonth(t time.Time) int {
	year, month, _ := t.Date()
	day := 0
	switch time.Month(month) {
	case time.April, time.June, time.September, time.November:
		day = 30
	case time.February:
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) { // leap year
			day = 29
		}
		day = 28
	default:
		day = 31
	}

	return day
}
