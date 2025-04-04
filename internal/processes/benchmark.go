package processes

import (
	"fmt"
	"time"

	"github.com/mphartzheim/f1tray/internal/config"
)

var benchmarks = make(map[string]time.Time)

func StartBenchmark(label string) {
	if !config.Get().Debug.Enabled {
		return
	}
	benchmarks[label] = time.Now()
}

func EndBenchmark(label string) {
	if !config.Get().Debug.Enabled {
		return
	}
	if start, ok := benchmarks[label]; ok {
		duration := time.Since(start)
		fmt.Printf("[BENCHMARK] %s: %s\n", label, duration.Round(time.Millisecond))
	}
}
