# D E V E L O P E R

## OS: Unix/Linux/Darwin


## OS: Windows

**NOTE:** *Keep in mind that the repository is case-sensitive; therefore, 
make sure you use the same nomenclature. Above all, do not change the name
of files!.

To build on windows CLI:

 * `go build -buildmode=exe -o bin/gocaesar.exe -gcflags all=-N ./cmd/cli/main.go`
 * Without `-gcflags all=-N` it will produce a non-executable output!

### Showcase Concepts

For Go newbies there are a few interesting things done in this module,
that beyond the basic, may be worth noting:

* Using the GO `flag` package with custom flags: `cmd.rune_flag.go`
* Modular `main()` which makes it look orderly
* Avoiding pitfalls of indexing strings by byte instead of runes
* using interfaces
* Command Pattern with Piped commands (chained commands)
* Custom application logging with `app/mlog.go` which is my improved
   version of `log/slog`

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