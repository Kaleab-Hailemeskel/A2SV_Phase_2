package main

import (
	"fmt"
	"strings"
)

func splitString(s string) map[string]int {
	last_word := ""
	res := make(map[string]int)
	
	for right_index:= 0; right_index < len(s); right_index ++{
		if s[right_index] == ' '{
			res[last_word] += 1
			last_word = ""
		}else{
			last_word += string(s[right_index])
		}
	}
	
	if last_word != ""{
		res[last_word] += 1
	}
	return res
}
func checkPalindrome(s string) bool{
	s = strings.ToLower(s)
	new_s := ""
	for _, char := range(s){
		if char != ' '{
			new_s += string(char)
		}
	}
	for i := 0; i < len(new_s); i++{
		rev_i := len(new_s) - i - 1
		if new_s[rev_i] != new_s[i]{
			return false
		}
	}
	return true
}
func main(){

	var my_list = []string{"hallo how are you", "what is your name", "my name is kaleab hailemeskel", "how old are you ?"}
	for _, val := range my_list{
		fmt.Println(splitString(val))
	}
	
	testCases := [][]string{{
		"Racecar",           // Classic, mixed case
		"Madam",             // Another classic, mixed case
		"A man a plan a canal Panama", // Famous, multi-word, mixed case
		"No lemon no melon",  // Another multi-word, mixed case
		"Was it a car or a cat I saw", // Long, complex, mixed case
		"Deleveled",         // Single word, mixed case
		"Eva can I see bees in a cave", // Long, mixed case
		"Go hang a salami Im a lasagna hog", // Another long one, mixed case
		"Level",             // Simple, Capitalized
		"raceCAR",           // Simple, mixed case
	}, {"Hello World",             // Clearly not a palindrome, includes space
		"RacecarX",                // Almost a palindrome, but with an extra char
		"Madam I'm Adam",          // Classic phrase, but not a palindrome when processed this way
		"Abcde",                   // Simple non-palindrome
		"Google",                  // Another simple non-palindrome
		"A man a plan a canal Panana", // Close to the famous one, but wrong ending
		"Top spot",                // Palindrome if case-insensitive, but false for case-sensitive
		"Was it a cat I saw or a car", // Reversed order of a palindrome, so not a palindrome
		"No melon no lemon",       // Reversed order of the "No lemon..." palindrome
		"LevelUp",  }}

	
	for _, list_of_string := range testCases{
		fmt.Println()
		for _, val := range list_of_string{
			fmt.Println(checkPalindrome(val), "\t",val)
		}
	}
	
	
}