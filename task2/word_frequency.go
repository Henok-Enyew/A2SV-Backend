package main

import (
	"fmt"
	"regexp"
	"strings"
)

func WordFrequency(text string) map[string]int {
	re := regexp.MustCompile(`[^\w\s]`)
	cleaned := re.ReplaceAllString(text, " ")
	words := strings.Fields(strings.ToLower(cleaned))
	
	frequency := make(map[string]int)
	for _, word := range words {
		frequency[word]++
	}
	
	return frequency
}

func main() {
	text1 := "Hello world hello"
	fmt.Println(WordFrequency(text1))
	
	text2 := "The quick brown fox jumps over the lazy dog. The dog was lazy!"
	fmt.Println(WordFrequency(text2))
	
	text3 := "Go is great! Go, go, GO!"
	fmt.Println(WordFrequency(text3))
}

