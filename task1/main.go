package main

import "fmt"

func SumOfNumbers(numbers []int) int {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return sum
}

func main() {
	numbers1 := []int{1, 2, 3, 4, 5}
	fmt.Printf("Sum of %v = %d\n", numbers1, SumOfNumbers(numbers1))

	numbers2 := []int{-1, 0, 1}
	fmt.Printf("Sum of %v = %d\n", numbers2, SumOfNumbers(numbers2))

	numbers3 := []int{}
	fmt.Printf("Sum of %v = %d\n", numbers3, SumOfNumbers(numbers3))

	numbers4 := []int{10}
	fmt.Printf("Sum of %v = %d\n", numbers4, SumOfNumbers(numbers4))
}

