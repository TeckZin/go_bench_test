package main

import (
	"fmt"
	"main/bench"
)

func main() {
	fmt.Println("start")
	benchmark := bench.New()

	tests := map[string]func(){
		"fibonacci(45)":   func() { bench.Fibonacci(45) },
		"primes(1000000)": func() { bench.SieveOfEratosthenes(1000000) },
	}

	benchmark.Compare(tests, 5).PrintResults()
}
