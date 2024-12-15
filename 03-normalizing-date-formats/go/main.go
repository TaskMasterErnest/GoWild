package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var dateNums = [3]string{}

// match the month number to the abbreviated name
func monthNumToName(num int) string {
	switch num {
	case 1:
		return "Jan"
	case 2:
		return "Feb"
	case 3:
		return "Mar"
	case 4:
		return "Apr"
	case 5:
		return "May"
	case 6:
		return "Jun"
	case 7:
		return "Jul"
	case 8:
		return "Aug"
	case 9:
		return "Sep"
	case 10:
		return "Oct"
	case 11:
		return "Nov"
	case 12:
		return "Dec"
	default:
		return fmt.Sprintf("Unknown month value: %d", num)
	}
}

func main() {
	// read arguments from Stdin
	fmt.Print("Enter date: ")
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("error reading input: ", err)
		os.Exit(1)
	}

	// check if format has / or - in dates and sanitize
	// replace these with empty spaces
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "/", " ")
	input = strings.ReplaceAll(input, "-", " ")

	// split into components and perform length validation
	// use Fields to handle multiple spaces
	dateValues := strings.Fields(input)
	if len(dateValues) != 3 {
		fmt.Println("Invalid date format. Expected MM/DD/YYYY, MM DD YYYY or MM-DD-YYYY.")
		os.Exit(1)
	}

	// check if month is either a string or an int
	// using regex to check for these
	letters := regexp.MustCompile(`^[a-zA-Z]+$`)
	digits := regexp.MustCompile(`^[0-9]+$`)

	switch {
	case letters.MatchString(dateValues[0]):
		if len(dateValues[0]) > 3 {
			// Take first three letters
			monthStr := dateValues[0][:3]
			// Convert to title case
			monthStr = strings.ToLower(monthStr)
			dateValues[0] = string(unicode.ToUpper(rune(monthStr[0]))) + monthStr[1:]
		}
	case digits.MatchString(dateValues[0]):
		month, _ := strconv.Atoi(dateValues[0])
		chMonth := monthNumToName(month)
		dateValues[0] = chMonth
	default:
		fmt.Println("Invalid month format specification.")
		os.Exit(1)
	}

	// check if day is a valid day
	day, err := strconv.Atoi(dateValues[1])
	if err != nil || day < 1 || day > 31 {
		fmt.Println("Day must be between 1 and 31")
		os.Exit(1)
	}

	// check if year given is a valid year
	year, err := strconv.Atoi(dateValues[2])
	if err != nil || year <= 999 {
		fmt.Println("Year must be a valid 4-digit value.")
		os.Exit(1)
	}

	fmt.Printf("%s %s %d\n", dateValues[0], dateValues[1], year)

}
