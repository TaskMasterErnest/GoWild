# Word Frequency Analyzer in Go

This program analyzes word frequencies in a text file and identifies the N most frequently used words.

## Features

* Reads text reviews from a specified file path.
* Cleans the review text by converting it to lowercase and removing punctuation.
* Calculates the frequency of each word in the cleaned text.
* Identifies the top N most frequent words based on user input.

## Getting Started

### Prerequisites

* Go installed and configured on your system.

### Usage

1. Save the code as `word_frequency_analyzer.go`.
2. Open a terminal and navigate to the directory containing the file.
3. Run the program:
```bash
go run word_frequency_analyzer.go
```
This will prompt you to enter the path to your review file.

4. Enter the path to your review file and press Enter.
5. The program will clean the text, analyze word frequencies, and prompt you to enter the number of top words to display.
6. Enter the desired number (N) and press Enter.

The program will then print the top N most frequent words and their corresponding counts.