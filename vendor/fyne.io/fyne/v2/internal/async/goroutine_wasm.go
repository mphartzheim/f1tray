//go:build wasm

package async

// IsMainGoroutine returns true if it is called from the main goroutine, false otherwise.
func IsMainGoroutine() bool {
	return true
}
