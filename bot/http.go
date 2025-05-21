package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pokerdroid/poker"
	"github.com/pokerdroid/poker/table"
)

type HTTP struct {
	Advisor Advisor
	Logger  poker.Logger
}

func NewHTTP(advisor Advisor) *HTTP {
	return &HTTP{
		Advisor: advisor,
		Logger:  poker.VoidLogger{},
	}
}

type HttpResponse struct {
	Action table.DiscreteAction `json:"action"`
}

func (a HTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bd, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	a.Logger.Printf("received request: %d bytes", len(bd))

	state := State{}

	err = state.UnmarshalBinary(bd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.Logger.Printf("advising from http state: %s", state.String())

	action, err := a.Advisor.Advise(r.Context(), a.Logger, state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	actionBytes, err := json.Marshal(HttpResponse{
		Action: action,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(actionBytes)
}

func Request(url string, state State) (table.DiscreteAction, error) {
	data, err := state.MarshalBinary()
	if err != nil {
		return table.DNoAction, err
	}

	resp, err := http.Post(url, "application/octet-stream", bytes.NewReader(data))
	if err != nil {
		return table.DNoAction, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return table.DNoAction, err
	}

	if resp.StatusCode != http.StatusOK {
		return table.DNoAction, fmt.Errorf("err %d: %s", resp.StatusCode, string(body))
	}

	var res HttpResponse

	err = json.Unmarshal(body, &res)
	if err != nil {
		return table.DNoAction, fmt.Errorf("decode response: %w", err)
	}

	return res.Action, nil
}
