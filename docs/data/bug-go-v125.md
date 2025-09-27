<!-- Please answer these questions before submitting your issue. Thanks! -->

### What version of Go are you using (`go version`)?

<pre>
$ go version
go version go1.25.0 linux/amd64
</pre>

### Does this issue reproduce with the latest release?

Yes

### What operating system and processor architecture are you using (`go env`)?

<details><summary><code>go env</code> Output</summary><br><pre>
$ go env
AR='ar'
CC='gcc'
CGO_CFLAGS='-O2 -g'
CGO_CPPFLAGS=''
CGO_CXXFLAGS='-O2 -g'
CGO_ENABLED='1'
CGO_FFLAGS='-O2 -g'
CGO_LDFLAGS='-O2 -g'
CXX='g++'
GCCGO='gccgo'
GO111MODULE=''
GOAMD64='v1'
GOARCH='amd64'
GOAUTH='netrc'
GOBIN='/home/lordofscripts/go/bin'
GOCACHE='/home/lordofscripts/.cache/go-build'
GOCACHEPROG=''
GODEBUG=''
GOENV='/home/lordofscripts/.config/go/env'
GOEXE=''
GOEXPERIMENT=''
GOFIPS140='off'
GOFLAGS=''
GOGCCFLAGS='-fPIC -m64 -pthread -Wl,--no-gc-sections -fmessage-length=0 -ffile-prefix-map=/tmp/go-build3576437392=/tmp/go-build -gno-record-gcc-switches'
GOHOSTARCH='amd64'
GOHOSTOS='linux'
GOINSECURE=''
GOMOD='/dev/null'
GOMODCACHE='/home/lordofscripts/go/pkg/mod'
GONOPROXY=''
GONOSUMDB=''
GOOS='linux'
GOPATH='/home/lordofscripts/go'
GOPRIVATE=''
GOPROXY='https://proxy.golang.org,direct'
GOROOT='/usr/local/go'
GOSUMDB='sum.golang.org'
GOTELEMETRY='local'
GOTELEMETRYDIR='/home/lordofscripts/.config/go/telemetry'
GOTMPDIR=''
GOTOOLCHAIN='auto'
GOTOOLDIR='/usr/local/go/pkg/tool/linux_amd64'
GOVCS=''
GOVERSION='go1.25.0'
GOWORK=''
PKG_CONFIG='pkg-config'
uname -sr: Linux 6.12.41+deb13-amd64
Distributor ID:	Debian
Description:	Debian GNU/Linux 13 (trixie)
Release:	13
Codename:	trixie
/lib/x86_64-linux-gnu/libc.so.6: GNU C Library (Debian GLIBC 2.41-12) stable release version 2.41.
gdb --version: GNU gdb (Debian 16.3-1) 16.3
</pre></details>

### What did you do?

First of all I want to say something IMPORTANT. I think it is GREAT that the GO language is consistent in this: that the German Es-tzet character `ß`  (a single-character **rune**) has an equivalent uppercase `ẞ` **single-character rune** rather than `SS`. I say it because THAT is consistent, if your algorithm is handling **runes** it expects **runes** as return. Why? because Lowercase `ß` (rune) has uppercase `ẞ` (also a rune)  is programmatically correct, whereas saying lowercase `ß` is uppercase `SS` (double-rune i.e. **string**) is not programmatically correct. We don't want to break algorithms with exceptions.

Now, it so happens I spent hours debugging (long live unit tests!) an application where German was one of the input languages. I found out that the culprit was
that the standard GO library is not quite consistent in its doings when it comes
to this special character:

<!--
If possible, provide a recipe for reproducing the error.
A complete runnable program is good.
A link on go.dev/play is best.
-->

I prepared a Test Case that shows `strings.ToUpper` does not follow the correct GO assumption (speaking as a programmer, not a linguist!):

```
func Test_GermanCharacter(t *testing.T) {
	const LOWER_RUNE rune = 'ß' // 223
	const UPPER_RUNE rune = 'ẞ' // 7338
	const LOWER_STRING string = "daß"
	const UPPER_STRING string = "DAẞ"

	// just to make sure, these two PASS
	if LOWER_RUNE == UPPER_RUNE { // passes
		t.Errorf("This isn't supposed to happen in 1.25")
	}
	if LOWER_STRING == UPPER_STRING { // passes
		t.Errorf("This isn't supposed to happen in 1.25")
	}

	fmt.Printf("Lowercase es-tzet: %c (%U)\n", LOWER_RUNE, LOWER_RUNE)
	fmt.Printf("Uppercase es-tzet: %c (%U)\n", UPPER_RUNE, UPPER_RUNE)

	// Now the STD Lib inconsistencies
	var gotS, expectS string

	// (a) strings.ToUpper("daß") should be "DAẞ"
	expectS = UPPER_STRING
	gotS = strings.ToUpper(LOWER_STRING)
	if gotS != expectS { // FAILS with Std. Lib. v1.25
		t.Errorf("strings.ToUpper failure Got:'%s' Expect:'%s'", gotS, expectS)
	}

	// (b) strings.ToLower("DAẞ") should be "daß"
	expectS = LOWER_STRING
	gotS = strings.ToLower(UPPER_STRING)
	if gotS != expectS {
		t.Errorf("strings.ToLower failure Got:'%s' Expect:'%s'", gotS, expectS)
	}

	var gotR, expectR rune

	// (c) unicode.ToUpper('ß') should be 'ẞ'
	expectR = UPPER_RUNE
	gotR = unicode.ToUpper(LOWER_RUNE)
	if gotS != expectS {
		t.Errorf("unicode.ToUpper failure Got:'%c' Expect:'%c'", gotR, expectR)
	}

	// (d) unicode.ToLower('ẞ') should be 'ß'
	expectR = LOWER_RUNE
	gotR = unicode.ToLower(UPPER_RUNE)
	if gotS != expectS {
		t.Errorf("unicode.ToLower failure Got:'%c' Expect:'%c'", gotR, expectR)
	}
}
```

Even better, here is a different non-test version:  a [Go Play example](https://go.dev/play/p/b7UxbQManzI) demonstrating the findings.

### What did you expect to see?

Conversions of  `unicode.ToUpper()`, `unicode.ToLower() ` and `strings.ToUpper()`, `strings.ToLower()`  should be consistent that a lower to upper must be uppercased, and that an upper to lower should result in lowercase.

### What did you see instead?

Given the CORRECT programming-wise assumption that a conversion of a rune to/from Upper/Lowercase must **also** be a RUNE type, the library functions must be consistent in its translations. For the German ` ß` character I observed:

* `unicode.ToUpper()` gives an incorrect value (does not convert)
* `strings.ToUpper()` gives an improper conversion if it contains `ß` in the string

