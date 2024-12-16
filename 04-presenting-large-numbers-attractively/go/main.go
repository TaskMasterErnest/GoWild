package main

import (
	"flag"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func niceNumber(num float64, decDelim, thousandDelim string) string {
	// separate the number into decimal and integer parts
	intPart, fracPart := math.Modf(num)

	// convert integer part to string
	intStr := strconv.FormatInt(int64(intPart), 10)
	// // convert fraction part to string
	// fracStr := strconv.FormatInt(int64(fracPart), 10)

	// add commas
	var result strings.Builder
	strLength := len(intStr)
	for i, digit := range intStr {
		if i > 0 && (strLength-i)%3 == 0 {
			result.WriteString(thousandDelim)
		}
		result.WriteRune(digit)
	}

	// handle the decimal part
	formatted := result.String()
	if fracPart != 0 {
		// convert decimal part to string
		decStr := strconv.FormatFloat(fracPart, 'f', 5, 64)[2:]
		formatted += decDelim + decStr
	}

	return formatted

}

func main() {
	decDelimiter := flag.String("d", ".", "decimal delimiter")
	thousandDelimiter := flag.String("t", ",", "thousands delimiter")
	number := flag.Float64("n", 0.0, "number to format")

	flag.Parse()

	fmt.Println(niceNumber(*number, *decDelimiter, *thousandDelimiter))
}
