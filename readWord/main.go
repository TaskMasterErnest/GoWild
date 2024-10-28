/*
*
In this example, I will attempt to test two inbuilt functions with different test-cases
These functions are the bufio.NewReader.ReadString and the regex.MustCompile("[^\\s]+")
Both are used to read a string of words in a file, one word at a time.
*
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

func useScanWords(file string) (int, error) {
	// open the file
	f, err := os.Open(file)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)

	var wordNumber int

	for scanner.Scan() {
		wordNumber++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return wordNumber, nil
}

func useRegex(file string) (int, error) {
	f, err := os.Open(file)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	reg := regexp.MustCompile(`\S+`) // Matches any non-whitespace sequences

	var wordNumber int

	for {
		line, err := reader.ReadString('\n')
		words := reg.FindAllString(line, -1)
		wordNumber += len(words)

		if err == io.EOF {
			break
		} else if err != nil {
			return 0, err
		}
	}

	return wordNumber, nil
}

func main() {
	filename := flag.String("f", "", "file to read")
	flag.Parse()

	if *filename == "" {
		fmt.Println("Filename is required")
		os.Exit(1)
	}

	fmt.Println(useScanWords(*filename))
}
