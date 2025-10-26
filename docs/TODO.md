# THINGS TO DO


## KNOWN BUGS

-[X] -variant caesar ignores -num option
-[X] cmd/tabularecta -demo GERMAN fails for key ß. `SpecialCases` deal with library errors.
-[X] tabularecta_caesar_test#57 Key ẞ fails for ẞ `SpecialCases`
-[X] caesar_test.go WithAlphabet() fails for Latin/Greek/Cyrillic but okay for English.
-[X] GH-010 sporadic Vigenere bin file decryption failure *to be verified* (since v1.1-0-RC3)
-[X] Affine decryption failure with `tests/testdata/input.bin` (since v1.1-0-RC4)

## ENHANCEMENTS

-[·] use Sequencer for Affine as well. v1.1.0-RC4 used in Binary
-[X] ciphers.IPipe in CaesarCipherCommand{} `IPipe` & `Pipe`
-[X] {cmn} allow built-in alphabets (namely German) have its  own equivalents to unicode.ToUpper/Lower(rune)
     and strings.ToUpper/Lower() by overriding pointers to standard functions. This allows
     for transparent handling rather than implement exception logic in the standard algorithm.
     See `cmn.SpecialCaseHandlers`.
-[ ] HTTP API module that uses command patterns to do encodings as a web service
-[ ] {cmd} generate a day-of-the-month encoding parameter table like ENIGMA that never changes for
     the same input seed. In GitHub there is a Go table formatter package.
-[ ] User configuration file to set preferred `-alpha`, `-variant` and the parameters for the current
     day if the *day-of-the-month* feature is implemented.
-[X] [Build RPM](https://infotechys.com/how-to-create-and-build-rpm-packages/) or perhaps with [workflow](https://www.spencersmolen.com/creating-an-automated-rpm-build-pipeline-using-github-actions/) **v1.1.2**

## TESTS

-[X] {ciphers} NewGroupingCommand (Trigrams/Quartets/Quintets) `Test_NgramFormatter`
-[ ] {util} HasUniqueRunes() used by Alphabet.Check
-[X] {util} IntersectInt
-[ ] {util} IsNotBlank()
-[ ] {cmn} Alphabet.Check, Alphabet.WithSpecialCase
-[X] {ciphers} chained alphabets

-[X] {affine} AffineHelper common coprimes, check if common
-[X] {affine} AffineEncoder, AffineDecoder
-[X] {affine} WithChain()

## REFACTORING ETC.

-[X] Re-Audit usage of string[from:to] because that refers to indices and does not
     work with multi-byte rune strings!. **Watch the source of a few bugs here**
-[ ] Use caesarx.CipherVariant enum instead of cmd/caesar/VariantID
-[·] Refactor using `internal` tree to export only what needs to be exported
-[ ] Error code for exit value reorganization so that it helps pinpoint errors
-[ ] Use predefined `ErrXXXXX` errors
-[ ] Custom errors `TabulaError`, `CipherError`, `ErrorInvalidValue`, `ErrorCliBadParam`
-[ ] Encrypt/DecryptBinaryFile in `tabula_caesar.go` and `affine_crypto.go` are the same
     but only need the `IKeySequencer` instance.

# Project Milestones

- 27.Sep.2025 v1.0 String message encryption with Caesar,Didimus,Fibonacci,Bellaso,Vigenere & Affine.
- 09.Oct.2025 v1.1-RC1 Text file encryption support for all ciphers
- 14.Oct.2025 Binary file encryption support for all ciphers except Affine.

