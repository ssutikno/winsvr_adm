package network

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/net"
)

// NetStats holds network interface statistics
type NetStats struct {
	Interface string
	RxBytes   uint64
	TxBytes   uint64
}

// GetNetworkUtilization returns the network utilization for all interfaces
func GetNetworkUtilization() ([]NetStats, error) {
	ioCounters, err := net.IOCounters(true)
	if err != nil {
		return nil, fmt.Errorf("failed to get network IO counters: %w", err)
	}

	var netStats []NetStats
	for _, io := range ioCounters {
		netStats = append(netStats, NetStats{
			Interface: io.Name,
			RxBytes:   io.BytesRecv,
			TxBytes:   io.BytesSent,
		})
	}

	return netStats, nil
}

// CalculateNetworkUtilization calculates the network utilization in bytes per second
func CalculateNetworkUtilization(prevStats, currStats []NetStats, duration time.Duration) ([]NetStats, error) {
	if len(prevStats) != len(currStats) {
		return nil, fmt.Errorf("mismatched number of interfaces")
	}

	var utilization []NetStats
	for i := range currStats {
		if prevStats[i].Interface != currStats[i].Interface {
			return nil, fmt.Errorf("mismatched interface names")
		}

		rxBytesPerSecond := float64(currStats[i].RxBytes-prevStats[i].RxBytes) / duration.Seconds()
		txBytesPerSecond := float64(currStats[i].TxBytes-prevStats[i].TxBytes) / duration.Seconds()

		utilization = append(utilization, NetStats{
			Interface: currStats[i].Interface,
			RxBytes:   uint64(rxBytesPerSecond),
			TxBytes:   uint64(txBytesPerSecond),
		})
	}

	return utilization, nil
}
