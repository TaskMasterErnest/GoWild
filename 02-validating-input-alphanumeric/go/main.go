package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

/***
Validate input
Input has consist of only alphanumeric values
***/

func isAlphaNumeric(input string) {
	str := regexp.MustCompile(`^[a-zA-Z0-9]+$`)

	if str.MatchString(input) {
		fmt.Println("Input is valid!")
	} else {
		fmt.Println("Enter only letters and numerals!")
	}

}

func main() {
	// prompt user for input
	fmt.Print("Enter input: ")

	// Create a new reader
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}

	// Trim the newline character from the input
	input = input[:len(input)-1]

	if input == "" {
		fmt.Println("Invalid input!")
		os.Exit(1)
	}

	// validate input
	isAlphaNumeric(input)
}
