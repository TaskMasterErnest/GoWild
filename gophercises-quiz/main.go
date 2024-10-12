package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ErrFileNotFound = errors.New("file not found!")
	input           string
	sumOfCorrect    int
	totalRecords    int
)

func readCSV(file io.Reader, done chan bool) {
	reader := csv.NewReader(file)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		question := record[0]
		answer, _ := strconv.Atoi(record[1])

		// asking user for input
		fmt.Print("What is ", question, "?\n")
		fmt.Scan(&input)

		// convert user input to int
		userInput, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			fmt.Printf("value %d is not an integer.", err)
		}

		if userInput == answer {
			sumOfCorrect++
		}
	}

	done <- true
}

func main() {
	csvFile := flag.String("f", "", "CSV quiz file")
	timing := flag.Int("t", 30, "Timer for quiz")
	flag.Parse()

	if *csvFile == "" {
		fmt.Println("filename is required!")
		flag.PrintDefaults()
		os.Exit(1)
	}

	file, err := os.Open(*csvFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	// get total file size
	scanFile(file)
	file.Seek(0, io.SeekStart)

	done := make(chan bool)
	go readCSV(file, done)

	// timing
	timer := time.NewTimer(time.Duration(*timing) * time.Second)
	select {
	case <-timer.C:
		fmt.Printf("\nTimer expired! You got %d right out of %d.\n", sumOfCorrect, totalRecords)
	case <-done:
		fmt.Printf("Quiz completed! You got %d right out of %d.\n", sumOfCorrect, totalRecords)
	}
}

func scanFile(file io.Reader) error {
	// get the total number of records
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		totalRecords++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return nil
}
