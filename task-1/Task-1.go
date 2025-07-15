package main

import (
	"fmt"
)

func calcAverage(sub_grade map[string]float32) float32 {
	res := float32(0)
	for _, val := range sub_grade {
		res += val
	}
	return res / float32(len(sub_grade))
}
func main() {
	var num_grade int
	var subject_grade = map[string]float32{}

	fmt.Println("Welcome to Grade Calculator")
	fmt.Print("How many courses do you want to enter: ")
	fmt.Scanf("%d", &num_grade)
	fmt.Printf("Enter all %d subject name with your respective grade\n", num_grade)
	fmt.Println("Enter the values using space between \n\tthem and after you finish each subject-grade entery press ENTER")

	var last_subject string
	var last_grade float32

	for i := 0; i < num_grade; {
		fmt.Print(i+1, ") ")
		_, err := fmt.Scan(&last_subject, &last_grade)
		if last_grade >= 0 && last_grade < 101 && last_subject != "" && err == nil {
			i++ // only accept the next sub-grade pair if the current is valid, otherwise prompt again and again
			fmt.Println("Saved")
			subject_grade[last_subject] = last_grade
		} else {
			if last_subject == "" {
				fmt.Println("Empty Subject Stirng is Not Valid")
			} else if last_grade < 0 || last_grade > 100 {
				fmt.Printf("Invalid Grade report %f\n", last_grade)
			} else {
				fmt.Println("Invalid grade value")
			}
		}
	}

	fmt.Println("Your Average result is: ", calcAverage(subject_grade))
}
