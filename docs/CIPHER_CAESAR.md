# Caesar Cipher

[![Go Reference](https://pkg.go.dev/badge/github.com/lordofscripts/caesarx.svg)](https://pkg.go.dev/github.com/lordofscripts/caesarx)
[![GitHub release (with filter)](https://img.shields.io/github/v/release/lordofscripts/caesarx)](https://github.com/lordofscripts/caesarx/releases/latest)
[![GitHub License](https://img.shields.io/github/license/lordofscripts/caesarx)](https://github.com/lordofscripts/caesarx/blob/master/LICENSE)
[![Go Report](https://goreportcard.com/badge/github.com/lordofscripts/caesarx)](https://goreportcard.com/report/github.com/lordofscripts/caesarx)

![](./assets/caesarx_header.jpg)


## History

The **Caesar cipher** was invented by Julius Caesar (yes, the Roman emperor) to
convey secret messages to his generals in the field. At the time the Roman
alphabet had 23 letters (ABCDEFGHIKLMNOPQRSTVXYZ). And they used Roman
numerals for their numbers (unlike our decimal digits), so there was no
problem in transmitting secret messages with letters and numbers.

The Caesar cipher is rudimentary for today's standards —yet still useful— but
we ought to remember that in the era of the Roman Empire, education was
a luxury few could afford, just the elites. Even a plain text may have
looked like a cipher to someone who cannot read.

While it was an ancient cipher, you can see in [Duke's library](https://people.duke.edu/~ng46/collections/crypto-disk-strip-ciphers.htm)
that even during World War I the Bulgarians used a Cyrillic Caesar cipher disc
and the U.S. Army used a similar disc.

The so called *Rot-13* encoding used in some applications nowadays, is nothing more
than a Caesar with a shift=13. However, it is usually based on the English alphabet only.

## Strengths & Weaknesses

Strengths:
* Very simple to understand
* Quick encryption and decryption, not even a computer is needed

Weaknesses:
* Can be easily broken with frequency analysis due to the statistical distribution of letters.
* Extremely vulnerable to brute-force attacks. Plain Caesar has only 25 possible keys (A does not count because it is a *shift=0*). Just the spaces serve as tell-tale sign unless you group the encrypted output in trigram/quartet/quintets. 
* For a given key, the output character is always the same.
* There are 25 possibilities

## Encryption & Decryption

The Plain Caesar cipher is a simple substitution of one letter for another depending
on the chosen Caesar key. The Caesar key is a **single** letter that must be present
in the encoding alphabet. That key determines the number of shift positions.

You can learn more by watching the [Cryptography: Caesar cipher with shift](https://www.youtube.com/watch?v=F6vBdvt8Ctw) on YouTube. Or you can visualize and shift a Caesar [cipher disk online](https://computerscienced.co.uk/site/caesar-cipher-wheel/caesar-cipher/).

Once you agreed on the single-character key, all you have to do is take the reference
alphabet and shift it that many positions. It's easy to do with just pen and paper.
However, helper devices have been used throughout history such as:

* The Caesar disk with two concentric discs. The inner disc is rotated the amount of shift positions and to encrypt you lookup letters in the outer disc and convert it to the letter in the inner disc. For decryption you do the same shift but look up the ciphered letter in the inner disc and retrieve the decrypted letter from the outer disc. Alternatively, you can encode from inner to outer and decode from outer to inner, whatever workflow suits you best.
* Use a Tabula Recta made of a square matrix with the reference alphabet atop (horizontally) and vertically. Here you have all possible key combinations in a single table. In the [LANGUAGES](./LANGUAGES.md) document you can browse the Tabula Rectas for each of the built-in alphabets.
* A simple transliteration table with the plain alphabet in the top row accompanied by the numbers 0..n and in the lower (2nd) row the same alphabet but shifted K positions.

Below is a *transliteration table* (character substitution) based on the **English** built-in
alphabet and a Caesar *Key* "G". From the top row you can see that the "G" character represents
a *shift* by six (6) positions.

```
	0                   1                   2
	0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5
	A B C D E F G H I J K L M N O P Q R S T U V W X Y Z
	- - - - - - - - - - - - - - - - - - - - - - - - - - 
	G H I J K L M N O P Q R S T U V W X Y Z A B C D E F
```

For that Key you can do the exercise of encoding "Disability 2025" and you should end up with "Joyghoroze 2025".

## Using it with GoCaesarX

These are the CLI options you need for using the Plain Caesar cipher:

* `-alpha spanish` to select your built-in language, English is the default if you do not provide this option.
* `-variant caesar` if you want to be explicit, although that is the default.
* `-key <LETTER>` specifying a **single** letter from the *chosen alphabet* as your shift key. It is case-insensitive.
* `-d` if your last parameter is a text to be *Decoded/Decrypted*, if not given the last parameter in the CLI is the text to be *Encoded/Encrypted*.
* (optional) `-ngram <2|3|4|5>` if you want to reformat the cipher output in groups of 2,3,4 or 5 characters in order to disguise the tell-tale space character. Only for encryption.



