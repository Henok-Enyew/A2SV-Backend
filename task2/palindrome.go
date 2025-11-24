package main

import (
	"fmt"
	"regexp"
	"strings"
)

func IsPalindrome(text string) bool {
	re := regexp.MustCompile(`[^\w]`)
	cleaned := re.ReplaceAllString(strings.ToLower(text), "")
	
	left := 0
	right := len(cleaned) - 1
	
	for left < right {
		if cleaned[left] != cleaned[right] {
			return false
		}
		left++
		right--
	}
	
	return true
}

func main() {
	fmt.Println(IsPalindrome("racecar"))
	fmt.Println(IsPalindrome("A man a plan a canal Panama"))
	fmt.Println(IsPalindrome("hello"))
	fmt.Println(IsPalindrome("Madam"))
	fmt.Println(IsPalindrome("12321"))
	fmt.Println(IsPalindrome("No 'x' in Nixon"))
}

