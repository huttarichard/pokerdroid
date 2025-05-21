package poker

import (
	"fmt"
	"strings"
	"testing"
)

var (
	Name      = "PokerDroid"
	GitCommit = "unknown"
)

// Logger is minimal interface for logging.
// Complient with standard library logger.
type Logger interface {
	Printf(string, ...any)
}

type LoggerPrefix struct {
	Logger
	Prefix string
}

func (l LoggerPrefix) Printf(format string, v ...interface{}) {
	l.Logger.Printf(l.Prefix+format, v...)
}

// VoidLogger will not log anything.
type VoidLogger struct{}

func (l VoidLogger) Printf(format string, v ...interface{}) {}

// TestingLogger could be used with *testing.T.
type TestingLogger struct {
	T testing.TB
}

func (t *TestingLogger) Printf(format string, v ...interface{}) {
	t.T.Logf(format, v...)
}

type SBLogger struct {
	sb strings.Builder
}

func (l *SBLogger) Printf(format string, v ...interface{}) {
	l.sb.WriteString(fmt.Sprintf(format, v...))
}

func (l *SBLogger) String() string {
	return l.sb.String()
}

func (l *SBLogger) Bytes() []byte {
	return []byte(l.String())
}
