# GO Caesar X Encryption

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/lordofscripts/caesarx)
[![Go Report Card](https://goreportcard.com/badge/github.com/lordofscripts/caesarx?style=flat-square)](https://goreportcard.com/report/github.com/lordofscripts/caesarx)
![Build](https://github.com/lordofscripts/caesarx/actions/workflows/go.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/lordofscripts/caesarx.svg)](https://pkg.go.dev/github.com/lordofscripts/caesarx)
[![GitHub release (with filter)](https://img.shields.io/github/v/release/lordofscripts/caesarx)](https://github.com/lordofscripts/caesarx/releases/latest)
[![License: CC BY-NC-ND 4.0](https://img.shields.io/badge/License-CC_BY--NC--ND_4.0-lightgrey.svg)](https://creativecommons.org/licenses/by-nc-nd/4.0/)


![](./assets/caesarx_header.jpg)


The ultimate application for modern-day usage of the ancient Caesar cipher and some of its variants. While Caesar-type ciphers are no match and certainly not a replacement for modern-day cryptography, there are situations like in family and friend circles, even for games, where it is handy to encrypt some text that can —should disaster strike— be able to decrypt it by hand. And sometimes you just need to hide in plain view if it is not that critical.

**Go Caesar X** is a pure GO implementation of the famous ancient Caesar cipher and some of its known variants used in the XIX century plus a couple of variations I created just for fun. The truth is, most Caesar implementations out there, are still based on the limited ASCII alphabet (A-Z). Those are unsuitable for modern-day usage because it does not work with symbols and multi-byte characters used in other languages like Spanish, German, Greek or the Cyrillic alphabets. Well, Go·Caesar·X does!

Better yet, it does so by using concentric discs or slave alphabets. With plain Caesar encoding you can see numbers and spaces; therefore, with some skill you can start attacking it and predict the encryption key. My implementation —while no substitute for today's encryption algorithms— makes those attacks more difficult.

This fun application and library is merely an educational experiment I did for fun in my free time. Just because I love cryptography and have been implementing modern-day encryption algorithms in several languages like C, C++, C#, Java and now in Go.

### Features:

* Implements plain Caesar, Didimus, Fibonacci, Bellaso, Vigenère and Affine.
* Includes several built-in modern-day alphabets: English (plain ASCII), Latin/Spanish, German, Greek and Cyrillic.
* Supports custom alphabets
* Does not break with Unicode multi-byte characters, specially designed for this!
* Preserves upper/lowercase, no need to extend the alphabet.
* Supports custom casing rules (more details in the technical document).
* Supports symbols and chained alphabets.
* The command-pattern based API allows piping results into other commands. For example double-algorithm encryption or output grouping.
* Supports grouping encrypted output in groups of 2/3/4/5 characters for CLI text encryption (message-based).
* Supports **text file encrytion** in all algorithm modalities! You are not limited to short messages anymore.
* Supports **binary file encryption** for all algorithms. Now you can encrypt images and other binary data.
* Lots of test cases included

|     | Show your support   |
| --- | :---: | 
| [ ![AllMyLinks](./assets/allmylinks.png)](https://allmylinks.com/lordofscripts)      | visit <br> Lord of Scripts&trade; <br> on [AllMyLinks.com](https://allmylinks.com/lordofscripts)                  |
| [ ![Buy me a coffee](./assets/buymecoffee.jpg)](https://allmylinks.com/lordofscripts)|  buy Lord of Scripts&trade; <br> a Capuccino on <br>[BuyMeACoffee.com](https://www.buymeacoffee.com/lostinwriting)| 

#### License

You are **not** granted permission to use this library, application or derivative works for profit or commercial purposes. Just to let you know... Read the [LICENSE](../LICENSE.md) for more details.

### Installation

These instructions depend whether you are going to use this fine module as an
end-user application or integrate the functionality in your **free** software.
Respect the license please...

#### For End-Users

To install the `tabularecta`, `caesarx` and `affine` executables in your system:

`go get github.com/lordofscripts/caesarx@latest`

Or you can install the Debian package (recommended):

`sudo apt-get install go-caesarx******.deb` 

On **Linux** you have the `caesarx`, `tabularecta` and `affine` as main CLI apps, 
with `didimus`, `fibonacci`, `bellaso` and `vigenere` probably appearing as
*soft symbolic links* to `caesarx`.

On **Windows** you have the `caesarx.exe`, `tabula.exe` and `affine.exe` 
executable files.

##### Usage

Both the `tabularecta` and the `caesarx` are CLI (command-line interface) application with some common options to get you started exploring this exciting world of ancient ciphers. Additionally there are special aliases to `caesarx` called `affine`, `bellaso` and `vigenere` that are pre-configured for those algorithms; thus sparing you from having to specify the algorithm through the CLI.

The `-demo` CLI option executes a demonstration of what it can do. Likewise, the `-help` CLI option shows you all the parameters expected by the application(s).

The `-alpha ALPHABET` CLI option lets you specify the target language/alphabet, that is very important for your text to be properly encrypted. Valid values are: `english, latin, german, greek, cyrillic` where `spanish` is the same as `latin`. It defaults to English which is the plain ASCII A-Z without accented character support. I prepared a [LANGUAGES](./LANGUAGES.md) page with more details.

The `-list` option in CaesarX lists all supported cipher variants (not present in the `affine` application).

The `-F` parameter tells the application that the free argument at the end
(by encryption and decryption actions) is a *filename*  rather than an
explicit (plain/cipher) text string. (as of v1.1) This applies only to
*encode* and *decode* operations of all applications except `tabularecta`.

To encrypt a plain text file `secret.txt` use:

>
> program_name {parameters} [options] -F secret.txt
>

the output file would be `secret_txt.EXT` where EXT is any of `cae, did, fib, bel, vig`
depending on the chosen algorithm/variant.

Similarly, to decrypt `ciphered_txt.EXT` to a plain text `plain.txt`
you would use two free arguments instead of one:

>
> program_name {parameters} [options] -F ciphered_txt.EXT plain.txt
>

#### For Integrating in your own FREE software

The usual:

`go get github.com/lordofscripts/caesarx`

And then RTFM :) plenty of info here and the source code is documented in the 
[Pkg Info](https://pkg.go.dev/github.com/lordofscripts/caesarx).

### The Ciphers Explained

As I indicated, these are *ancient* and XIX century ciphers. Even in the 
ancient era of Julius Caesar there was a need for encrypted communication. 
Ever since then the curious minds in the world have been fascinated by it, 
and brilliant minds have developed them over the centuries.

Do you remember the famous German Enigma machine that kept the allied forces 
on their toes until they cracked it? I also implemented a modern-day Enigma 
library in Go, but it isn't released yet.

Now, before further ado, read more about the ciphers supported by this library 
and application. It is up to you to put them in practice. Me and my family use 
it, and soon I will get my friends too. Simple solutions for modern-day problems.

These are the ciphers supported by my application. Please read the appropriate 
document to know about their strengths, weaknesses and how they differ from 
other implementations.

* Plain [Caesar](./CIPHER_CAESAR.md) cipher used by Julius Caesar, the Roman Emperor over 2000 years ago.
* [Didimus](./CIPHER_DIDIMUS.md) cipher is a polysyllabic variation of Caesar
* [Fibonacci](./CIPHER_FIBONACCI.md) cipher is another polysyllabic variation of Caesar I came up with for fun.
* [Bellaso](./CIPHER_BELLASO.md) cipher is a repeated-key cipher based on a secret word or phrase which builds upon the Caesar cipher.
* [Vigenère](./CIPHER_VIGENERE.md) cipher is an auto-key variation of the Bellaso cipher
* [Affine](./CIPHER_AFFINE.md) cipher is similar, despite it being invented much later, Caesar is a variation of Affine.

#### Common Concepts

All of these ciphers are based on an alphabet that determines what gets encoded. If a symbol or character is not in the reference alphabet, it is passed as-is to the output.

They are based on assigning a numerical value to each letter, the value represents the number of shifted positions in the *transliterated alphabet*. So, if we have the plain English alphabet 
`ABCDEFGHIJKLMNOPQRSTUVWXYZ` we have N=26 (number of characters), with the *shift value* starting
with **zero** and incrementing with each of the letters *in that precise order*. Therefore, A=0,
B=1,C=2 and so on until Z=25 in the reference **English** alphabet. 

A particularity is that the first letter (in the English alphabet the letter "A") has a shift of
zero positions. That means that the output is exactly the same as the input. Therefore, you can't
really encrypt with *shift=0* being letter "A" in English, Spanish and German alphabets, "Α" in 
Greek and "А" in cyrillic. Well, it is a coincidence that they look the same, but their Unicode
values are not the same!

You can define the alphabet in another order, or mixing letters and numbers, but it is important
that you use exactly the same order of characters during encryption and decryption.

***
Copyright &copy;2025 Lord of Scripts


