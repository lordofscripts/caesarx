# D E V E L O P E R

As of v1.1.0RC5 the application checks for a minimum GO version v1.20.
The reason being that in that version Math.Rand()'s behavior changed
and we rely on that change to generate proper temporary filenames.

## OS: Unix/Linux/Darwin

Nothing to be reported, just disappointed to learn that MacOS is a
case-insensitive OS.

## OS: Windows

**NOTE:** *Keep in mind that the repository is case-sensitive; therefore, 
make sure you use the same nomenclature. Above all, do not change the name
of files!.

To build on windows CLI:

 * `go build -buildmode=exe -o bin/gocaesar.exe -gcflags all=-N ./cmd/cli/main.go`
 * Without `-gcflags all=-N` it will produce a non-executable output!

## Showcase Concepts

For Go newbies there are a few interesting things done in this module,
that beyond the basic, may be worth noting:

* Using the GO `flag` package with custom flags: `cmd.rune_flag.go`
* Modular `main()` which makes it look orderly
* Avoiding pitfalls of indexing strings by byte instead of runes
* using interfaces
* Command Pattern with Piped commands (chained commands)
* Custom application logging with `app/mlog.go` which is my improved
   version of `log/slog`
* Test cases like `Test_Affine_Exit()` which exercise various CLI execution
  parameter combinations to check the application return value. You no
  longer need to run the CLI manually prior to every release.
* Skipping tests conditionally. `GITHUBLOS=true go test ./...` because
  the CLI application exit value test cannot be resolved on GitHub server
  because I don't know where the executable is in that server build.
  See `tabularecta_vigenere_test.go, tabularecta_caesar_test.go, go.yml
* Custom YAML/v3 (de)  serialization.
* Using YAML to properly serialize an enumeration `caesarx.CipherVariant`
  and a rune as a character instead of number `prefs.Rune`
  
More odd stuff:

* `Makefile` get absolute path of current makefile to get project's BIN dir   

## Built-in Languages

CaesarX comes with several [predefined alphabets](./LANGUAGES.md) that are suited for most
commonly used languages.

When defining new *built-in* alphabets keep in mind that the same structure
should be used in all of them. Each alphabet knows about *letters* such as
A..Z in English, and *vowels* such as AEIOU in English.

Many non-English languages have special characters, some are letters with 
special adornments, and others are *vowels with diacriticals*.

When we define a built-in language, the alphabet string first contains the
letters and letters with adornments followed by the vowels with diacriticals.

In the letter section we place a letter with adornment right next to the
corresponding letter without adornment, i.e. "MNÑO" as you can see there
there is an "N" followed by a "N tilde".

At the end of the alphabet string we find all vowels with diacriticals. The
normal vowels are in the letter section.

A clear example can be appreciated with the Czech alphabet: `ABCČDĎEFGHIJKLMNŇOPQRŘSŠTŤUVWXYÝZŽÁÉÍÓÚĚŮ` 
where you can see that the letters C,D,N,R,S,T,Y and Z have both a regular
and adorned version next to each other, in the *same order* as the **official
language**. And the last part contains the vowels with diacriticals grouped
by diacritical form (acute, grave, etc.)

### Languages/Alphabets and dealing with Files

As of `v1.0` CaesarX supports file encryption for both *text* files (v1.1.0-RC1)
and *binary* files (v1.1.0-RC4). This implies using the `-F` CLI argument **and**
using the `-alpha {ID}` argument to select the alphabet/language. If the `{ID}`
is set to `binary` the it will assume the files are binary files.

For developers, keep in mind that instantiating any of the *ciphers* can be done
with any language. But when calling the `EncryptBinFile` or `DecryptBinFile`
methods, the command instance will *ignore* the current human language and
use its own *Binary Tabula* instead of the language chosen for the instance.
This means, you don't have to create a specific instance to deal with binary
files.

## Package Building

The `Makefile` now has build targets to build DEB and RPM packages. The
local build paths are specified in the Makefile which currently default to:

* `~/Develop/Distrib/Build/caesarx/rpmbuild` for RPM
* `~/Develop/Distrib/Build/caesarx/DEBIAN` for DEB

That is for manual creation of packages. I am currently (v1.1.1) trying to
get the [Packaging Workflow](../.github/workflows/packaging.yml) to work. It
still gives error on GitHub.

### DEB Package

I use **Debian** so it is my native build environment.

> make debian

and the resulting Debian package can be found at `~/Develop/Distrib/Build/caesarx/`.

### RPM Package

I use **Debian** so it is my native build environment. In order to build RPM packages
on Debian you must install the `rpmbuild` package and rebuild its empty database so
that it doesn't give SQLite errors during build:

> sudo apt update
> sudo apt install rpmbuild
> rpmbuild --version
> sudo rpm --rebuilddb

The RPM package can be build with:

> make rpmclean; make rpm

and the resulting RPM package can be found at the directory
`~/Develop/Distrib/Build/caesarx/rpmbuild/BUILD/RPMS/x86_64/`.

*NOTE: I don't have a RedHat/Fedora system at my disposal; therefore, I am
currently unable to test the actual installation of the RPM package*

To list the files that would be installed by the RPM package:

> rpm -qlp ../RPMS/x86_64/caesarx-1.1.1-1.x86_64.rpm

which currently lists:

```
/usr/bin/affine
/usr/bin/bellaso
/usr/bin/caesarx
/usr/bin/didimus
/usr/bin/fibonacci
/usr/bin/tabularecta
/usr/bin/vigenere
/usr/lib/.build-id
/usr/lib/.build-id/14
/usr/lib/.build-id/14/e9dec8dd643856378523f48095d1d6e0324d8a
/usr/lib/.build-id/1e
/usr/lib/.build-id/1e/51cf2bef833e85d783074902000226682988f9
/usr/lib/.build-id/b4
/usr/lib/.build-id/b4/8fc2fc3f8436637d5335e3d65897c0a272dbe7
/usr/share/licenses/caesarx-1.1.1
/usr/share/licenses/caesarx-1.1.1/LICENSE.md
/usr/share/man/man1/affine.1.gz
/usr/share/man/man1/caesarx.1.gz
```

## Other 

### Debugging

Set the `LOG_LEVEL` environment variable to any of trace, debug, info, warn,
error or fatal. Log output will appear on `stderr`.

### Skipping Tests

Certain working tests will not be executed on GitHub servers because the
location of the app binary is unknown. For this see Issue #11 on this 
repository. I modified the `go.yaml` GitHub workflow file like this:

>
>    - name: Test
>      run: GITHUBLOS=true go test -v -coverprofile=profile.cov ./...
>

Then on the susceptible tests (those using the built executable) I added
code to skip them:

>
>	// @note We set this on go.yml so that this test is SKIPPED on GitHub servers
>	if os.Getenv("GITHUBLOS") != "" {
>		t.Skip("Skipping working test due to missing executable")
>	}
>

## Quirks

### The German Es-Tzet 'ß' Conondrum

There is a particular quirk with 'ß' that make us take a guess sometimes.
In at least 5 other programming languages, 'ß' is lowercase while 'SS' is
its uppercase counterpart. The problem there is that you rely on runes,
that uppercase version has **two** runes instead of one. That would be a
problem in how we handle things here.

On the other hand, in GO the 'ß' character is both the lowercase **and**
uppercase version. Good for the algorithm because it handles individual
runes (ONE character). 

Thus in German messages be careful how you type the Es-Tzet character 
because the lowercase version `ß` looks **almost** exactly the same as the
uppercase version `ẞ` and in my built-in German alphabet I had inadvertently
typed it as lowercase within the UPPERCASE alphabet! Spent quite some
time debugging that!

Special handling is done in all (German and custom alphabets) containing
the 'ß' character. Also notice that it is NOT the same as the
Greek 'β' (lowercase Beta) which has 'B' as its uppercase equivalent.

BTW the same applies for some (sometimes many) characters in other
languages, for example the Cyrillic alphabet. It's a mine field!

```go
package main

