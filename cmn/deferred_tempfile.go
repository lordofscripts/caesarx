/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * An object that manages the lifetime of a temporary file, from
 * creation, obtaining its name, closing it and optionally to
 * automatically remove it when it has been closed. This is
 * platform agnostic.
 *-----------------------------------------------------------------*/
package cmn

import (
	"errors"
	"fmt"
	"lordofscripts/caesarx/app/mlog"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const defaultTempFilePattern string = "tempfile-*"

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type DeferredTempFile struct {
	fd        *os.File // temporary file descriptor when created
	directory string   // directory where temp file should be created
	pattern   string   // filename pattern (prefix)
	closed    bool     // (internal) whether it has been already closed
	remove    bool     // remove temporary file on Close()
	lastName  string   // the last known temp filename in case user needs the file
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// Sets up an instance for deferring the creation of a temporary
// file starting with pattern in the name at the dir directory.
func NewDeferredTempFileWith(dir, pattern string) *DeferredTempFile {
	return &DeferredTempFile{
		fd:        nil,
		directory: dir,
		pattern:   pattern,
		closed:    true,
		remove:    false,
		lastName:  "",
	}
}

// Sets up an instance for deferring the creation of a temporary
// file at the dir directory.
func NewDeferredTempFileAt(dir string) *DeferredTempFile {
	return NewDeferredTempFileWith(dir, defaultTempFilePattern)
}

// Sets up an instance for deferring the creation of a temporary
// file at the OS temporary files directory.
func NewDeferredTempFile() *DeferredTempFile {
	return NewDeferredTempFileWith(os.TempDir(), defaultTempFilePattern)
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer and gives information about the temp file
func (d *DeferredTempFile) String() string {
	active := "no"
	name := ""
	if d.fd != nil {
		name = d.fd.Name()
		active = "yes"
	}

	return fmt.Sprintf("Active:%s Closed:%t %s", active, d.closed, name)
}

// request removal of this instance's temporary file. The removal is
// done when the (deferred) Close() method is called.
func (d *DeferredTempFile) WithRemoval() *DeferredTempFile {
	d.remove = true
	return d
}

// if the file has not been created yet, add the ext file extension
// to the pattern. If ext has no leading "." then it is prepended by
// this method.
func (d *DeferredTempFile) WithExtension(ext string) *DeferredTempFile {
	if d.fd == nil && len(ext) > 0 {
		const EXT_SEP string = "."
		ext = strings.Trim(ext, " \t")
		if !strings.HasPrefix(ext, EXT_SEP) {
			ext = EXT_SEP + ext
		}

		d.pattern = d.pattern + ext
	}

	return d
}

// creates the temporary file with the specification given at the constructor
// and returns nil on success. The caller must defer the call to this
// object's Close() method to dispose of the temporary file.
func (d *DeferredTempFile) Create() error {
	if d.fd != nil {
		return errors.New("this DeferredTempFile is already activated")
	}

	var err error = nil
	d.lastName = ""
	d.fd, err = os.CreateTemp(d.directory, d.pattern)
	if err != nil {
		mlog.Error("unable to create temporary file", err, mlog.String("Dir", d.directory))
	} else {
		d.lastName = d.fd.Name()
	}

	return err
}

// get the temporary file's filename or empty if none.
func (d *DeferredTempFile) GetFilename() string {
	return d.lastName
}

// get the temporary file's descriptor for operating on the file.
// The caller should limit itself to writing to the file and NOT close it
// because that is handled by this object's instance.
func (d *DeferredTempFile) GetFile() *os.File {
	return d.fd
}

// closes the temporary file and if requested, delete the file.
func (d *DeferredTempFile) Close() {
	if d.fd != nil && !d.closed {
		filename := d.fd.Name()
		d.fd.Close()
		d.fd = nil
		d.closed = true

		// the user requested removal?
		if d.remove {
			os.Remove(filename)
		}
	}
}

// generates a temporary filename (not a file!) where the base filename
// is a pattern where the "*" is replaced by a random number. Any directory
// part will be removed because this function uses the system's Temp directory.
func GenerateTemporaryFileName(pattern string) string {
	getRandom := func() int {
		// Generate a non-deterministic random integer between 0 and 100
		randomInt := rand.Intn(100_000) // Specify the upper limit (exclusive)
		return randomInt
	}

	const STAR string = "*"
	if len(pattern) == 0 {
		pattern = "tempfile-*"
	} else {
		// we will use the system's TEMP directory
		pattern = path.Base(pattern)
		// add a pattern placeholder if none given
		if !strings.Contains(pattern, STAR) {
			pattern = pattern + STAR
		}
	}

	for {
		star := strings.Index(pattern, "*")
		if star == -1 {
			break
		}

		pattern = pattern[:star] + strconv.Itoa(getRandom()) + pattern[star+1:]
	}

	return filepath.Join(os.TempDir(), pattern)
}

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/

/*
// Demonstrates using the DeferredTempFile object.
func DeferredTempFileDemo() {
	// we want a temporary file to dispose at the end.
	// or eliminate WithRemoval() if you have further use for the file.
	dtf := NewDeferredTempFile().WithRemoval()
	// let the object close and remove when exiting
	defer dtf.Close()

	// when you need it
	if dtf.Create() == nil {
		fmt.Println("Temp filename is:", dtf.GetFilename())

		fd := dtf.GetFile()
		fd.WriteString("dummy data")
		// do something else before discarding it
	}
}
*/
