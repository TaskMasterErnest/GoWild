# readWord
This is a function that reads the words in a file and outputs the number of words.

## Purpose
To benchmark two different pieces of code that perform the same function.

## Prerequisites
- the Go programming language
- the Go perf package
- a web browser, preferably Firefox

## How To Run
1. First test one function, either the `useScanLn` or `useRegex`. In the main routine and the benchmark function, change the functions to match the selected one.
```go
  // main.go
	fmt.Println(useScanLn(*filename))
  
  // main_test.go
  wordCount, err := useScanLn(filename)
```
2. Run the benchmark
```bash
go test -bench . -benchtime=10x -run ^$ -benchmem | tee benchresults01m.txt
```
3. Get the cpuprofile
```bash
go test -bench . -benchtime=10x -run ^$ -cpuprofile cpu01.pprof
```
4. Run the tests again for the alternate function.
5. Use the `benchstat` tool to compare the benchmark files
```bash
benchstat benchresults01m.txt benchresults02m.txt
```
6. Use the `pprof` to gain insight into the CPU profiles
```bash
go tool pprof cpu01.pprof
```