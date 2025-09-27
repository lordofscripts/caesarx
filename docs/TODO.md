# THINGS TO DO


## KNOWN BUGS

-[X] -variant caesar ignores -num option
-[X] cmd/tabularecta -demo GERMAN fails for key ß. `SpecialCases` deal with library errors.
-[X] tabularecta_caesar_test#57 Key ẞ fails for ẞ `SpecialCases`
-[X] caesar_test.go WithAlphabet() fails for Latin/Greek/Cyrillic but okay for English.

## ENHANCEMENTS

-[ ] use Sequencer for Affine as well.
-[X] ciphers.IPipe in CaesarCipherCommand{} `IPipe` & `Pipe`
-[X] {cmn} allow built-in alphabets (namely German) have its  own equivalents to unicode.ToUpper/Lower(rune)
     and strings.ToUpper/Lower() by overriding pointers to standard functions. This allows
     for transparent handling rather than implement exception logic in the standard algorithm.
     See `cmn.SpecialCaseHandlers`.
-[ ] HTTP API module that uses command patterns to do encodings as a web service
-[ ] {cmd} generate a day-of-the-month encoding parameter table like ENIGMA that never changes for
     the same input seed. In GitHub there is a Go table formatter package.

## TESTS

-[X] {ciphers} NewGroupingCommand (Trigrams/Quartets/Quintets) `Test_NgramFormatter`
-[ ] {util} HasUniqueRunes() used by Alphabet.Check
-[X] {util} IntersectInt
-[ ] {util} IsNotBlank()
-[ ] {cmn} Alphabet.Check, Alphabet.WithSpecialCase
-[ ] {ciphers} chained alphabets

-[X] {affine} AffineHelper common coprimes, check if common
-[X] {affine} AffineEncoder, AffineDecoder
-[X] {affine} WithChain()

## REFACTORING ETC.

-[X] Re-Audit usage of string[from:to] because that refers to indices and does not
     work with multi-byte rune strings!. **Watch the source of a few bugs here**
-[ ] Refactor using `internal` tree to export only what needs to be exported

