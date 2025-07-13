package main

import (
	"fmt"
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
func main(){

	var my_list = []string{"hallo how are you", "what is your name", "my name is kaleab hailemeskel", "how old are you ?"}
	for _, val := range my_list{
		fmt.Println(splitString(val))
	}
}