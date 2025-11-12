/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * YAML.v3 (un)marshallers for Rune type. It makes the rune type,
 * which in GO is an alias for int32, to be rendered as a character
 * rather than as a number.
 *-----------------------------------------------------------------*/
package prefs

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ yaml.Unmarshaler = (*Rune)(nil)
var _ yaml.Marshaler = (*Rune)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// Rune represents a rune (UTF-8 character) and implements YAML
// serialization/deserialization.
type Rune rune

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// MarshalYAML converts the rune to a YAML character.
func (r Rune) MarshalYAML() (any, error) {
	return string(r), nil
}

// UnmarshalYAML converts a YAML node back to a rune.
func (r *Rune) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode {
		return fmt.Errorf("expected scalar node, got %v", node.Kind)
	}
	runeString := node.Value
	*r = Rune([]rune(runeString)[0]) // Handle multi-byte characters properly
	return nil
}

/*
// UnmarshalYAML converts a YAML character back to a rune.
func (r *Rune) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var char string
	if err := unmarshal(&char); err != nil {
		return err
	}
	*r = Rune([]rune(char)[0]) // Handle multi-byte characters properly
	return nil
}

func demo() {
	All := []Rune{'😊', '👍', 'Ш', 'Σ', 'ẞ', 'A', 'É'}

	for i, chr := range All {
		fmt.Printf("#%d Item: %c\n", i+1, chr)
		// Serialization
		example := Example{Char: chr} // Using a multi-byte character
		data, err := yaml.Marshal(example)
		if err != nil {
			panic(err)
		}
		fmt.Println("\tSerialized YAML:")
		fmt.Println("\t", string(data))

		// Deserialization
		//yamlData := "char: '👍'" // UTF-8 multi-byte character in YAML
		var deserializedExample Example
		if err := yaml.Unmarshal([]byte(data), &deserializedExample); err != nil {
			panic(err)
		}
		fmt.Printf("\tDeserialized Char: %c\n", deserializedExample.Char)
	}

}
*/
