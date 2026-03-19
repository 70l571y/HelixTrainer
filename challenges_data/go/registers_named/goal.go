package main

import "fmt"

func funcC() {
	fmt.Println("Function C")
}

func funcA() {
	fmt.Println("Function A")
}

func funcB() {
	fmt.Println("Function B")
}

func main() {
	funcC()
	funcA()
	funcB()
}
