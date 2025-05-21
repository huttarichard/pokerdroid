package memtils

import (
	"context"
	"log"
	"os"
	"runtime"
	"syscall"
	"time"
)

type MemStats struct {
	Runtime *runtime.MemStats
	System  *syscall.Rusage
}

func GetMemStats() *MemStats {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	var rusage syscall.Rusage
	syscall.Getrusage(syscall.RUSAGE_SELF, &rusage)

	return &MemStats{
		Runtime: &ms,
		System:  &rusage,
	}
}

func MonitorMemory(ctx context.Context, interval time.Duration, logger *log.Logger) {
	pid := os.Getpid()
	lastPageFaults := int64(0)

	report := func() {
		stats := GetMemStats()

		// Calculate page fault delta
		pageFaultDelta := stats.System.Minflt + stats.System.Majflt - lastPageFaults
		lastPageFaults = stats.System.Minflt + stats.System.Majflt

		logger.Printf("\n=== Memory Stats (PID: %d) ===", pid)
		logger.Printf("Heap In-Use: %v MB", stats.Runtime.HeapInuse/1024/1024)
		logger.Printf("Heap Allocated: %v MB", stats.Runtime.HeapAlloc/1024/1024)
		logger.Printf("Total Allocated: %v MB", stats.Runtime.TotalAlloc/1024/1024)
		logger.Printf("System Memory: %v MB", stats.Runtime.Sys/1024/1024)
		logger.Printf("GC Cycles: %v", stats.Runtime.NumGC)
		logger.Printf("Page Faults (cumulative): %v", stats.System.Minflt+stats.System.Majflt)
		logger.Printf("Page Faults (delta): %v", pageFaultDelta)
		logger.Printf("Max RSS: %v MB\n", stats.System.Maxrss/(1024*1024))
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			report()
		}
	}
}
