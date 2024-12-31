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

// var dateNums = [3]string{}

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

func normalizeDate(input string) (string, error) {
	// Sanitize input
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "/", " ")
	input = strings.ReplaceAll(input, "-", " ")

	// Split into components
	dateValues := strings.Fields(input)
	if len(dateValues) != 3 {
		return "", fmt.Errorf("invalid date format. Expected MM/DD/YYYY, MM DD YYYY or MM-DD-YYYY")
	}

	// Process month
	letters := regexp.MustCompile(`^[a-zA-Z]+$`)
	digits := regexp.MustCompile(`^[0-9]+$`)

	switch {
	case letters.MatchString(dateValues[0]):
		if len(dateValues[0]) > 3 {
			monthStr := dateValues[0][:3]
			monthStr = strings.ToLower(monthStr)
			dateValues[0] = string(unicode.ToUpper(rune(monthStr[0]))) + monthStr[1:]
		}
	case digits.MatchString(dateValues[0]):
		month, _ := strconv.Atoi(dateValues[0])
		chMonth := monthNumToName(month)
		dateValues[0] = chMonth
	default:
		return "", fmt.Errorf("invalid month format specification")
	}

	// Process year
	year, err := strconv.Atoi(dateValues[2])
	if err != nil || year <= 999 {
		return "", fmt.Errorf("year must be a valid 4-digit value")
	}

	// Return formatted date
	return fmt.Sprintf("%s %s %d", dateValues[0], dateValues[1], year), nil
}

func main() {
	// Interactive mode
	fmt.Fprintf(os.Stderr, "Enter date: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	result, err := normalizeDate(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(result)
}
