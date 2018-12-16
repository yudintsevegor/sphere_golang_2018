package main

import (
	"fmt"
)

func main() {
	one := "one"
	two := 2
	switch{
	case 1 < 2:
		fmt.Println(one)
	case 2 < 1:
		fmt.Println(two)
	}

}
