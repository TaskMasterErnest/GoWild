package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Product struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Desc     string   `json:"description"`
	Price    float64  `json:"price"`
	Category string   `json:"category"`
	Images   []string `json:"images"`
}

func main() {
	var filepath string
	fmt.Print("Enter file path: ")
	fmt.Scan(&filepath)

	content, err := readFile(filepath)

	var p Product

	err = json.Unmarshal(content, &p)
	if err != nil {
		panic(err)
	}

	w, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(w))

	writeToFile("result.json", w)

}

func readFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return data, nil
}

func writeToFile(filename string, data []byte) (int, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		panic(err)
	}

	return 0, nil
}
