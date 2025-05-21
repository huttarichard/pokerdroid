package slumbot

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/table"
)

const (
	Host       = "slumbot.com"
	NumStreets = 4
	SmallBlind = 50
	BigBlind   = 100
	StackSize  = 20000
)

var split = regexp.MustCompile(`(c|b\d+|k|f|/)`)

type PlayerAction struct {
	table.ActionAmount
	Player uint8
}

func (p PlayerAction) String() string {
	return fmt.Sprintf("#P%d: %s $%s", p.Player, p.Action.String(), p.Amount.StringFixed(2))
}

type PlayerActions []PlayerAction

func (pp PlayerActions) String() string {
	var ss strings.Builder
	ss.WriteString("Player Actions:\n")
	for _, a := range pp {
		ss.WriteString("\t" + a.String() + "\n")
	}
	return ss.String()
}

type State struct {
	OldAction    string   `json:"old_action"`
	ErrMsg       *string  `json:"error_msg"`
	Action       string   `json:"action"`
	ClientPos    int      `json:"client_pos"`
	HoleCards    []string `json:"hole_cards"`
	Board        []string `json:"board"`
	BotHoleCards []string `json:"bot_hole_cards,omitempty"`
	WonPot       float64  `json:"won_pot,omitempty"`
}

type HandResponse struct {
	State   `json:"state,inline"`
	Token   string        `json:"token"`
	Actions PlayerActions `json:"-"`
}

func (s State) GetHoleCards() card.Cards {
	var ppc card.Cards
	for _, c := range s.HoleCards {
		ppc = append(ppc, card.Parse(c))
	}
	return ppc
}

func (s State) GetBotHoleCards() card.Cards {
	var ppc card.Cards
	for _, c := range s.BotHoleCards {
		ppc = append(ppc, card.Parse(c))
	}
	return ppc
}

func (s State) GetBoardCards() card.Cards {
	var ppc card.Cards
	for _, c := range s.Board {
		ppc = append(ppc, card.Parse(c))
	}
	return ppc
}

type ActRequest struct {
	Token string `json:"token"`
	Incr  string `json:"incr"`
}

type ActResponse struct {
	Token string `json:"token"`
	State State  `json:"state"`
}

type ErrorResponse struct {
	ErrorMsg string `json:"error_msg"`
}

func EncodeAction(a table.ActionKind, amount chips.Chips) string {
	switch a {
	case table.Fold:
		return "f"
	case table.Bet, table.Raise, table.AllIn:
		return "b" + amount.StringFixed(0)
	case table.Call:
		return "c"
	case table.Check:
		return "k"
	default:
		panic("illegal action")
	}
}

func ParseActionString(actionStr string, clientPos int) []PlayerAction {
	splitActions := split.FindAllString(actionStr, -1)

	var actions []PlayerAction

	btnPos := uint8(1)
	if clientPos == 1 {
		btnPos = uint8(0)
	}
	currentPlayer := btnPos
	nextPlayer := (btnPos + 1) % 2

	street := 0
	accumulator := [NumStreets]int{}

	var bets [NumStreets][2]int
	sb := SmallBlind
	bb := BigBlind

	if btnPos == 0 {
		bets[street][0] += sb
		bets[street][1] += bb
	} else {
		bets[street][0] += bb
		bets[street][1] += sb
	}

	for _, action := range splitActions {
		if action == "/" {
			street++
			continue
		}

		if street >= 1 && accumulator[street] == 0 {
			currentPlayer = (btnPos + 1) % 2
			nextPlayer = btnPos
		}

		switch {
		case action[0] == 'c':
			amount := bets[street][nextPlayer] - bets[street][currentPlayer]
			bets[street][currentPlayer] += amount

			actions = append(actions, PlayerAction{
				ActionAmount: table.ActionAmount{
					Action: table.Call,
					Amount: chips.NewFromInt(int64(amount)),
				},
				Player: currentPlayer,
			})
		case action[0] == 'b':
			actx := table.Raise
			sameBets := bets[street][currentPlayer] == bets[street][nextPlayer]
			if sameBets {
				actx = table.Bet
			}
			amount, _ := strconv.Atoi(action[1:])
			amount -= bets[street][currentPlayer]
			bets[street][currentPlayer] += amount

			total := 0
			for _, b := range bets {
				total += b[currentPlayer]
			}
			if total == StackSize {
				actx = table.AllIn
			}

			actions = append(actions, PlayerAction{
				ActionAmount: table.ActionAmount{
					Action: actx,
					Amount: chips.NewFromInt(int64(amount)),
				},
				Player: currentPlayer,
			})
		case action[0] == 'k':
			actions = append(actions, PlayerAction{
				ActionAmount: table.ActionAmount{
					Action: table.Check,
					Amount: chips.Zero,
				},
				Player: currentPlayer,
			})
		case action[0] == 'f':
			actions = append(actions, PlayerAction{
				ActionAmount: table.ActionAmount{
					Action: table.Fold,
					Amount: chips.Zero,
				},
				Player: currentPlayer,
			})
		}
		accumulator[street] += 1

		nextPlayer = currentPlayer
		currentPlayer = (currentPlayer + 1) % 2
	}

	return actions
}

