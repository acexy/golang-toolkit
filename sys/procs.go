package sys

import (
	"os"
	"runtime"
	"strconv"
	"strings"
)

const (
	LimitHalf  = -1 // 使用一半核心数，最小1，向下取整
	LimitMax   = -2 // 全部
	Limit1Left = -3 // 保留一个核心 最小1
)

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

// DetectCPULimit 检测操作系统/容器 可运行CPU核数
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

// SetGoMaxProc 设置GOMAXPROCS
func SetGoMaxProc(limit int) {
	if limit > 0 {
		runtime.GOMAXPROCS(limit)
	}
}

// SetGoMaxProcType 设置GOMAXPROCS
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
