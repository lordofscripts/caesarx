/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Length of BIP39 sentence as enumeration.
 *-----------------------------------------------------------------*/
package bip39

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	// The number of words in the BIP39 mnemonic sentence
	Bip39Words12 Bip39Length = iota
	Bip39Words15
	Bip39Words18
	Bip39Words21
	Bip39Words24
)

/* ----------------------------------------------------------------
 *							L o c a l s
 *-----------------------------------------------------------------*/

var toString = map[Bip39Length]string{
	Bip39Words12: "Bip39Words12",
	Bip39Words15: "Bip39Words15",
	Bip39Words18: "Bip39Words18",
	Bip39Words21: "Bip39Words21",
	Bip39Words24: "Bip39Words24",
}

var toEnum = map[string]Bip39Length{
	// from their enumeration name
	"Bip39Words12": Bip39Words12,
	"Bip39Words15": Bip39Words15,
	"Bip39Words18": Bip39Words18,
	"Bip39Words21": Bip39Words21,
	"Bip39Words24": Bip39Words24,
	// from their integer values
	"12": Bip39Words12,
	"15": Bip39Words15,
	"18": Bip39Words18,
	"21": Bip39Words21,
	"24": Bip39Words24,
}

var toSize = map[Bip39Length]int{
	Bip39Words12: 12,
	Bip39Words15: 15,
	Bip39Words18: 18,
	Bip39Words21: 21,
	Bip39Words24: 24,
}

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type Bip39Length int

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// String returns the string version of the cipher name
// @implements fmt.Stringer interface
func (c Bip39Length) String() string {
	var name string = ""
	if c.IsValid() {
		name = toString[c]
	}

	return name
}

// Given an enumeration value return the actual BIP39 size
// (12..24) or -1 if it is an invalid enum value.
// @note GO enumerations are not quite perfect yet
func (c Bip39Length) ToSize() int {
	if val, ok := toSize[Bip39Length(c)]; !ok {
		return -1
	} else {
		return val
	}
}

// whether the contained value is within the valid range of the enumeration
func (c Bip39Length) IsValid() bool {
	if c == Bip39Words12 ||
		c == Bip39Words15 ||
		c == Bip39Words18 ||
		c == Bip39Words21 ||
		c == Bip39Words24 {
		return true
	}

	return false
}

// Convert takes any integer and tries to convert it to a Bip39Length
// enumeration value if it is within the range.
func (c Bip39Length) Convert(v int) (Bip39Length, error) {
	if val, ok := toEnum[strconv.Itoa(v)]; !ok {
		return 0, fmt.Errorf("cannot convert '%d' as Bip39Length", v)
	} else {
		return val, nil
	}
}

// Parse takes a cipher name string and attempts to parse it to
// convert it to the CipherVariant enumeration. Case insensitive.
func (c Bip39Length) Parse(str string) (Bip39Length, error) {
	var err error = nil
	var v Bip39Length = 0
	var ok bool

	if v, ok = toEnum[str]; !ok {
		// try case-insensitive map key search
		found := false
		for k, val := range toEnum {
			if strings.EqualFold(k, str) {
				v = val
				found = true
				break
			}
		}

		if !found {
			err = fmt.Errorf("invalid BIP39 name '%s'", str)
		}
	}

	return v, err
}

// Custom YAML unmarshalling of enumeration string to its numeric value.
func (c *Bip39Length) UnmarshalYAML(value *yaml.Node) error {
	// it should be a string
	var name string
	if err := value.Decode(&name); err != nil {
		return err
	}

	// it was a string, continue
	if v, ok := toEnum[name]; !ok {
		return fmt.Errorf("parse literal has invalid Bip39Length '%s'", name)
	} else {
		*c = v
		return nil
	}
}

// Custom YAML marshalling of enumeration, otherwise it appears as integer.
func (c Bip39Length) MarshalYAML() (any, error) {
	return toString[c], nil
}
