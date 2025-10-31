/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Cipher algorithm enumeration with YAML (un)marshalling.
 *-----------------------------------------------------------------*/
package caesarx

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

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

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ yaml.Unmarshaler = (*CipherVariant)(nil)
var _ yaml.Marshaler = NoCipher

/* ----------------------------------------------------------------
 *							L o c a l s
 *-----------------------------------------------------------------*/

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
// convert it to the CipherVariant enumeration. Case insensitive.
func (c CipherVariant) Parse(str string) (CipherVariant, error) {
	var err error = nil
	var v CipherVariant = NoCipher
	var ok bool

	if v, ok = stringToCipher[str]; !ok {
		// try case-insensitive map key search
		found := false
		for k, val := range stringToCipher {
			if strings.EqualFold(k, str) {
				v = val
				found = true
				break
			}
		}

		if !found {
			err = fmt.Errorf("invalid cipher name '%s'", str)
		}
	}

	return v, err
}

// Custom YAML unmarshalling of enumeration string to its numeric value.
func (c *CipherVariant) UnmarshalYAML(value *yaml.Node) error {
	// it should be a string
	var name string
	if err := value.Decode(&name); err != nil {
		return err
	}

	// it was a string, continue
	if v, ok := stringToCipher[name]; !ok {
		return fmt.Errorf("parse literal has invalid CipherVariant '%s'", name)
	} else {
		*c = v
		return nil
	}
}

// Custom YAML marshalling of enumeration, otherwise it appears as integer.
func (c CipherVariant) MarshalYAML() (interface{}, error) {
	return cipherToString[c], nil
}

func MaxCipher() int {
	return int(AffineCipher)
}
