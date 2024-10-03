// Package getopt provides the good old getopt function as a Go iterator.
package getopt

import (
	"errors"
	"fmt"
	"iter"
	"strings"
	"unicode/utf8"
)

type Option struct {
	Name rune
	Arg  *string
}

func parseOptstring(optstring string) map[rune]int {
	opts := map[rune]int{}
	var last rune
	for _, r := range optstring {
		if r == ':' {
			if last != 0 {
				opts[last]++
			}
			continue
		}
		if _, ok := opts[r]; !ok {
			opts[r] = 0
		}
		last = r
	}
	return opts
}

var RestArgs []string

func Getopt(args []string, optstring string) iter.Seq2[Option, error] {
	opts := parseOptstring(optstring)
	return func(yield func(Option, error) bool) {
		for i := 0; i < len(args); i++ {
			curr := args[i]
			if !strings.HasPrefix(curr, "-") {
				// No options.
				RestArgs = args[i:]
				return
			}
			switch curr {
			case "-":
				RestArgs = args[i:]
				yield(Option{Name: 0, Arg: nil}, errors.New("a single \"-\" is not supported"))
				return
			case "--":
				RestArgs = args[i+1:]
				return
			default:
				curr = curr[1:]
			}

			for j, opt := range curr {
				optCnt, ok := opts[opt]
				// An option ":" is not allowed.
				// Options not included in opts are not allowed.
				if opt == ':' || !ok {
					RestArgs = args[i+1:]
					if !yield(Option{Name: opt, Arg: nil}, fmt.Errorf("illegal option: '%c'", opt)) {
						return
					}
					continue
				}
				// No arguments required.
				if optCnt == 0 {
					RestArgs = args[i+1:]
					if !yield(Option{Name: opt, Arg: nil}, nil) {
						return
					}
					continue
				}
				// Treat the rest of "curr" as the argument, if available.
				if tail := j + utf8.RuneLen(opt); tail < len(curr) {
					RestArgs = args[i+1:]
					arg := curr[tail:]
					if !yield(Option{Name: opt, Arg: &arg}, nil) {
						return
					}
					break
				}
				// Use next args as the argument, if available.
				if i+1 < len(args) {
					RestArgs = args[i+2:]
					if !yield(Option{Name: opt, Arg: &args[i+1]}, nil) {
						return
					}
					i++
					break
				}
				// No arguments found.
				RestArgs = nil
				yield(Option{Name: opt, Arg: nil}, fmt.Errorf("no arguments supplied: '%c'", opt))
				return
			}
		}
	}
}
