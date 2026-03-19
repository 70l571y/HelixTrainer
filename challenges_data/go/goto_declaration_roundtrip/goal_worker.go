//go:build ignore

package main

type workerRunner struct{}

func (workerRunner) Execute(input string) string {
	return "runner:" + input
}
