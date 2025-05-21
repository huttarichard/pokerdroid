package bot

import (
	"context"
	"errors"

	"github.com/pokerdroid/poker"
	"github.com/pokerdroid/poker/table"
)

type Combined struct {
	Agents []Advisor
}

func NewCombined(agents ...Advisor) *Combined {
	return &Combined{
		Agents: agents,
	}
}

func (a *Combined) Advise(ctx context.Context, loggr poker.Logger, state State) (table.DiscreteAction, error) {
	for _, agent := range a.Agents {
		loggr.Printf("advising from agent: %T", agent)
		action, err := agent.Advise(ctx, loggr, state)
		if err != nil {
			loggr.Printf("error advising from agent: %T: %s", agent, err.Error())
			continue
		}

		loggr.Printf("action advised from agent: %T: %s", agent, action)
		return action, nil
	}

	return table.DNoAction, errors.New("no agent returned action")
}
