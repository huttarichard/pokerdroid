package studiotree

import (
	"context"
	"errors"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pokerdroid/poker"
	"github.com/pokerdroid/poker/bot"
	"github.com/pokerdroid/poker/card"
	"github.com/pokerdroid/poker/chips"
	"github.com/pokerdroid/poker/frand"
	"github.com/pokerdroid/poker/table"
)

var (
	ErrGameNotFound  = errors.New("game not found")
	ErrUserNotFound  = errors.New("user not found")
	ErrInvalidAction = errors.New("invalid action")
	ErrNotPlayerTurn = errors.New("not player's turn")
)

type Advisors map[string]bot.Advisor

// GameSession represents an active poker game
type GameSession struct {
	ID          int64
	Params      table.GameParams
	Players     []bot.Advisor
	Rng         frand.Rand
	LastUpdated time.Time
	Rounds      int
	Winnings    chips.List
	LastWinners []int
	Round       int
}

//	Run runs a full heads‑up game session until the game is finished.
//
// It assumes that the session has been created through GameManager.New and that
// the GameSession contains two players (both implementing bot.Advisor),
// a valid Deck, and a Game with initial bets already made.
// It returns the net winnings (as chips.List) for each player.
func (session *GameSession) Run(ctx context.Context, logger poker.Logger) error {
	prms := session.Params.Clone()

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		logger.Printf("round %d", session.Round)

		winnings := bot.Benchmark(ctx, bot.BenchmarkParams{
			Params:   prms,
			Advisors: session.Players,
			Logger:   logger,
			Workers:  1,
			Rounds:   session.Rounds,
			Rand:     session.Rng,
			OnAction: func(state bot.State, action table.DiscreteAction) {
				session.LastUpdated = time.Now()
			},
		})

		var winners []int
		for i := range session.Players {
			session.Winnings[i] += winnings[i]

			if winnings[i] > 0 {
				winners = append(winners, i)
			}
		}

		session.LastWinners = winners

		session.Round++
	}
}

// GameManager manages poker games against bots
type GameManager struct {
	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	games  map[int64]*GameSession
	Logger poker.Logger
	Rng    frand.Rand
}

// NewGameManager creates a new game manager
func NewGameManager(logger poker.Logger) *GameManager {
	if logger == nil {
		logger = poker.VoidLogger{}
	}

	ctx, cancel := context.WithCancel(context.Background())

	gm := &GameManager{
		ctx:    ctx,
		cancel: cancel,
		games:  make(map[int64]*GameSession),
		Logger: logger,
		Rng:    frand.NewHash(),
	}

	return gm
}

var id atomic.Int64

type NewGameParams struct {
	Players   []bot.Advisor
	StackSize chips.Chips
	SB        chips.Chips
}

// CreateGame starts a new heads-up poker game against a selected bot
func (gm *GameManager) New(p NewGameParams) (int64, error) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	// Find closest root
	params := table.NewGameParams(uint8(len(p.Players)), p.StackSize)
	params.SbAmount = p.SB
	params.BtnPos = 0

	deck := card.NewDeck(card.All())
	deck.Shuffle(gm.Rng)

	id.Add(1)

	// Create session
	session := &GameSession{
		ID:          id.Load(),
		Params:      params,
		Rng:         gm.Rng,
		Rounds:      math.MaxInt,
		Players:     p.Players,
		LastUpdated: time.Now(),
		Winnings:    chips.NewListAlloc(len(p.Players)),
	}

	gm.games[id.Load()] = session

	go session.Run(gm.ctx, gm.Logger)

	return id.Load(), nil
}

func (gm *GameManager) Get(gameID int64) (*GameSession, error) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	session, ok := gm.games[gameID]
	if !ok {
		return nil, ErrGameNotFound
	}

	return session, nil
}

func (gm *GameManager) User(gameID int64) (*UserPlayer, error) {
	session, err := gm.Get(gameID)
	if err != nil {
		return nil, err
	}

	for _, player := range session.Players {
		if user, ok := player.(*UserPlayer); ok {
			return user, nil
		}
	}

	return nil, ErrUserNotFound
}

func (gm *GameManager) Action(gameID int64, action table.DiscreteAction) error {
	user, err := gm.User(gameID)
	if err != nil {
		return err
	}
	user.Action(action)
	return nil
}

func (gm *GameManager) State(gameID int64) (bot.State, error) {
	user, err := gm.User(gameID)
	if err != nil {
		return bot.State{}, err
	}
	return user.GetState(), nil
}

// UserPlayer implements bot.Advisor interface for human players
type UserPlayer struct {
	lastState bot.State
	actionCh  chan table.DiscreteAction
}

func NewUserPlayer() *UserPlayer {
	return &UserPlayer{
		actionCh: make(chan table.DiscreteAction, 1),
	}
}

// Advise implements bot.Advisor interface
func (p *UserPlayer) Advise(ctx context.Context, loggr poker.Logger, state bot.State) (table.DiscreteAction, error) {
	p.lastState = state

	// Wait for client action
	select {
	case action := <-p.actionCh:
		return action, nil
	case <-ctx.Done():
		return table.DNoAction, ctx.Err()
	}
}

// Action allows client to submit an action
func (p *UserPlayer) Action(action table.DiscreteAction) {
	p.actionCh <- action
}

// GetState allows client to get current game state
func (p *UserPlayer) GetState() bot.State {
	return p.lastState
}
