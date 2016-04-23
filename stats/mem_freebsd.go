package stats

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	totalMemKey      = "hw.realmem"
	pageSizeKey      = "hw.pagesize"
	pagesInactiveKey = "vm.stats.vm.v_inactive_count"
	pagesCacheKey    = "vm.stats.vm.v_cache_count"
	pagesFreeKey     = "vm.stats.vm.v_free_count"
	swapInfoKey      = "vm.swap_info"
	swapInfoSize     = 20
	swapInfoVersion  = 1
)

func sysctlError(key string, err error) error {
	return fmt.Errorf("error reading sysctl %s: %s", key, err)
}

// GetMemoryStats returns memory and swap statistics
func GetMemoryStats() (*MemoryStats, error) {
	totalMem, err := unix.SysctlUint64(totalMemKey)
	if err != nil {
		return nil, sysctlError(totalMemKey, err)
	}

	pageSize, err := unix.SysctlUint32(pageSizeKey)
	if err != nil {
		return nil, sysctlError(pageSizeKey, err)
	}

	pagesInactive, err := unix.SysctlUint32(pagesInactiveKey)
	if err != nil {
		return nil, sysctlError(pagesInactiveKey, err)
	}

	pagesCache, err := unix.SysctlUint32(pagesCacheKey)
	if err != nil {
		return nil, sysctlError(pagesCacheKey, err)
	}

	pagesFree, err := unix.SysctlUint32(pagesFreeKey)
	if err != nil {
		return nil, sysctlError(pagesFreeKey, err)
	}

	n := 0
	var swapInfoUsed, swapInfoTotal uint64

	// Iterate over swap devices and sum used and total pages
	for {
		swapInfo, err := unix.SysctlRaw(swapInfoKey, n)
		if err != nil {
			if os.IsNotExist(err) {
				// Swap device not found, stop searching for swap devices
				break
			} else {
				return nil, sysctlError(swapInfoKey, err)
			}
		}

		if len(swapInfo) != swapInfoSize {
			return nil, fmt.Errorf("unexpected size for %s", swapInfoKey)
		}

		// The value of the vm.swap_info sysctl matches a struct with the following layout:
		// struct xswdev {
		//     u_int   xsw_version;
		//     dev_t   xsw_dev;
		//     int     xsw_flags;
		//     int     xsw_nblks;
		//     int     xsw_used;
		// };

		values := *(*[5]int32)(unsafe.Pointer(&swapInfo[0])) // xsw_version

		if values[0] != swapInfoVersion {
			return nil, fmt.Errorf("unexpected version %d for %s", values[0], swapInfoKey)
		}

		swapInfoTotal = swapInfoTotal + uint64(values[3]) // xsw_nblks
		swapInfoUsed = swapInfoUsed + uint64(values[4])   // xsw_used
		n++
	}

	memoryAvailable := uint64(pageSize) * (uint64(pagesInactive) + uint64(pagesCache) + uint64(pagesFree))
	memoryUsed := totalMem - memoryAvailable
	memoryUtilization := 100 * (float64(memoryUsed) / float64(totalMem))

	swapAvailable := uint64(pageSize) * (swapInfoTotal - swapInfoUsed)
	swapUsed := uint64(pageSize) * swapInfoUsed

	stats := &MemoryStats{
		MemoryAvailable:   memoryAvailable,
		MemoryUsed:        memoryUsed,
		MemoryUtilization: memoryUtilization,
		SwapAvailable:     swapAvailable,
		SwapUsed:          swapUsed,
	}

	if swapInfoTotal > 0 {
		stats.SwapUtilization = 100 * (float64(swapInfoUsed) / float64(swapInfoTotal))
	}

	return stats, nil
}
