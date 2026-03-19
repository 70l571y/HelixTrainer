//go:build ignore

package main

type workerRunner struct{}

func (workerRunner) Run(input string) string {
	return "runner:" + input
}
