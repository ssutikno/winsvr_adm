package storage

import (
	"fmt"

	"github.com/shirou/gopsutil/disk"
)

// UsageStat represents the storage usage statistics for a partition.
type UsageStat struct {
	Path        string  `json:"path"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
}

// GetStorageUsage retrieves and returns storage usage information in JSON format.
func GetStorageUsage() ([]UsageStat, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, fmt.Errorf("failed to get disk partitions: %w", err)
	}

	var usageStats []UsageStat
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to get disk usage for %s: %w", partition.Mountpoint, err)
		}

		usageStats = append(usageStats, UsageStat{
			Path:        usage.Path,
			Total:       usage.Total,
			Used:        usage.Used,
			Free:        usage.Free,
			UsedPercent: usage.UsedPercent,
		})
	}

	// jsonData, err := json.Marshal(usageStats)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to marshal storage usage data: %w", err)
	// }

	return usageStats, nil
}
