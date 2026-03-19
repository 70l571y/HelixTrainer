package main

import "fmt"

func processData(data []int) {
	if len(data) > 0 {
		if data[0] > 10 {
			fmt.Println("Large value")
			if data[0] > 100 {
				fmt.Println("Very large value")
			}
		}
	}
}

func main() {
	processData([]int{50})
}
