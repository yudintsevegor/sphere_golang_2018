package main

import "fmt"

func main() {
	test := []int{1,2,3,4,5}
	test = test[:len(test) - 1]
	fmt.Println(test)
}
