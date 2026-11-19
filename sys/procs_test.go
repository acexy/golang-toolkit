package sys

import (
	"math"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestDetectCPULimit(t *testing.T) {
	for _, tc := range []struct {
		name     string
		v1Quota  string
		v1Period string
		v2Max    string
		want     int
	}{
		{"v1 subcore takes precedence", "50000", "100000", "200000 100000", 1},
		{"v2 subcore", "-1", "100000", "50000 100000", 1},
		{"v1 rounds down", "250000", "100000", "max 100000", 2},
		{"v2 rounds down", "-1", "100000", "250000 100000", 2},
		{"large quota fits int", "9223372036854775807", "1", "max 100000", math.MaxInt},
		{"v2 after invalid v1 period", "200000", "0", "300000 100000", 3},
		{"unlimited", "-1", "100000", "max 100000", runtime.NumCPU()},
		{"invalid", "invalid", "100000", "9223372036854775808 100000", runtime.NumCPU()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			cpuDir := filepath.Join(dir, "cpu")
			if err := os.Mkdir(cpuDir, 0700); err != nil {
				t.Fatal(err)
			}
			for path, content := range map[string]string{
				filepath.Join(cpuDir, "cpu.cfs_quota_us"):  tc.v1Quota,
				filepath.Join(cpuDir, "cpu.cfs_period_us"): tc.v1Period,
				filepath.Join(dir, "cpu.max"):              tc.v2Max,
			} {
				if err := os.WriteFile(path, []byte(content+"\n"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			if got := detectCPULimit(dir); got != tc.want {
				t.Fatalf("detectCPULimit() = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestSetGOMAXPROCS(t *testing.T) {
	t.Cleanup(runtime.SetDefaultGOMAXPROCS)
	SetGoMaxProc(1)
	if got := runtime.GOMAXPROCS(0); got != 1 {
		t.Fatalf("GOMAXPROCS = %d, want 1", got)
	}
	SetGoMaxProc(0)
	if got := runtime.GOMAXPROCS(0); got != 1 {
		t.Fatalf("GOMAXPROCS after zero limit = %d, want 1", got)
	}
	want := DetectCPULimit()
	SetGoMaxProcType(LimitMax)
	if got := runtime.GOMAXPROCS(0); got != want {
		t.Fatalf("GOMAXPROCS after LimitMax = %d, want %d", got, want)
	}
}
