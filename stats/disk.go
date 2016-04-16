package stats

import (
	"fmt"
	"syscall"
)

// GetDiskStats retrieves disk statistics for the specified path
func GetDiskStats(path string) (*DiskStats, error) {
	stat := &syscall.Statfs_t{}
	err := syscall.Statfs(path, stat)
	if err != nil {
		return nil, fmt.Errorf("error retrieving disk stats for '%s': %s", path, err)
	}

	stats := &DiskStats{
		MountPath:            path,
		DiskSpaceAvailable:   stat.Bfree * uint64(stat.Bsize),
		DiskSpaceUsed:        (stat.Blocks - stat.Bfree) * uint64(stat.Bsize),
		DiskSpaceUtilization: 100 * (float64((stat.Blocks - stat.Bfree)) / float64(stat.Blocks)),
	}

	return stats, nil
}
