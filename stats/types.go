package stats

// MemoryStats represents snapshot of memory statistics at a given point in time
type MemoryStats struct {
	MemoryAvailable   uint64
	MemoryUsed        uint64
	MemoryUtilization float64
	SwapAvailable     uint64
	SwapUsed          uint64
	SwapUtilization   float64
}

// DiskStats represents a snapshot of disk statistics for a path at a given point in time
type DiskStats struct {
	MountPath            string
	DiskSpaceAvailable   uint64
	DiskSpaceUsed        uint64
	DiskSpaceUtilization float64
}
