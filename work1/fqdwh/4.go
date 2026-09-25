package main

import (
	"fmt"
	"math"
)

func isPrime(x int) bool {
	if x == 2 {
		return true
	}
	a := int(math.Sqrt(float64(x)))
	for i := 2; i <= a; i++ {
		if x%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	var n int
	fmt.Scan(&n)
	if isPrime(n) {
		fmt.Println("Yes")
	} else {
		fmt.Println("No")
	}
}
