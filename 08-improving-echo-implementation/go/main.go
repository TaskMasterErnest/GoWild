package main

import (
	"flag"
	"fmt"
	"os/exec"
)

func main() {
	// get input from Stdin
	input := flag.String("i", "", "string to pass")
	flag.Parse()

	cmd := exec.Command("/usr/bin/printf", *input)

	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(output))
}