import (
 "fmt"
"unicode"
"strings"
)

func main() {
  const ch rune = 'ß' // typed with Lowercase
  const chCL rune = 'ß' // typed with Lowercase
  const chCU rune = 'ẞ' // typed with Uppercase

  fmt.Println("The German ß ambiguity when dealing with Runes") 

  chL := unicode.ToLower(ch)
  fmt.Printf("LC(%c) isLower %t\n", chL, unicode.IsLower(chL))
  fmt.Printf("LC(%c) isUpper %t\n", chL, unicode.IsUpper(chL))
  fmt.Printf("%c isLower %t (should be TRUE)\n", chCL, unicode.IsLower(chCL))

chU := unicode.ToUpper(ch)
  fmt.Printf("LC(%c) isLower %t\n", chU, unicode.IsLower(chU))
  fmt.Printf("LC(%c) isUpper %t\n", chU, unicode.IsUpper(chU))
  fmt.Printf("%c isUpper %t (should be TRUE)\n", chCL, unicode.IsUpper(chCU))
  // Now the issue of ß in strings.ToUpper()
  const SL = "ß"
  const SU = "ẞ"
  fmt.Printf("Must be true: ToUpper(%s) %t\n", SL, strings.ToUpper(SL) == SU) // FAILS
  fmt.Printf("Must be true: ToLower(%s) %t\n", SU, strings.ToLower(SU) == SL)

converted := unicode.ToUpper(chCL)
  if converted != chCU { // FAILS TOO
        fmt.Printf("GO std lib ToUpper failure: %c != %c\n", converted, chCU)
  }
  converted = unicode.ToLower(chCU)
  if converted != chCL {
        fmt.Printf("GO std lib ToLower failure: %c != %c\n", converted, chCL)
  }
}
```

### ONLINE TOOLS

For some of the standard algorithms used here, you can easily obtain
test data at [Cryptii](https://cryptii.com/) by choosing the *Caesar Cipher*
using the same alphabet defined here and used the corresponding text Tabula
generated with `tabularecta` or use the one provided in the documentation:
[German](./data/german_tabula.txt), [English](./data/english_tabula.txt), [Spanish](./data/latin_tabula.txt),
[Greek](./data/greek_tabula.txt) and [Cyrillic](./data/cyrillic_tabula.txt).