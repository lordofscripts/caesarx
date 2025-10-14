/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Cipher algorithm enumeration
 *-----------------------------------------------------------------*/
package caesarx

import "fmt"

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	NoCipher CipherVariant = iota
	CaesarCipher
	DidimusCipher
	FibonacciCipher
	BellasoCipher
	VigenereCipher
	AffineCipher
)

var cipherToString = map[CipherVariant]string{
	NoCipher:        "None",
	CaesarCipher:    "Caesar",
	DidimusCipher:   "Didimus",
	FibonacciCipher: "Fibonacci",
	BellasoCipher:   "Bellaso",
	VigenereCipher:  "Vigenere",
	AffineCipher:    "Affine",
}

var stringToCipher = map[string]CipherVariant{
	"None":      NoCipher,
	"Caesar":    CaesarCipher,
	"Didimus":   DidimusCipher,
	"Fibonacci": FibonacciCipher,
	"Bellaso":   BellasoCipher,
	"Vigenere":  VigenereCipher,
	"Affine":    AffineCipher,
}

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type CipherVariant uint8

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// String returns the string version of the cipher name
// @implements fmt.Stringer interface
func (c CipherVariant) String() string {
	_, name, _ := c.Convert(int(c))
	return name
}

// Convert takes any integer and tries to convert it to a CipherVariant
// enumeration value if it is within the range.
func (c CipherVariant) Convert(v int) (CipherVariant, string, error) {
	if val, ok := cipherToString[CipherVariant(v)]; !ok {
		return NoCipher, "", fmt.Errorf("cannot convert '%d' as CipherVariant", v)
	} else {
		return CipherVariant(v), val, nil
	}
}

// Parse takes a cipher name string and attempts to parse it to
// convert it to the CipherVariant enumeration.
func (c CipherVariant) Parse(str string) (CipherVariant, error) {
	if v, ok := stringToCipher[str]; !ok {
		return NoCipher, fmt.Errorf("invalid cipher name '%s'", str)
	} else {
		return v, nil
	}
}
