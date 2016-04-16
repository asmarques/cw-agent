package stats

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

const memInfoPath = "/proc/meminfo"

var memInfoExp = regexp.MustCompile(`^(.+):\s+(\d+)`)

// GetMemoryStats returns memory and swap statistics
func GetMemoryStats(sharedBuffersAsUsed bool) (*MemoryStats, error) {
	file, err := os.Open(memInfoPath)
	if err != nil {
		return nil, fmt.Errorf("error retrieving memory stats from %s: %s", memInfoPath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	mem := map[string]uint64{}

	for scanner.Scan() {
		result := memInfoExp.FindStringSubmatch(scanner.Text())
		if result != nil && len(result) == 3 {
			value, err := strconv.ParseUint(result[2], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing %s: %s", memInfoPath, err)
			}
			mem[result[1]] = value * 1024
		} else {
			return nil, fmt.Errorf("invalid format for %s", memInfoPath)
		}
	}

	freeMem := mem["MemFree"]
	if !sharedBuffersAsUsed {
		freeMem += mem["Buffers"] + mem["Cached"]
	}

	stats := &MemoryStats{
		MemoryAvailable:   freeMem,
		MemoryUsed:        mem["MemTotal"] - freeMem,
		MemoryUtilization: 100 * float64((mem["MemTotal"] - freeMem)) / float64(mem["MemTotal"]),
		SwapAvailable:     mem["SwapFree"],
		SwapUsed:          mem["SwapTotal"] - mem["SwapFree"],
	}

	if mem["SwapTotal"] > 0 {
		stats.SwapUtilization = 100 * float64(stats.SwapUsed) / float64(mem["SwapTotal"])
	}

	return stats, nil
}
