package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var daysInMonth = map[string]int{
	"Jan": 31, "Feb": 28, "Mar": 31,
	"Apr": 30, "May": 31, "Jun": 30,
	"Jul": 31, "Aug": 31, "Sep": 30,
	"Oct": 31, "Nov": 30, "Dec": 31,
}

// check if a year is a Leap year
func isLeapYear(n int) bool {
	if (n%4 == 0 && n%100 != 0) || (n%400 == 0) {
		return true
	} else {
		return false
	}
}

// ensure the days in the month presented are valid
func exceedsDaysInMonth(month string, day int) bool {
	// lookup the number of days in the calendar month
	// validate it against the day input given
	days, exists := daysInMonth[month]
	if !exists {
		fmt.Fprintf(os.Stderr, "unknown month name: %s\n", month)
		return false
	}

	if day < 1 || day > days {
		return false
	} else {
		return true
	}
}

func main() {
	// run the date normalization command
	pathToNormDate := "../../03-normalizing-date-formats/go/normdate"
	cmd := exec.Command(pathToNormDate)
	// Connect stdin for interactive input
	cmd.Stdin = os.Stdin
	// Connect stderr to show prompts
	cmd.Stderr = os.Stderr
	// Capture only stdout for the result
	output, err := cmd.Output()
	if err != nil {
		os.Exit(1)
	}
	// Use the captured output as date
	date := strings.TrimSpace(string(output))

	// split up output string into components
	dateValues := strings.Split(date, " ")
	// validate date length
	if len(dateValues) != 3 {
		fmt.Fprintf(os.Stderr, "invalid date format: %s\n", date)
		os.Exit(1)
	}

	month := dateValues[0]
	day, err := strconv.Atoi(dateValues[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid day value: %s\n", dateValues[1])
		os.Exit(1)
	}
	year, err := strconv.Atoi(dateValues[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid year value: %s\n", dateValues[2])
		os.Exit(1)
	}

	// start validation
	if !exceedsDaysInMonth(month, day) {
		if month == "Feb" && day == 29 {
			if !isLeapYear(year) {
				fmt.Fprintf(os.Stderr, "%d is not a leap year, so Feb does not have 29 days.\n", year)
				os.Exit(1)
			}
		} else {
			fmt.Fprintf(os.Stderr, "invalid day value: %s does not have %d days.\n", month, day)
			os.Exit(1)
		}
	}

	fmt.Fprintf(os.Stdin, "Valid date: %s %d %d\n", month, day, year)
}
