package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir() && info.Mode().Perm()&0111 != 0
}

func inPath(command, pathDir string) (string, error) {
	// Check if the command is in the PATH directories
	dirs := strings.Split(pathDir, string(os.PathListSeparator))
	for _, directory := range dirs {
		cmdPath := filepath.Join(directory, command)
		if isExecutable(cmdPath) {
			return fmt.Sprintf("%s found in PATH\n", command), nil
		}
	}
	return "", fmt.Errorf("%s not found in PATH or not executable", command)
}

func main() {
	// Use flags to get command parameter
	command := flag.String("p", "", "command to find and execute")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -p <command>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	// Get system PATH
	basePath := os.Getenv("PATH")

	// Check if command is specified
	if *command == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Find the command in PATH
	fullPath, err := inPath(*command, basePath)
	if err != nil {
		log.Fatal(err)
	}

	// Execute command
	fmt.Print(fullPath)
}
