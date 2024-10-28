# readWord: Go Word Counter Performance Analysis

A performance comparison of two different approaches to counting words in text files using Go. This project provides detailed benchmarking and profiling of different word counting implementations.

## Purpose

To benchmark and analyze two different implementations that perform the same function:
1. Scanner-based approach using `bufio.Scanner`
2. Regex-based approach using `regexp.MustCompile`

## Prerequisites

- Go programming language installed
- Go perf package
- Firefox browser (preferred for viewing profiles)
- `benchstat` tool for comparing benchmark results

## Implementation Details

### Scanner Implementation (`useScanWords`)
```go
scanner := bufio.NewScanner(f)
scanner.Split(bufio.ScanWords)
```
This implementation uses Go's `bufio.Scanner` with `ScanWords` as the split function, providing a streamlined approach to word tokenization.

### Regex Implementation (`useRegex`)
```go
reg := regexp.MustCompile(`\S+`)
```
This implementation uses regular expressions to match non-whitespace sequences, offering a more flexible but potentially more resource-intensive approach.

## How to Run

1. **Select Implementation**: Choose either `useScanWords` or `useRegex` by modifying the main routine and benchmark function:
```go
// main.go
fmt.Println(useScanWords(*filename))

// main_test.go
wordCount, err := useScanWords(filename)
```

2. **Run Benchmarks**: Execute the benchmark with memory statistics:
```bash
go test -bench . -benchtime=10x -run ^$ -benchmem | tee benchresults01m.txt
```

3. **Generate CPU Profile**:
```bash
go test -bench . -benchtime=10x -run ^$ -cpuprofile cpu01.pprof
```

4. **Test Alternative Implementation**: Repeat steps 1-3 with the other function

5. **Compare Results**: Use the `benchstat` tool to compare the benchmark files:
```bash
benchstat benchresults01m.txt benchresults02m.txt
```

6. **Analyze CPU Profile**: Use the `pprof` tool for detailed CPU analysis:
```bash
go tool pprof cpu01.pprof
```

## Usage

Basic command-line usage:
```bash
go run main.go -f <filename>
```

### Flags
- `-f`: Specify the input file path (required)

## Project Structure
```
.
├── main.go         # Main implementation file
├── main_test.go    # Benchmark tests
├── benchresults*.txt  # Benchmark results
├── cpu*.pprof     # CPU profiles
└── testdata/      
    └── benchmark/  # CSV files for benchmarking
```

## Performance Considerations

When choosing between these implementations, consider:
1. File size - Scanner implementation may be more memory-efficient for large files
2. Word pattern complexity - Regex implementation offers more flexibility for complex patterns
3. Processing requirements - Scanner implementation generally provides better performance for simple word counting

## Benchmarking Details

The benchmark suite:
- Processes multiple CSV files from the `testdata/benchmark` directory
- Excludes setup time from measurements using `b.StopTimer()` and `b.StartTimer()`
- Runs multiple iterations for statistical significance
- Includes memory allocation statistics
- Generates CPU profiles for detailed performance analysis
