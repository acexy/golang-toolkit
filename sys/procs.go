package sys

import (
	"os"
	"runtime"
	"strconv"
	"strings"
)

const (
	// LimitHalf uses half the detected CPU count, rounded down with a minimum of one.
	LimitHalf = -1
	// LimitMax uses all detected CPUs.
	LimitMax = -2
	// Limit1Left reserves one detected CPU, with a minimum of one for Go.
	Limit1Left = -3
)

// CPULimitType selects a GOMAXPROCS policy.
type CPULimitType int

func readCgroupV1CPU() (quota int64, period int64) {
	quotaPath := "/sys/fs/cgroup/cpu/cpu.cfs_quota_us"
	periodPath := "/sys/fs/cgroup/cpu/cpu.cfs_period_us"
	quota = readIntFromFile(quotaPath)
	period = readIntFromFile(periodPath)
	return
}

func readCgroupV2CPU() int {
	cpuMaxPath := "/sys/fs/cgroup/cpu.max"
	content, err := os.ReadFile(cpuMaxPath)
	if err != nil {
		return 0
	}
	parts := strings.Fields(string(content))
	if len(parts) >= 2 && parts[0] != "max" {
		quota, _ := strconv.ParseInt(parts[0], 10, 64)
		period, _ := strconv.ParseInt(parts[1], 10, 64)
		if quota > 0 && period > 0 {
			return int(float64(quota) / float64(period))
		}
	}
	return 0
}

func readIntFromFile(path string) int64 {
	content, err := os.ReadFile(path)
	if err != nil {
		return -1
	}
	n, err := strconv.ParseInt(strings.TrimSpace(string(content)), 10, 64)
	if err != nil {
		return -1
	}
	return n
}

// DetectCPULimit checks cgroup CPU quotas and falls back to runtime.NumCPU.
func DetectCPULimit() int {
	if quota, period := readCgroupV1CPU(); quota > 0 && period > 0 {
		limit := int(float64(quota) / float64(period))
		if limit > 0 {
			return limit
		}
	}
	if cpus := readCgroupV2CPU(); cpus > 0 {
		return cpus
	}
	return runtime.NumCPU()
}

// SetGoMaxProc sets process-wide GOMAXPROCS when limit is positive.
func SetGoMaxProc(limit int) {
	if limit > 0 {
		runtime.GOMAXPROCS(limit)
	}
}

// SetGoMaxProcType applies a supported CPU-limit policy to GOMAXPROCS.
func SetGoMaxProcType(limit CPULimitType) {
	switch limit {
	case LimitHalf:
		num := DetectCPULimit() / 2
		if num < 1 {
			num = 1
		}
		SetGoMaxProc(num)
	case LimitMax:
		SetGoMaxProc(DetectCPULimit())
	case Limit1Left:
		num := DetectCPULimit() - 1
		if num < 1 {
			num = 1
		}
		SetGoMaxProc(num)
	}
}
