package main

import (
	"fmt"
	"os"
	"os/exec"
)

var daysInMonth = map[string]int{
	"Jan": 31, "Feb": 28, "Mar": 31,
	"Apr": 30, "May": 31, "Jun": 30,
	"Jul": 31, "Aug": 31, "Sep": 30,
	"Oct": 31, "Nov": 30, "Dec": 31,
}

// check if a year is a Leap year
func isLeapYear(n int) bool {
	if (n%4 == 0) && (n%400 == 0) && (n%100 != 0) {
		return true
	} else {
		return false
	}
}

// ensure the days in the month presented are valid
func exceedsDaysInMonth(month string, day int) (bool, int) {
	// lookup the number of days in the calendar month
	// validate it against the day input given
	days, exists := daysInMonth[month]
	if !exists {
		fmt.Fprintf(os.Stderr, "unknown month name: %s\n", month)
		return false, 1
	}

	if day < 1 || day > days {
		return false, 1
	} else {
		return true, 0
	}
}

func main() {
	// check the input arguments
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s month day year\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Typical input format is August 3 1903 or 8 3 1903.")
		os.Exit(1)
	}

	// run the date normalization command
	pathToNormDate := "../../03-normalizing-date-formats/go/normdate"

	cmd := exec.Command(pathToNormDate)

	// stdin, err := cmd.StdinPipe()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// stdin.Write([]byte(os.Args[1]))

	out, err := cmd.CombinedOutput()
	if err != nil {
		os.Exit(1)
	}

	fmt.Println(out)
}
