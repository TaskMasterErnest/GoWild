package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ErrFileNotFound = errors.New("file not found!")
	input           string
	dataStore       []int
	sumOfCorrect    int
	sumOfWrong      int
)

func readCSV(file io.Reader) {
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
		} else {
			sumOfWrong++
		}
	}

	fmt.Printf("You got %d right out of %d.\n", sumOfCorrect, sumOfWrong+sumOfCorrect)

}

func main() {
	csvFile := flag.String("f", "", "CSV quiz file")
	timer := flag.Duration("t", 30*time.Second, "Timer for quiz")
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

	// start timer function
	var wg sync.WaitGroup

	// handle the timer routine
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(*timer)
		fmt.Println("Timer Expired. Stopping Quiz!!!")
		os.Exit(1)
	}()

	// handle the main routine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			readCSV(file)
			os.Exit(1)
		}
	}()

	wg.Wait()
}
