// +build !linux

package stats

import (
	"fmt"
	"runtime"
)

// GetMemoryStats returns memory and swap statistics
func GetMemoryStats(sharedBuffersAsUsed bool) (*MemoryStats, error) {
	return nil, fmt.Errorf("memory usage not supported for OS '%s'", runtime.GOOS)
}
