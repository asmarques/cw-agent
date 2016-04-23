package stats

import (
	"fmt"

	"golang.org/x/sys/unix"
)

// GetDiskStats retrieves disk statistics for the specified path
func GetDiskStats(path string) (*DiskStats, error) {
	stat := &unix.Statfs_t{}
	err := unix.Statfs(path, stat)
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
