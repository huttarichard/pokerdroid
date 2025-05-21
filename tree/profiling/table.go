package profiling

import (
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/pokerdroid/poker/table"
)

type Formatter func(s Node, i int) string

func FormatAveragePolicy(s Node, i int) string {
	return fmt.Sprintf("%.2f%%", s.Strategy[i]*100)
}

func FormatCurrentPolicy(s Node, i int) string {
	if s.Policy == nil {
		return "-"
	}
	return fmt.Sprintf("%.2f%%", s.Policy.Strategy[i]*100)
}

func FormatRegretSummary(s Node, i int) string {
	if s.Policy == nil {
		return "-"
	}
	return fmt.Sprintf("%.10f", s.Policy.RegretSum[i])
}

func FormatStrategySummary(s Node, i int) string {
	if s.Policy == nil {
		return "-"
	}
	return fmt.Sprintf("%.10f", s.Policy.StrategySum[i])
}

func FormatBaseline(s Node, i int) string {
	if s.Policy == nil {
		return "-"
	}
	return fmt.Sprintf("%.10f", s.Policy.Baseline[i])
}

func (s Profile) WriteTable(f Formatter, w io.Writer) {
	tb := tablewriter.NewWriter(w)
	headers := []string{"Path", "Cards", "Fold", "Check", "Call", "AllIn"}

	for _, b := range s.BetSizes[0] {
		headers = append(headers, fmt.Sprintf("BET Pot %.2f", b))
	}

	tb.SetHeader(headers)
	tb.SetBorder(false)

	m := make(map[table.DiscreteAction]int)

	for i, x := range s.BetSizes[0] {
		m[table.DiscreteAction(x)] = 6 + i
	}

	for _, st := range s.Nodes {

		s := make([]string, len(headers))
		for i := range headers {
			s[i] = "-"
		}
		s[0] = st.Runes.String()
		s[1] = st.Cards.String()

		for i, a := range st.Actions {
			stx := f(st, i)
			switch a {
			case table.DFold:
				s[2] = stx
			case table.DCheck:
				s[3] = stx
			case table.DCall:
				s[4] = stx
			case table.DAllIn:
				s[5] = stx
			default:
				s[m[a]] = stx
			}
		}

		tb.Append(s)
	}

	tb.Render()
}
