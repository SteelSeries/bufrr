package bufrr

import (
	// "bufio"
	"fmt"
	// "io"
	"strings"
	// "testing"
)

// TODO
//
// peek
// read
//  read
//   peek
//  unread
//  peek
// loop until done

func ExampleReader() {

	s := strings.NewReader("abc")
	b := NewReader(s)

	r, _, _ := b.ReadRune()
	fmt.Printf("%d\n", r)
	r, _, _ = b.ReadRune()
	fmt.Printf("%d\n", r)
	r, _, _ = b.ReadRune()
	fmt.Printf("%d\n", r)
	r, _, _ = b.ReadRune()
	fmt.Printf("%d\n", r)
	r, _, _ = b.ReadRune()
	fmt.Printf("%d\n", r)

	// Output:
	// 97
	// 98
	// 99
	// -1
	// -1
}

func ExampleReader2() {

	s := strings.NewReader("abc")
	b := NewReader(s)

	r, _, _ := b.ReadRune()
	fmt.Printf("%d\n", r)
	r, _, _ = b.ReadRune()
	fmt.Printf("%d\n", r)
	r, _, _ = b.ReadRune()
	fmt.Printf("%d\n", r)
	b.UnreadRune()
	r, _, _ = b.ReadRune()
	fmt.Printf("%d\n", r)

	// Output:
	// 97
	// 98
	// 99
	// 99
}

func Example_Peek() {
	testPeek("s")
	testPeek("str")
	// Output:
	// String: s
	// Read: U+0073 's'
	// Peek: U+FFFFFFFFFFFFFFFF
	// Unread
	// Peek: U+0073 's'
	// Read: U+0073 's'
	//
	// String: str
	// Read: U+0073 's'
	// Peek: U+0074 't'
	// Unread
	// Peek: U+0073 's'
	// Read: U+0073 's'
}

func testPeek(str string) {

	fmt.Printf("String: %s\n", str)

	// init
	s := strings.NewReader(str)
	b := NewReader(s)

	// read a rune
	r, _, err := b.ReadRune()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Read: %#U\n", r) // s

	// look ahead
	r, _, err = b.PeekRune()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Peek: %#U\n", r) // t

	// unread a rune
	b.UnreadRune()
	fmt.Println("Unread")

	// look ahead
	r, _, err = b.PeekRune()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Peek: %#U\n", r) // t

	// read rune again
	r, _, err = b.ReadRune()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Read: %#U\n\n", r) // t
}
