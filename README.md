# koron-go/getopt

[![PkgGoDev](https://pkg.go.dev/badge/github.com/koron-go/getopt)](https://pkg.go.dev/github.com/koron-go/getopt)
[![Actions/Go](https://github.com/koron-go/getopt/workflows/Go/badge.svg)](https://github.com/koron-go/getopt/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron-go/getopt)](https://goreportcard.com/report/github.com/koron-go/getopt)

Package getopt provides the good old getopt function as a Go iterator.

## Getting started

Install and upgrade:

```
go get github.com/koron-go/getopt@latest
```

Short example:

```go
package main

import (
	"log"
	"os"

	"github.com/koron-go/getopt"
)

func usage() {
	// TODO: Show usage and exit.
	os.Exit(1)
}

func main() {
	var bflag bool
	var file *os.File

	for opt, err := range getopt.Getopt(os.Args[1:], "bf:") {
		if err != nil {
			log.Fatalf("getopt failed: opt=%+v: %s", opt, err)
		}
		switch opt.Name {
		case 'b':
			bflag = true
		case 'f':
			f, err := os.Open(*opt.Arg)
			if err != nil {
				log.Fatalf("failed to open a file: %s", err)
			}
			defer f.Close()
			file = f
		default:
			usage()
		}
	}

	// TODO: Do your task with bflag, file, and getopt.RestArgs
	log.Printf("Do your task with bflag=%t f=%+v getopt.RestArgs=%+v", bflag, file, getopt.RestArgs)
}
```
