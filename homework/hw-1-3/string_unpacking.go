package main

import (
	"errors"
	"log"
	"strings"
	"unicode"
)

func main() {
	unpackStringArray([]string{"a4bc2d5e", "abcd", "45"})
	unpackStringArray([]string{`qwe\4\5`, `qwe\45`, `qwe\\5`})
}

func unpackStringArray(sa []string) {
	for _, s := range sa {
		us, err := unpackString(s)
		if err != nil {
			log.Println("Incorrect string:", s)
		} else {
			log.Println("Unpacked string:", us)
		}
	}
}

const emptyRune rune = 0x00

var (
	bsMode   bool = false
	prevRune rune = emptyRune
)

func switchBackslashMode(input string, flag bool) string {
	bsMode = flag
	prevRune = emptyRune
	return input
}

func unpackString(s string) (us string, err error) {
	// Check for empty string
	if len(s) == 0 {
		return us, errors.New("incorrect string")
	}

	var rs []rune = []rune(s)
	prevRune = rs[0]

	for i := 1; i < len(rs); i++ {
		r := rs[i]

		// Backslash mode logic
		if bsMode {
			if prevRune == emptyRune {
				prevRune = r
			} else {
				if unicode.IsLetter(r) {
					return us, errors.New("incorrect string")
				}
				count := int(r - '0')
				us += switchBackslashMode(strings.Repeat(string(prevRune), count), false)
			}
		} else {
			if r == 0x5C {
				us += switchBackslashMode(string(prevRune), true)
				continue
			}
			// Simple string mode
			switch {
			case unicode.IsDigit(r):
				if unicode.IsDigit(prevRune) {
					return us, errors.New("incorrect string")
				}
				count := int(r - '0')
				us += strings.Repeat(string(prevRune), count)

			case unicode.IsLetter(r) && i == len(rs)-1:
				if unicode.IsLetter(prevRune) {
					us += string(prevRune)
				}
				us += string(r)

			case unicode.IsLetter(prevRune):
				us += string(prevRune)
			}
			prevRune = r
		}
	}

	if bsMode && prevRune != emptyRune {
		us += string(prevRune)
	}

	return us, nil
}