// def get_parsed_actions_new(
//     action_str: str, old_str: str, client_pos: int
// ) -> List[Tuple[int, Action, Chips]]:
//     a = parse_action_string(action_str, client_pos)
//     b = parse_action_string(old_str, client_pos)
//     return a[len(b) :], a, b

func ParseActionStringDiff(a, o string, cp int) []PlayerAction {
	ax := ParseActionString(a, cp)
	ox := ParseActionString(o, cp)
	return ax[len(ox):]
}

func ParseResponse(r *http.Response) (HandResponse, error) {
	var handResp HandResponse
	var state State
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return HandResponse{}, err
	}
	err = json.Unmarshal(data, &state)
	if err != nil {
		return HandResponse{}, err
	}
	if state.ErrMsg != nil {
		return HandResponse{}, errors.New(*state.ErrMsg)
	}
	handResp.State = state

	handResp.Actions = ParseActionStringDiff(
		state.Action,
		state.OldAction,
		state.ClientPos,
	)
	return handResp, nil
}

func NewHand(ctx context.Context, token string) (HandResponse, error) {
	data := map[string]string{}
	if token != "" {
		data["token"] = token
	}

	b, err := json.Marshal(data)
	if err != nil {
		return HandResponse{}, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("https://%s/api/new_hand", Host),
		bytes.NewBuffer(b),
	)
	if err != nil {
		return HandResponse{}, err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return HandResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return HandResponse{}, err
		}
		return HandResponse{}, fmt.Errorf("non success response: %s", errResp.ErrorMsg)
	}

	handResp, err := ParseResponse(resp)
	if err != nil {
		return HandResponse{}, err
	}

	return handResp, nil
}

func Act(ctx context.Context, token string, action string, hand HandResponse) (HandResponse, error) {
	actReq := ActRequest{
		Token: token,
		Incr:  action,
	}

	payload, err := json.Marshal(actReq)
	if err != nil {
		return HandResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("https://%s/api/act", Host), nil)
	if err != nil {
		return HandResponse{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Body = ioutil.NopCloser(bytes.NewBuffer(payload))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return HandResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return HandResponse{}, err
		}
		var errResp ErrorResponse
		err = json.Unmarshal(data, &errResp)
		if err != nil {
			return HandResponse{}, err
		}
		return HandResponse{}, fmt.Errorf("non success response: %s", errResp.ErrorMsg)
	}

	actResp, err := ParseResponse(resp)
	if err != nil {
		return HandResponse{}, err
	}

	return actResp, nil
}

func Login(username, password string) (string, error) {
	loginData := map[string]string{
		"username": username,
		"password": password,
	}

	payload, err := json.Marshal(loginData)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(fmt.Sprintf("https://%s/api/login", Host), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return "", err
		}
		return "", fmt.Errorf("non success response: %s", errResp.ErrorMsg)
	}

	var loginResp struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return "", err
	}

	return loginResp.Token, nil
}
