package main

import (
	"encoding/json"
	"net/http"

	// "network"
	// "storage"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"

	// "github.com/shirou/gopsutil/v4/process"

	// import modules from the network and storage packages
	"modules/network"
	"modules/storage"
	"modules/winprocess"
)

type ServerStatus struct {
	CPUUsage     float64                  `json:"cpuUsage"`
	MemoryUsage  float64                  `json:"memoryUsage"`
	Processes    []winprocess.ProcessInfo `json:"processes"`
	Network      []network.NetStats       `json:"network"`
	StorageUsage []storage.UsageStat      `json:"storage"`
}

func getServerStatus() ServerStatus {
	cpuPercentages, _ := cpu.Percent(0, false)
	memoryStat, _ := mem.VirtualMemory()
	processes, _ := winprocess.GetProcesses()
	networkStats, _ := network.GetNetworkUtilization()
	storageStats, _ := storage.GetStorageUsage()

	// convert the process list to the correct type
	var processList []winprocess.ProcessInfo
	for _, p := range processes {
		processList = append(processList, *p)
	}

	return ServerStatus{
		CPUUsage:     cpuPercentages[0],
		MemoryUsage:  memoryStat.UsedPercent,
		StorageUsage: storageStats,
		Network:      networkStats,
		Processes:    processList,
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	status := getServerStatus()
	json.NewEncoder(w).Encode(status)
}

// write server memory utilization function and return json
func cpuUtilization(w http.ResponseWriter, r *http.Request) {
	cpuPercentages, _ := cpu.Percent(0, false)

	json.NewEncoder(w).Encode(cpuPercentages[0])
}

// write network utilization function and return json
func networkUtilization(w http.ResponseWriter, r *http.Request) {
	// prepare the variable
	var netStats []network.NetStats
	// var netUsage []NetStats

	// get network utilization
	netStats, _ = network.GetNetworkUtilization()

	// return json
	json.NewEncoder(w).Encode(netStats)
}

// write network calculation function and return json
func networkCalculation(w http.ResponseWriter, r *http.Request) {
	// prepare the variable
	var netStats []network.NetStats
	// var netUsage []network.NetStats

	// get network utilization
	netStats, _ = network.GetNetworkUtilization()

	// return json
	json.NewEncoder(w).Encode(netStats)
}

// write storage usage function and return json
func storageUsage(w http.ResponseWriter, r *http.Request) {
	// prepare the variable
	var usageStats []storage.UsageStat

	// get storage usage
	usageStats, _ = storage.GetStorageUsage()

	// return json
	json.NewEncoder(w).Encode(usageStats)
}

// write memory usage function and return json
func memoryUsage(w http.ResponseWriter, r *http.Request) {
	// prepare the variable
	var memoryStat *mem.VirtualMemoryStat

	// get memory usage
	memoryStat, _ = mem.VirtualMemory()

	// return json
	json.NewEncoder(w).Encode(memoryStat)
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	// prepare the variable
	var processes []*winprocess.ProcessInfo

	// get processes
	processes, _ = winprocess.GetProcesses()

	// return json
	json.NewEncoder(w).Encode(processes)
}

func main() {
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/cpu", cpuUtilization)
	http.HandleFunc("/network", networkCalculation)
	http.HandleFunc("/storage", storageUsage)
	http.HandleFunc("/memory", memoryUsage)
	http.HandleFunc("/processes", processHandler)

	// serving static files
	http.Handle("/", http.FileServer(http.Dir("./static")))

	println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
	// write on the console that the server is running

}
