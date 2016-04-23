// +build !linux,!freebsd

package stats

import (
	"fmt"
	"runtime"
)

// GetMemoryStats returns memory and swap statistics
func GetMemoryStats() (*MemoryStats, error) {
	return nil, fmt.Errorf("memory usage not supported for OS '%s'", runtime.GOOS)
}
