//go:build ignore

package main

type upperFormatter struct{}

func (upperFormatter) Format(input string) string {
	return "[" + input + "]"
}
