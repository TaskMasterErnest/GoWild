package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type wordCount struct {
	word  string
	count int
}

func readReview(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var content []string

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		content = append(content, line)
	}

	return content
}

func cleanReview(s []string) []string {
	var cleaned []string

	for _, v := range s {
		v = strings.ToLower(v)

		// capture all punctuations
		re := regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
		// replace puctuations with empty string ""
		v = re.ReplaceAllString(v, "")

		cleaned = append(cleaned, v)
	}

	return cleaned
}

func analyzeWordFrequency(s []string) map[string]int {
	wordFrequency := make(map[string]int)

	for _, line := range s {
		words := strings.Fields(line)
		for _, word := range words {
			wordFrequency[word]++
		}
	}

	return wordFrequency
}

func getTopNWords(frequency map[string]int, n int64) []wordCount {
	var wordCountSlice []wordCount
	for word, count := range frequency {
		wordCountSlice = append(wordCountSlice, wordCount{word, count})
	}
	// sort the slice in descending order
	sort.Slice(wordCountSlice, func(i, j int) bool {
		return wordCountSlice[i].count > wordCountSlice[j].count
	})

	return wordCountSlice[:n]
}

func main() {
	var pathToFile string
	fmt.Print("Path to review file: ")
	fmt.Scan(&pathToFile)
	fmt.Println("----------xxx----------xxx----------xxx----------xxx----------xxx")
	data := analyzeWordFrequency(cleanReview(readReview(pathToFile)))

	fmt.Println()
	var numberStr string
	fmt.Print("Insert number of top counts: ")
	fmt.Scan(&numberStr)
	fmt.Println("<<<---------->>>----------<<<---------->>>----------<<<---------->>>")
	number, _ := strconv.ParseInt(numberStr, 10, 64)
	result := getTopNWords(data, number)
	fmt.Print(result)

}
