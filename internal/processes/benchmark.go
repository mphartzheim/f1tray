package processes

import (
	"fmt"
	"time"

	"github.com/mphartzheim/f1tray/internal/config"
)

var benchmarks = make(map[string]time.Time)

// StartBenchmark records the current time for a given label if debug mode is enabled.
func StartBenchmark(label string) {
	if !config.Get().Debug.Enabled {
		return
	}
	benchmarks[label] = time.Now()
}

// EndBenchmark prints the elapsed time since StartBenchmark for the given label.
func EndBenchmark(label string) {
	if !config.Get().Debug.Enabled {
		return
	}
	start, exists := benchmarks[label]
	if !exists {
		fmt.Printf("[BENCHMARK] %s not found\n", label)
		return
	}
	elapsed := time.Since(start)
	fmt.Printf("[BENCHMARK] %s: %v\n", label, elapsed)
	delete(benchmarks, label)
}
