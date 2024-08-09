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
	Quantity int      `json:"quantity"`
	Price    float64  `json:"price"`
	Category string   `json:"category"`
	Images   []string `json:"images"`
}

func main() {
	var filepath string
	fmt.Print("Enter file path: ")
	fmt.Scan(&filepath)

	content, err := readProductsFromFile(filepath)
	if err != nil {
		fmt.Println("Error reading content:", err)
		return
	}
	fmt.Printf("%v : %T\n", content, content)

	writeProductsToFile("results.json", content)
}

func readProductsFromFile(filename string) ([]Product, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	// validate json
	if !json.Valid(data) {
		return nil, fmt.Errorf("invalid JSON data")
	}

	var products []Product

	err = json.Unmarshal(data, &products)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling products: %w", err)
	}

	return products, nil
}

func writeProductsToFile(filename string, data []Product) (int, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	products := data
	jsonEncode, err := json.MarshalIndent(products, "", "    ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if _, err := file.Write(jsonEncode); err != nil {
		panic(err)
	}

	return 0, nil
}
