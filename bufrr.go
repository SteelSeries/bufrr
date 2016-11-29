// Package bufrr provides a buffered rune reader,
// with both PeekRune and UnreadRune.
// It takes an io.Reader providing the source,
// buffers it by wrapping with a bufio.Reader,
// and creates a new Reader implementing
// the bufrr.RunePeeker interface
// (an io.RuneScanner interface plus an
// additional PeekRune method).
//
// Additionally, bufrr.Reader also translates
// io.EOF error into the invalid rune value of -1
// (defined as bufrr.EOF)
//
// Internally, bufrr.Reader is a bufio.Reader
// plus a single-rune peek buffer and a
// single-rune unread buffer.
package bufrr

import (
	"bufio"
	"errors"
	"io"
)

const EOF = -1

var ErrInvalidUnreadRune = errors.New("bufrr: invalid use of UnreadRune")

type Reader struct {
	br   *bufio.Reader
	peek runeCache
	last runeCache
}

type RunePeeker interface {
	io.RuneScanner
	PeekRune() (r rune, w int, err error)
}

func NewReader(rd io.Reader) *Reader {
	return &Reader{br: bufio.NewReader(rd)}
}

func NewReaderSize(rd io.Reader, size int) *Reader {
	return &Reader{br: bufio.NewReaderSize(rd, size)}
}

func (b *Reader) ReadRune() (r rune, w int, err error) {
	// if there's one in the bag, return that, and empty the bag
	if b.peek.valid() {
		r, w = b.peek.get()
		b.peek.invalidate()
	} else {
		// otherwise, do a read
		r, w, err = b.read()
		if err != nil {
			return
		}
	}
	b.last.set(r, w)
	return
}

func (b *Reader) PeekRune() (r rune, w int, err error) {
	// if there's one in the bag, return that
	if b.peek.valid() {
		r, w = b.peek.get()
	} else {
		// otherwise, do a read and put in the bag
		r, w, err = b.read()
		if err != nil {
			return
		}
		b.peek.set(r, w)
	}
	return r, w, nil
}

func (b *Reader) UnreadRune() error {
	if b.last.invalid() {
		return ErrInvalidUnreadRune
	}
	if b.peek.valid() {
		// this kicks bad bad news at end of stream
		// (i.e. not a simple EOF)
		// so deliberately ignore its error
		b.br.UnreadRune()
		// err := b.br.UnreadRune()
		// if err != io.EOF {
		// 	return err
		// }
	}
	b.peek = b.last
	b.last.invalidate()
	return nil
}

func (b *Reader) read() (rune, int, error) {
	// read next rune from underlying reader, masking EOF by converting it from an error to rune value of -1.
	r, w, err := b.br.ReadRune()
	if err == io.EOF {
		return EOF, 0, nil
	} else if err != nil {
		return 0, 0, err
	}
	return r, w, nil
}

// helper struct holding rune value and size
type runeCache struct {
	ru rune
	sz int
}

// all of these runeCache helpers will get inlined

func (r *runeCache) valid() bool {
	return r.ru != 0
}

func (r *runeCache) invalid() bool {
	return r.ru == 0
}

func (r *runeCache) invalidate() {
	r.ru = 0
}

func (r *runeCache) get() (rune, int) {
	return r.ru, r.sz
}

func (r *runeCache) set(ru rune, sz int) {
	r.ru, r.sz = ru, sz
}
