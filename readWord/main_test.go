package main

import (
	"path/filepath"
	"testing"
)

func BenchmarkReadWords(b *testing.B) {
	filenames, err := filepath.Glob("./testdata/benchmark/*.csv")
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		for _, filename := range filenames {
			b.StopTimer() // Stop the timer to exclude setup time
			wordCount, err := useScanLn(filename)
			b.StartTimer() // Restart the timer before measuring `useBufio` time

			if err != nil {
				b.Error(err)
			}

			_ = wordCount // we don’t use wordCount in the benchmark, but it’s computed
		}
	}
}
