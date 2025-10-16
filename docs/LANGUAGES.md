# Built-in Languages

[![Go Reference](https://pkg.go.dev/badge/github.com/lordofscripts/caesarx.svg)](https://pkg.go.dev/github.com/lordofscripts/caesarx)
[![GitHub release (with filter)](https://img.shields.io/github/v/release/lordofscripts/caesarx)](https://github.com/lordofscripts/caesarx/releases/latest)
[![License: CC BY-NC-ND 4.0](https://img.shields.io/badge/License-CC_BY--NC--ND_4.0-lightgrey.svg)](https://creativecommons.org/licenses/by-nc-nd/4.0/)
[![Go Report](https://goreportcard.com/badge/github.com/lordofscripts/caesarx)](https://goreportcard.com/report/github.com/lordofscripts/caesarx)

![](./assets/caesarx_header.jpg)


## Alphabets

The correct term is "Alphabets" because several languages share a similar, if not the same,
alphabet. While the ancient Roman alphabet didn't have 26 letters as we do today in English,
the modern English alphabet is an ASCII code. That means it can be represented with a single
byte of information with a value between 0..255 (0x00 to 0xFF in hexadecimal notation).

The built-in **English** alphabet is as you may expect, entirely ASCII-based. And while in
the GO language strings are Unicode (a letter/symbol may be 1 to 4 bytes long in Unicode),
an ASCII string is easy to use. Most implementations in the wild rely on that and would break
or fail as soon as you feed a multi-byte character like "ß", "Л" or "Σ". Therefore, my 
implementation not only relies on GO's Unicode strings, but it handles them as runes to
properly process the text to be encoded/decoded.

Keep in mind that CaesarX while CaesarX preserves case, it stores its built-in alphabets
in **uppercase**.

But I won't bore you with more details, let's get to know the built-in alphabets in CaesarX.
The alphabet names given below are what the application expects from you as the `-alpha` 
parameter:

* `english` (ASCII)
* `latin` which is the same as `spanish` (UTF8)
* `german` (UTF8)
* `greek` (UTF8)
* `cyrillic` is an alias for `ukranian` and `russian` (UTF8)
* `binary` a universal alphabet that serves primarily to work with non-text files
  (images, PDF documents, binary files in general, etc.) that uses an *alphabet*
  composed of bytes `0x00` to `0xFF`.

All the ciphers in this module other than the *plain standard Caesar* use one of
the above named alphabets as **primary** *reference alphabet* and the Extended Numeric
**slave** *reference alphabet* `numbers+` which contains the decimal digits, a few 
commonly-used symbols and the space character. With that, your cipher is stronger
against the common thieves or casual non-tech bystanders. Why? because a space on
the encrypted text is NOT a space in the plain message, and the common symbols hide
tell-tale signs like "a number or email address follows!".

### English

**Alphabet:** ABCDEFGHIJKLMNOPQRSTUVWXYZ
**Number of items:** 26 runes (26 bytes)
**Type:** ASCII (single-byte per character)
**Special casing rules:** None.

Say "I love cryptography!". Here is the corresponding [Tabula Recta](./data/english_tabula.txt).

### Spanish

**Alphabet:** ABCDEFGHIJKLMNÑOPQRSTUVWXYZÁÉÍÓÚÜ
**Number of items:** 33 runes (40 bytes)
**Type:** UTF8 (multi-byte per character)
**Special casing rules:** None.

It is like the English alphabet but contains accute-accented
vocals á, é, í, ó and ú, the umlauted ü ("vergüenza") and the
well known Spanish N with tilde ñ.

BTW, *you can still use this alphabet to encrypt your English
text*, in fact, it will make it slightly stronger against 
brute force and character frequency attacks because the accented
characters are part of the primary alphabet.

Now say "Amo la criptografía!". Here is the corresponding [Tabula Recta](./data/latin_tabula.txt).

### German

**Alphabet:** ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÜẞ
**Number of items:** 30 runes (34 bytes) 
**Type:** UTF8 (multi-byte per character)
**Special casing rules:** Yes.

It has several characters with umlaut and the *eszet* symbol that **looks like**
the Greek *lowercase beta* but it is not the same Unicode code point/value!

Now, in all programming languages out there, converting the lowercase "ß" to
uppercase would yield "SS" which is TWO letters. That's an abomination obviously
not invented by programmers. Luckily, GO was designed by smart people that opted
to think differently. I won't go in details or discussions, but in the GO programming
language, the lowercase "ß" eszet has a single-letter (single-rune) uppercase
value "ẞ" that **looks like** the same lowercase, but keen eyes would notice it 
is fatter and slightly different, in fact, it has its own Unicode value different
from its lowercase equivalent and from the similar Greek lowercase Beta.

**NOTE:** Despite that, I had to take extra measures and devise this *Special Case Handling*
in my library and application because at least up to GO v1.25 converting the lowercase
"ß" into its uppercase "ẞ" via the standard library `strings.ToUpper()` `strings.ToLower()` 
and `unicode.ToUpper()` are broken and do not convert the *Eszet* correctly or consistently
as the rest of the library. I spent quite some time debugging the mysterious problem until
I found it. I filed the bug to the GO code repository.

Can you say "Daß liebe hübschen Mädchen"? I haven't spoken German in many years, so I
lost practice with German declinations.  Here is the corresponding [Tabula Recta](./data/german_tabula.txt).

### Greek

**Alphabet:** ΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩ
**Number of items:** 24 runes (48 bytes)
**Type:** UTF8 (multi-byte per character)
**Special casing rules:** None.

Say "Λατρεύω την κρυπτογραφία!". Here is the corresponding [Tabula Recta](./data/greek_tabula.txt).

### Cyrillic

**Alphabet:** АБВГДЕËЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ
**Number of items:** 33 runes (66 bytes)
**Type:** UTF8 (multi-byte per character)
**Special casing rules:** None.

The letters in the cyrillic alphabet are shared by the Ukranian and Russian languages and
are used (some, maybe not all) in the Serbian language.

Say "Я люблю криптографию!" just don't tell Puttin about it! Here is the corresponding [Tabula Recta](./data/cyrillic_tabula.txt).

### Numbers & Punctuation

There are several built-in alternate alphabets that contain important characters for
modern-day communications. These are suitable as *Slave Discs* to extend the basic
letter alphabets introduced above (English, Spanish, etc.)

* Standard numbers `0123456789` identified by the `numbers` alphabet handle.
* Eastern numbers `٠١٢٣٤٥٦٧٨٩`  which are used in the Hindi language. It is a right-to-left system though. It goes by the `numbers_east` alphabet handle.
* Extended numeric `0123456789 #$%+-@` goes by the `numbers+` handle. It has the commonly used arabic numbers used in most of the world, plus an essential list of symbols that would make your encrypted messages more difficult to figure out, thus leaving no tell-tale signs of what kind of information is encrypted in the message. This library was meant for modern-day usage. It also includes the SPACE character, thus making it difficult to know where the word boundaries of the encrypted text are.
* Symbols/Punctuation disk `¡!\"#$%&'()*+,-./0123456789:;<=>¿?@[]` contains most symbols for common languages. It may be a useful slave disk.

 Here is the corresponding [Tabula Recta](./data/numeric_tabula.txt).

### Binary

**Alphabet:** ASCII codes 0..255
**Number of items:** 256 bytes
**Type:** Binary
**Special casing rules:** None.

A special alphabet that treats input as mere 8-bit bytes. It is useful for 
any binary files such as images, executables, etc. When you are going to
encrypt/decrypt binary files with `affine` or `caesarx` and its aliases,
use the CLI options `-alpha binary -F` and it will treat the filename
arguments as *binary* files rather than text.

***
Copyright &copy;2025 Lord of Scripts
