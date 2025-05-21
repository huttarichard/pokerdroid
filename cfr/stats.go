package cfr

import (
	"fmt"
	"time"
)

// Stats collects and reports cumulative training statistics.
type Stats struct {
	// Start is the timestamp when training started.
	Start time.Time
	// It is the number of iterations executed in the last epoch.
	It uint64
	// TotIt is the total iterations executed so far.
	TotIt uint64
	// Up is the total update count in the last epoch.
	Up uint64
	// EV is the average EV from the last epoch.
	EV float64
	// Exploit is the average exploit.
	Exploit float64
	// States is the current number of states in the game.
	States uint32
	// Nodes is the current number of nodes in the game.
	Nodes uint32
	// Epoch is the current epoch count.
	Epoch uint64
}

// String returns a nicely formatted string of the stats.
func (s *Stats) String() string {
	now := time.Now()
	diff := now.Sub(s.Start).Seconds()

	return fmt.Sprintf("ep: %-6d | it: %-10d | it/s: %-7d | up/s: %-7d | sts: %-7d | nodes: %-7d | exp: %.8f | ev: %.9f\n",
		s.Epoch,
		s.TotIt,
		uint32(float64(s.It)/diff),
		uint32(float64(s.Up)/diff),
		s.States,
		s.Nodes,
		s.Exploit,
		s.EV,
	)
}
