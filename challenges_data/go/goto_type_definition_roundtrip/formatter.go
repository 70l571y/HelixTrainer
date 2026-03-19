//go:build ignore

package main

type Formatter interface {
	Format(input string) string
}
