Task:

We need to create mapping function that will map real game actions to the abstracted tree actions.


Summary of Scenarios and Problems

Scenario 1: Bet Size Mismatch
Initial Stack: 100 units
Real Game Path (R): r:n:b3.00:c:n
Player bets 3.00 units.
Abstracted Tree Path (A): r:n:b2.00:c:n 
Player can only bet 2.00 units in the tree.

Problem:
- Mismatch in Bet Sizes: Real game bet doesn't match any available bet size in the tree.
- Mapping Challenge: No direct node corresponds to the real game action, making mapping difficult.


Scenario 2: Stack Depth Divergence Due to Sequential Betting
Initial Stack: 50 units
Real Game Path (R): r:n:b2.00:c:n:b6.00:c:n:b15.00:c:n
Player bets: 2.00 ➔ 6.00 ➔ 15.00 units.
Abstracted Tree Path (A): r:n:b2.00:c:n:allin
Tree forces an all-in due to limited options.

Problem:
- Early All-In: Tree lacks intermediate bet sizes.
- Strategic Divergence: Alters strategic progression.


Scenario 3: Pot Size Divergence Due to Multi-Street Betting
Initial Stack: 75 units
Real Game Path (R): r:n:b3.00:c:n:b12.00:c:n:b30.00
Increasing bets over streets.
Abstracted Tree Path (A): r:n:b2.00:c:n:b2.00:c:n:allin
Only small bets and an all-in.

Problem:
- Different Pot Growth: Pot size diverges significantly.
- Strategy Impact: Affects decisions due to pot odds.


Scenario 4: Min-Raise Mismatch
Initial Stack: 40 units
Real Game Path (R): r:n:b2.00:c:n:b4.00:r8.00:c
Player raises to 8.00 units.
Abstracted Tree Path (A): r:n:b2.00:c:n:allin
No intermediate raise options.

Problem:
- Missing Raise Options: Tree lacks certain legal actions.
- Mapping Limitation: Cannot match real game raises.


Scenario 5: Stack-to-Pot Ratio (SPR) Divergence
Initial Stack: 60 units
Real Game Path (R): r:n:b3.00:c:n:b12.00:c:n
SPR after flop: ~15.
Abstracted Tree Path (A): r:n:b2.00:c:n:b2.00:c:n
SPR after flop: ~20.

Problem:
- Different SPRs: Affects strategic decisions.
- Strategic Alignment Issue: Tree strategy may not be effective.


Scenario 6: Multi-Way Pot Abstraction
Initial Stack: 80 units
Real Game Path (R): r:n:b2.50:c:r10.00:c:n
Significant raise leads to larger pot.
Abstracted Tree Path (A): r:n:b2.00:c:allin:c:n
Simplifies betting to early all-in.

Problem:
- Pot Size and Betting Complexity: Tree doesn't capture dynamics.
- Mapping Inaccuracy: Strategic mismatch.


Scenario 7: Betting Pattern Divergence Across Streets
Initial Stack: 45 units
Real Game Path (R): r:n:b2.00:c:n:k:k:n:b15.00
Checks early, large bet later.
Abstracted Tree Path (A): r:n:b2.00:c:n:k:k:n:b2.00
Only small bet allowed later.

Problem:
- Limited Betting Options on Later Streets: Can't represent aggressive plays.
- Strategic Restriction: Affects validity.


Scenario 8: Late Street Stack Depth Mismatch
Initial Stack: 70 units
Real Game Path (R): r:n:b3.00:c:n:b9.00:c:n:b27.00
Large bets culminating in 27.00 units.
Abstracted Tree Path (A): r:n:b2.00:c:n:b2.00:c:n:allin
Forced all-in, no incremental bets.

Problem:
- Lack of Deep Stack Play Representation: Can't mirror deep stack strategies.
- Strategic Divergence: Affects gameplay.


Scenario 9: Bet Size Decrease Across Streets
Real Game Path (R): r:n:b10.00:c:n:b2.00:c:n:b20.00
Bet sizes fluctuate across streets.
Abstracted Tree Path (A): r:n:b2.00:c:n:b2.00:c:n:allin
Doesn't accommodate fluctuations.

Problem:
- Fluctuating Bet Sizes Not Supported: Can't represent legal bet size changes.
- Mapping Challenge: Strategic mismatch.





Defining the Overall Problem

The primary challenge is accurately mapping real poker game actions onto an abstracted game tree used for Counterfactual Regret Minimization (CFR) algorithms. Discrepancies arise due to:

1) Discrepancies in Betting Options:
- Limited Bet Sizes in the Tree: May not cover all real game bet amounts.
- Missing Intermediate Actions: Certain legal actions are absent.

2) Divergence in Game State Metrics:
- Stack Sizes: Differing bet sizes lead to varying remaining stacks.
- Pot Sizes: Cumulative bets cause pot size divergence.

3) Strategic Differences:
- SPR Variance: Affects optimal strategies.
- Aggression Levels: Can't represent aggressive plays.

4) Complex Betting Patterns:
- Multi-Way Pots: Tree may not reflect dynamics accurately.
- Cross-Street Variations: Bet fluctuations not always represented.

5) Mapping Challenges:
- No Direct Correspondence: Hard to map game state to tree.
- Potential for Mismatches: Incorrect mapping leads to errors.



Solution:
///////////////////////////////////////////
// Main Entry Point
///////////////////////////////////////////

function map_game_state_to_tree(game_state, tree_root)
    current_tree_node = tree_root
    
    for action in game_state.action_history
        effective_stack = calculate_effective_stack(game_state)
        pot_size = calculate_pot_size(game_state)
        
        // Map action to a corresponding tree action
        if action.type == 'bet' or action.type == 'raise'
            mapped_bet_size = map_bet_size_to_tree_bet(
                action.amount,
                current_tree_node.available_bets,
                effective_stack,
                game_state,
                action.street
            )
            if mapped_bet_size == null
                log_error("Unable to map bet size: no suitable node or too large divergence.")
                return null
            end if
            
            next_tree_node = find_matching_node(current_tree_node, action.type, mapped_bet_size)
        else
            // For call/check/fold actions
            next_tree_node = find_matching_node(current_tree_node, action.type, null)
        end if
        
        if next_tree_node == null
            log_error("No matching child node found for the given action.")
            return null
        end if
        
        current_tree_node = next_tree_node
        
        // Verify that the current node can represent the next set of legal actions
        if not verify_legal_actions(current_tree_node, game_state, effective_stack)
            log_error("Legal actions verification failed at the current node.")
            return null
        end if
        
        // Check if the game state has diverged too much from the abstracted tree state
        real_metrics = extract_metrics(game_state)
        tree_metrics = extract_metrics_from_tree(current_tree_node)
        
        if not check_divergence_thresholds(real_metrics, tree_metrics, action.street)
            log_error("State divergence exceeds permissible thresholds.")
            return null
        end if
    end for
    
    return current_tree_node
end function


///////////////////////////////////////////
// Stack and Pot Calculations
///////////////////////////////////////////

function calculate_effective_stack(game_state)
    if game_state.is_multiway
        // Multi-way scenario: consider complex logic for side pots if necessary
        return calculate_multiway_effective_stack(game_state)
    end if
    
    // Heads-up or simple scenario: smallest remaining stack among active players
    return minimum(player.stack for player in game_state.active_players)
end function

function calculate_pot_size(game_state)
    if game_state.is_multiway
        // In multi-way pots, pot might be split into main pot + side pots
        return sum(game_state.main_pot + sum(game_state.side_pots))
    end if
    
    return sum(player.total_bet for player in game_state.players)
end function

function calculate_multiway_effective_stack(game_state)
    // For each pot, consider the minimum stack of players eligible for that pot
    stacks = []
    for pot in game_state.all_pots
        stacks.append(minimum(player.stack for player in pot.eligible_players))
    end for
    // The effective stack could be the largest among these minimums depending on definition,
    // but logic may vary based on the specific CFR abstraction needs.
    return maximum(stacks)
end function


///////////////////////////////////////////
// Bet Size Mapping
///////////////////////////////////////////

function map_bet_size_to_tree_bet(real_bet_size, tree_bet_sizes, effective_stack, game_state, street)
    // Check if the player is effectively going all-in
    if real_bet_size >= effective_stack
        // Find if there's an all-in node available
        mapped_allin = find_matching_node_for_allin(tree_bet_sizes, effective_stack)
        return mapped_allin  // Could be 'allin' or a node representing all-in
    end if
    
    mapped_bet_size = find_nearest_bet_size(real_bet_size, tree_bet_sizes)
    if mapped_bet_size == null
        return null
    end if
    
    // Calculate tolerance dynamically based on street and pot size
    current_pot_size = calculate_pot_size(game_state)
    tolerance = calculate_street_based_tolerance(street, current_pot_size)
    
    // Check if the chosen bet size is within the allowed tolerance
    if abs(mapped_bet_size - real_bet_size) <= tolerance
        return mapped_bet_size
    end if
    
    return null
end function

function find_nearest_bet_size(real_bet_size, tree_bet_sizes)
    numeric_bets = [bet for bet in tree_bet_sizes if is_numeric(bet)]
    if numeric_bets is empty
        return null
    end if
    
    closest_bet = null
    smallest_diff = infinity
    for bet in numeric_bets
        diff = abs(bet - real_bet_size)
        if diff < smallest_diff
            smallest_diff = diff
            closest_bet = bet
        end if
    end for
    
    return closest_bet
end function

function find_matching_node_for_allin(tree_bet_sizes, effective_stack)
    // If all-in is an action represented in tree_bet_sizes, return it
    if 'allin' in tree_bet_sizes
        return 'allin'
    end if
    // Otherwise, attempt to find a bet size that is effectively all-in (close to effective_stack)
    allin_like_bets = [bet for bet in tree_bet_sizes if is_numeric(bet) and bet >= effective_stack]
    if not empty(allin_like_bets)
        return minimum(allin_like_bets)  // The smallest bet >= effective_stack
    end if
    
    return null
end function


///////////////////////////////////////////
// Node Navigation
///////////////////////////////////////////

function find_matching_node(current_node, action_type, bet_size)
    for child in current_node.children
        if action_type == 'bet' or action_type == 'raise'
            // Handle both numeric and 'allin' cases
            if child.action == action_type
                if bet_size == 'allin' and child.is_allin
                    return child
                else if child.bet_size == bet_size
                    return child
                end if
            end if
        else
            // For actions like 'call', 'check', or 'fold'
            if child.action == action_type
                return child
            end if
        end if
    end for
    
    return null
end function


///////////////////////////////////////////
// Legal Actions Verification
///////////////////////////////////////////

function verify_legal_actions(tree_node, game_state, effective_stack)
    tree_actions = tree_node.available_actions
    game_actions = game_state.get_legal_actions()
    
    for action in game_actions
        if action.type in ['bet', 'raise']
            mapped_bet = map_bet_size_to_tree_bet(
                action.amount,
                tree_node.available_bets,
                effective_stack,
                game_state,
                action.street
            )
            if mapped_bet == null or not is_valid_tree_bet(mapped_bet, tree_node)
                return false
            end if
        else
            // For call/check/fold actions, just ensure the tree can handle them
            if action.type not in tree_actions
                return false
            end if
        end if
    end for
    
    return true
end function

function is_valid_tree_bet(mapped_bet, tree_node)
    // Check if the mapped bet size or 'allin' exists in the node's bets
    if mapped_bet == 'allin'
        return 'allin' in tree_node.available_bets
    end if
    
    return mapped_bet in tree_node.available_bets
end function


///////////////////////////////////////////
// Divergence Checks
///////////////////////////////////////////

function check_divergence_thresholds(real_metrics, tree_metrics, street)
    thresholds = get_street_thresholds(street)
    
    // Check SPR divergence
    if not check_spr_divergence(real_metrics, tree_metrics, thresholds.spr)
        return false
    end if
    
    // Check stack and pot divergence
    stack_div = calculate_stack_divergence(real_metrics, tree_metrics)
    pot_div = calculate_pot_divergence(real_metrics, tree_metrics)
    
    if stack_div > thresholds.stack
        return false
    end if
    if pot_div > thresholds.pot
        return false
    end if
    
    return true
end function

function get_street_thresholds(street)
    // Example thresholds per street (tweak as needed)
    return switch street
        case 'preflop': return {stack: 0.25, pot: 0.30, spr: 0.20}
        case 'flop':    return {stack: 0.20, pot: 0.25, spr: 0.15}
        case 'turn':    return {stack: 0.15, pot: 0.20, spr: 0.12}
        case 'river':   return {stack: 0.10, pot: 0.15, spr: 0.10}
        default:        return {stack: 0.25, pot: 0.30, spr: 0.20}  // default fallback
    end switch
end function

function check_spr_divergence(real_metrics, tree_metrics, spr_threshold)
    if real_metrics.pot_size == 0 and tree_metrics.pot_size == 0
        // Both infinite SPR - treat as no divergence
        return true
    end if
    
    if real_metrics.pot_size == 0 or tree_metrics.pot_size == 0
        // One has infinite SPR, the other doesn't
        return false
    end if
    
    // Normal SPR calculation
    if real_metrics.SPR == 0
        // If real SPR is 0 and tree SPR is not zero, big divergence
        return tree_metrics.SPR == 0
    end if
    
    spr_div = abs(real_metrics.SPR - tree_metrics.SPR) / real_metrics.SPR
    return spr_div <= spr_threshold
end function

function calculate_stack_divergence(real_metrics, tree_metrics)
    if real_metrics.initial_stack == 0
        // Cannot compute meaningful divergence if initial stack is 0
        return infinity
    end if
    return abs(real_metrics.effective_stack - tree_metrics.effective_stack) / real_metrics.initial_stack
end function

function calculate_pot_divergence(real_metrics, tree_metrics)
    if real_metrics.pot_size == 0
        // If real pot size is zero, any difference is huge
        return tree_metrics.pot_size == 0 ? 0 : infinity
    end if
    return abs(real_metrics.pot_size - tree_metrics.pot_size) / real_metrics.pot_size
end function


///////////////////////////////////////////
// Metrics Extraction
///////////////////////////////////////////

function extract_metrics(game_state)
    eff_stack = calculate_effective_stack(game_state)
    p_size = calculate_pot_size(game_state)
    init_stack = game_state.initial_stack
    
    if p_size > 0
        spr = eff_stack / p_size
    else
        spr = infinity
    end if
    
    return {
        effective_stack: eff_stack,
        pot_size: p_size,
        initial_stack: init_stack,
        SPR: spr
    }
end function

function extract_metrics_from_tree(tree_node)
    // Ensure tree node has these metrics stored or can compute them
    // For simplicity, assume they are stored attributes
    eff_stack = tree_node.effective_stack
    p_size = tree_node.pot_size
    init_stack = tree_node.initial_stack
    
    if p_size > 0
        spr = eff_stack / p_size
    else
        spr = infinity
    end if
    
    return {
        effective_stack: eff_stack,
        pot_size: p_size,
        initial_stack: init_stack,
        SPR: spr
    }
end function


///////////////////////////////////////////
// Street-Based Tolerance
///////////////////////////////////////////

function calculate_street_based_tolerance(street, pot_size)
    base_tolerance = switch street
        case 'preflop': return 0.15
        case 'flop':    return 0.10
        case 'turn':    return 0.08
        case 'river':   return 0.05
        default:        return 0.10
    end switch
    
    return base_tolerance * pot_size
end function


///////////////////////////////////////////
// Logging and Error Handling
///////////////////////////////////////////

function log_error(message)
    // Log or print the error for debugging
    print("ERROR: " + message)
end function

// Note: We decided to just return null on errors, 
// no need for a separate handle_unmappable_actions() 
// since we log_error and then return null.


///////////////////////////////////////////
// Utility Functions
///////////////////////////////////////////

function is_numeric(value)
    // Check if value is a number
    return typeof(value) == 'number'
end function

function is_allin(bet_size)
    return bet_size == 'allin'
end function

function calculate_spr(game_state)
    p_size = calculate_pot_size(game_state)
    eff_stack = calculate_effective_stack(game_state)
    if p_size > 0
        return eff_stack / p_size
    else
        return infinity
    end if
end function

function has_required_metrics(tree_node)
    // Check if tree_node contains the needed attributes
    return (tree_node.effective_stack != null and 
            tree_node.pot_size != null and
            tree_node.initial_stack != null)
end function

function compute_tree_node_metrics(tree_node)
    // If needed, compute these metrics from the node's state or its ancestors
    // This is a placeholder. Actual logic depends on how the tree is built.
    return {
        effective_stack: tree_node.effective_stack,
        pot_size: tree_node.pot_size,
        initial_stack: tree_node.initial_stack,
        SPR: (tree_node.pot_size > 0) ? (tree_node.effective_stack / tree_node.pot_size) : infinity
    }
end function



To make real implementation here are key pieces:

```
// package chips
const epsilon = 1e-8

type Chips float32

var Zero = Chips(0)

// NewFromInt creates a new Chips value from an int.
func NewFromInt(val int64) Chips {
	return Chips(float32(val))
}

// NewFromFloat64 creates a new Chips value from a float64.
func NewFromFloat64(val float64) Chips {
	return Chips(float32(val))
}

// NewFromFloat creates a new Chips value from a float64.
func NewFromFloat(val float64) Chips {
	return Chips(float32(val))
}

// NewFromFloat32 creates a new Chips value from a float32.
func NewFromFloat32(val float32) Chips {
	return Chips(val)
}

// NewFromString creates a new Chips value from a string.
// Assumes the string is a valid float number.
func NewFromString(val string) Chips {
	var f float32
	fmt.Sscanf(val, "%f", &f)
	return Chips(f)
}

// String returns a string representation of c.
func (c Chips) String() string {
	return fmt.Sprintf("%f", c)
}

// Add adds a and b and returns the result.
func (c Chips) Add(b Chips) Chips {
	return c + b
}

// Sub subtracts b from a and returns the result.
func (c Chips) Sub(b Chips) Chips {
	return c - b
}

// Mul multiplies a and b and returns the result.
func (c Chips) Mul(b Chips) Chips {
	return c * b
}

// Div divides a by b and returns the result.
func (c Chips) Div(b Chips) Chips {
	if b == 0 {
		return Chips(math.Inf(1)) // return +Inf on division by zero
	}
	return c / b
}

// Equal checks if c and b are equal within a tolerance level.
func (c Chips) Equal(b Chips) bool {
	return math.Abs(float64(c-b)) < epsilon
}

// GreaterThan checks if c is greater than b.
func (c Chips) GreaterThan(b Chips) bool {
	return c > b
}

// GreaterThanOrEqual checks if c is greater than or equal b.
func (c Chips) GreaterThanOrEqual(b Chips) bool {
	return c >= b
}

// LessThan checks if c is less than b.
func (c Chips) LessThan(b Chips) bool {
	return c < b
}

// LessThanOrEqual checks if c is less than b.
func (c Chips) LessThanOrEqual(b Chips) bool {
	return c <= b
}

func (c Chips) StringFixed(places int) string {
	format := fmt.Sprintf("%%.%df", places)
	return fmt.Sprintf(format, c)
}

// Abs returns the absolute value of chips.
func (c Chips) Abs() Chips {
	return Chips(math.Abs(float64(c)))
}

// Pow returns the power of i.
func (c Chips) Pow(i Chips) Chips {
	return Chips(math.Pow(float64(c), float64(i)))
}

// Float32 gives float32 representation
func (c Chips) Float32() float32 {
	return float32(c)
}

// Float64 gives float64 representation
func (c Chips) Float64() float64 {
	return float64(c)
}

```

```go
// package tree

type NodeKind uint8
const (
    NodeKindRoot
    NodeKindChance 
    NodeKindTerminal
    NodeKindPlayer
)

type Node interface {
    GetID() uint32
    Kind() NodeKind
    GetParent() Node
    GetLastAction() (LastActionInfo, bool)
    Path() (b Runes)
}

type Root struct {
    ID uint32
    Nodes uint32
    Next Node
    NumPlayers uint8
    EffectiveStack float32
    BettingOptions []float32
    Terminal table.Street
    MaxActions uint8
}

type Player struct {
	ID         uint32
	Parent     Node
	PlayerID   uint8
	LastAction LastAction
	Actions    Actions
	Street     table.Street
}

type Actions struct {
    Parent Node
    Actions []table.DiscreteAction // pot multipliers
    Nodes []Node
}

type Terminal struct {
    ID uint32
    Parent Node
    LastAction LastAction
    Pot float32
    Paid PlayerAmount
    Status PlayerStatus
    Street table.Street
}

type Rune string
const (
    RuneRoot     Rune = "r"
    RuneChance   Rune = "n"
    RuneTerminal Rune = "t"
    RunePlayer   Rune = "p"
    RuneBet      Rune = "b"
    RuneFold     Rune = "f"
    RuneCall     Rune = "c"
    RuneCheck    Rune = "k"
    RuneAllIn    Rune = "a"
)

type Chance struct {
	ID         uint32
	Next       Node
	Parent     Node
	LastAction LastAction
	Street     table.Street
}

type Root struct {
	ID             uint32
	Nodes          uint32
	Next           Node
	NumPlayers     uint8
	EffectiveStack float32
	BettingOptions []float32
	Terminal       table.Street
	MaxActions     uint8
}


type Player struct {
	ID         uint32
	Parent     Node
	PlayerID   uint8
	LastAction LastAction
	Actions    Actions
	Street     table.Street
}

type Terminal struct {
	ID         uint32
	Parent     Node
	LastAction LastAction
	Pot        float32
	Paid       PlayerAmount
	Status     PlayerStatus
	Street     table.Street
}

```


```go
// package table

type DiscreteAction float32

const (
    DUnknown DiscreteAction = -4
    DAllIn   DiscreteAction = -3
    DFold    DiscreteAction = -2
    DCall    DiscreteAction = -1
    DCheck   DiscreteAction = 0
    // Positive values represent pot multipliers
)

type ActionKind uint8

// Real game actions
const (
    Bet ActionKind = iota
    SmallBlind
    BigBlind
    Fold
    Check
    Call
    Raise
    AllIn
)

type ActionAmount struct {
	Action ActionKind
	Amount chips.Chips
}


// Record represents a single poker action with complete state information.
// Used for tracking game history and analyzing player actions.
type Record struct {
    ID               int          // Unique identifier for the action
    PlayerID         uint8        // ID of player who took the action
    Street           Street       // Street when action occurred (Preflop/Flop/Turn/River)
    Action           ActionKind   // Type of action taken (Bet/Call/Raise etc)
    Amount           chips.Chips  // Amount of chips involved in action
    StackBefore      chips.Chips  // Player's stack before action
    StackAfter       chips.Chips  // Player's stack after action
    Addition         chips.Chips  // Additional amount above previous bet
    CommitedOnStreet chips.Chips  // Total amount player has committed on this street
}

// Records is a slice of Record objects representing a sequence of actions
type Records []Record

// Clone creates a deep copy of Records.
// Used when game state needs to be duplicated.
func (r Records) Clone() Records

// StreetRecord maps each street to its sequence of actions.
// Provides organized access to betting history by street.
type StreetRecord map[Street]Records

// Clone creates a deep copy of StreetRecord.
// Useful for game state snapshots and simulations.
func (s StreetRecord) Clone() StreetRecord

// String returns a human-readable representation of the StreetRecord.
// Format: "Street: {street}\nAction {n}: {action} ${amount} (cm ${committed})\n"
func (s StreetRecord) String() string

// StreetPaid tracks how much each player has paid on each street.
// Used for pot calculations and bet sizing.
type StreetPaid map[Street][]chips.Chips

// Clone creates a deep copy of StreetPaid.
// Important for maintaining independent game states.
func (s StreetPaid) Clone() StreetPaid

// PlayerAction pairs a player with their action for history tracking.
type PlayerAction struct {
    Player *Player
    Action ActionAmount
}

// StreetPlayerAction maps each street to its sequence of player actions.
// Primary history structure for game replay and analysis.
type StreetPlayerAction map[Street][]PlayerAction

// Clone creates a deep copy of StreetPlayerAction.
// Essential for game state management and simulation.
func (s StreetPlayerAction) Clone() StreetPlayerAction

// String returns a condensed representation of the action sequence.
// Format: "f:c:r/c:k/b:c/" (fold, call, raise/call, check/bet, call)
// Used for hand history analysis and debugging.
func (sx StreetPlayerAction) String() string

type PlayerAction struct {
    Player *Player
    Action ActionAmount
}

type StreetPlayerAction map[Street][]PlayerAction

// Player Types
type Status uint8

type Player struct {
    ID           uint8
    HoleCard     card.Cards
    InitialStack chips.Chips
    Stack        chips.Chips
    Paid         chips.Chips
    Rank         card.HandRank
    Status       Status
    PaidOnStreet StreetPaid
    History      StreetRecord
}

type Players []*Player

type Filter func(p *Player) bool

type State struct {
	Players            Players
	SbAmount           chips.Chips
	CommunityCard      card.Cards
	ActionCount        int
	Street             Street
	MaxActionsPerRound int
	BtnPos             uint8
	NextPlayerPos      uint8
	SbPos              uint8
	BbPos              uint8
	Winners            Players
	History            StreetPlayerAction
	PrizeMap           PrizeMap
	BetSizes           []float32
	TerminalStreet     Street
}

// LegalActions represents legal actions for a player.
// It is a map of ActionKind to chips.Chips.
// Chip amount is minimum of chips needed to perform action.
// Maximum in NLH is player stack.
type LegalActions map[ActionKind]chips.Chips        // Real game actions

// DiscreteLegalActions represents legal discrete actions for a player.
// It is similar in nature to LegalActions except chips are excat amounts
// to be used for action.
type DiscreteLegalActions map[DiscreteAction]chips.Chips // Tree actions

// Pot Types
type Pot struct {
    Amount  chips.Chips
    Players Players
}
type Pots []Pot



// Action Functions
func NewLegalActions(r *State) LegalActions
func NewDiscreteLegalActions(r *State) DiscreteLegalActions
func (l LegalActions) Validate(p *Player, action ActionKind, amount chips.Chips) error
func FromActionToDiscrete(r *State, action ActionKind, amount chips.Chips) DiscreteAction

// State Functions
func NewState(pp Players, sba chips.Chips, btnPos int) (*State, error)
func (r *State) Clone() *State
func (r *State) NextPlayer() *Player
func (r *State) CallAmount(player *Player) chips.Chips
func (r *State) Finished() bool

// Round Functions
func NewRound(s *State, d *card.Deck) (*Round, error)
func (r *Round) Action(aa Actioner) error
func (r *Round) State() *State

// Command Functions
// MakeAction executes a poker action for the next player in the game state.
// It validates the action, updates player status, handles all-in situations,
// and records the action in game history.
// Parameters:
//   - r: Current game state
//   - aa: Action to be performed (implements Actioner interface)
// Returns error if action is invalid or cannot be performed
func MakeAction(r *State, aa Actioner) error

// InitRound initializes a new round of poker by dealing hole cards
// and making initial blind bets.
// Parameters:
//   - r: Game state to initialize
//   - deck: Deck to deal cards from
// Returns error if initial bets cannot be made
func InitRound(r *State, deck *card.Deck) error

// MakeInitialBets posts the small and big blinds at the start of a hand.
// Automatically shifts turns between blind postings.
// Returns error if blind bets cannot be made
func MakeInitialBets(r *State) error

// RuleKind represents different game state transitions
type RuleKind uint8

const (
    RuleShiftTurn RuleKind = iota           // Move to next player
    RuleShiftStreet                         // Move to next street
    RuleShiftStreetUntilEnd                 // Process remaining streets
    RuleFinish                              // End the hand
)


// Rule determines the next game state transition based on current state.
// Returns RuleKind indicating what action should be taken:
//   - RuleShiftTurn: Move to next player
//   - RuleShiftStreet: Move to next street
//   - RuleShiftStreetUntilEnd: Move through streets until end
//   - RuleFinish: End the hand

func Rule(r *State) RuleKind

// Move advances the game state based on the current Rule.
// Handles turn progression, street changes, and game completion.
// Parameters:
//   - r: Current game state
//   - d: Deck for dealing community cards
// Returns error if state transition fails
func Move(r *State, d *card.Deck) error

// ShiftStreet advances the game to the next street (Flop/Turn/River).
// Deals appropriate community cards and resets betting positions.
// Parameters:
//   - r: Current game state
//   - deck: Deck to deal community cards from
// Returns error if street transition fails
func ShiftStreet(r *State, deck *card.Deck) error

// ShiftTurn moves the action to the next eligible player.
// Returns error if no valid next player is found
func ShiftTurn(r *State) error

// ShiftTurnStreetStart sets the first player to act on a new street.
// Handles both heads-up and multiway pot scenarios.
// Different logic for preflop vs postflop positions.
func ShiftTurnStreetStart(r *State)

// DealHoleCard deals two hole cards to each player in the game.
// Parameters:
//   - r: Current game state
//   - deck: Deck to deal from
func DealHoleCard(r *State, deck *card.Deck)

// DealCommunityCards adds specified cards to the board and updates
// player hand rankings.
// Parameters:
//   - r: Current game state
//   - cards: Cards to add to community cards
func DealCommunityCards(r *State, cards card.Cards)

// Finish ends the current hand, determines winners, distributes prizes,
// and updates player stacks.
// Returns error if prize distribution fails
func Finish(r *State) error

// Player Functions
func NewPlayer(id uint8, stack chips.Chips) *Player
func (p *Player) Clone() *Player
func (p *Player) PotOnStreet(street Street) chips.Chips
func (p *Player) Pay(street Street, amount chips.Chips)
func (p *Player) Add(amount chips.Chips)
func (p *Player) SetStatus(status Status)
func (p *Player) AddRecord(record Record) int
func (p *Player) UpdateRank(community card.Cards) error
func (p *Player) IsActive() bool
func (p *Player) IsWaitingAsk() bool
func (p *Player) PaidSum() chips.Chips
func (p *Player) HaveToAct(street Street, maxStreetBet chips.Chips) bool

// Players Functions
func NewPlayers(players ...*Player) Players
func (p Players) Clone() Players
func (p Players) Filter(a Filter) Players
func (p Players) GetIndexByID(id uint8) int
func (p Players) GetByID(id uint8) (*Player, bool)
func (p Players) FromPosition(startPos int) Players
func (p Players) MaxStackExcept(px *Player) chips.Chips
func (p Players) BiggestPotOnStreet(street Street) chips.Chips
func (p Players) BiggestRaiseOnStreet(street Street) (Record, bool)
func (p Players) HaveToAct(street Street) Players
func (p Players) AllIn() Players
func (p Players) Paid() chips.Chips
func (p Players) PaidMax() chips.Chips
func (p Players) Pots() Pots
func (p Players) PrizeDist() map[uint8]chips.Chips
func (p Players) Winners() Players


// EffectiveStack returns the smallest stack among active players, representing
// the maximum amount that can be won from any single player in the current hand.
// For heads-up play this is simply the smaller of the two stacks.
// For multiway pots it considers side pot eligibility.
func (p Players) EffectiveStack() chips.Chips {
    stack := p[0].Stack
    for _, player := range p[1:] {
        if player.Stack.LessThan(stack) {
            stack = player.Stack
        }
    }
    return stack
}

// Paid returns the total amount of chips committed to the pot by all players.
// This represents the current pot size.
func (p Players) Paid() chips.Chips {
    sum := chips.Zero
    for _, player := range p {
        sum = sum.Add(player.Paid)
    }
    return sum
}

// PaidMax returns the largest amount paid by any single player.
// Used for calculating pot-sized bets and verifying bet sizing.
func (p Players) PaidMax() chips.Chips {
    maxPaid := p[0].Paid
    for _, player := range p[1:] {
        if player.Paid.GreaterThan(maxPaid) {
            maxPaid = player.Paid
        }
    }
    return maxPaid
}

// Pots returns all pots (main pot and side pots) based on all-in players
// and their committed amounts. Used for handling all-in situations and
// calculating prize distributions.
// Returns slice of Pot structs containing amount and eligible players.
func (p Players) Pots() Pots {
    // Implementation details...
}

// BiggestPotOnStreet returns the largest amount committed by any player
// on the given street. Used for calculating legal bet sizes and raises.
func (p Players) BiggestPotOnStreet(street Street) chips.Chips {
    bets := make([]chips.Chips, len(p))
    for i, player := range p {
        bets[i] = player.PotOnStreet(street)
    }
    max := bets[0]
    for _, value := range bets[1:] {
        if value.GreaterThan(max) {
            max = value
        }
    }
    return max
}

// BiggestRaiseOnStreet returns the largest raise action and amount made
// on the given street. Used for enforcing minimum raise sizes.
// Returns the record of the raise action and whether one was found.
func (p Players) BiggestRaiseOnStreet(street Street) (Record, bool)

// InitStackSum returns the sum of all players' initial stacks.
// Used for tracking total chips in play and verifying game integrity.
func (p Players) InitStackSum() chips.Chips {
    sum := chips.Zero
    for _, player := range p {
        sum = sum.Add(player.InitialStack)
    }
    return sum
}

// CallAmount calculates how much the given player needs to call to match
// the current bet on the specified street.
// Returns 0 if player has already matched or exceeded the current bet.
func (p Players) CallAmount(street Street, player *Player) chips.Chips {
    largestBet := p.BiggestPotOnStreet(street)
    currentBet := player.PotOnStreet(street)
    callAmount := largestBet.Sub(currentBet)
    if callAmount.LessThan(chips.Zero) {
        return chips.Zero
    }
    return callAmount
}

// MaxStackExcept returns the largest stack size among all players except
// the specified player. Used for calculating maximum possible win amounts
// and effective stacks in multiway scenarios.
func (p Players) MaxStackExcept(px *Player) chips.Chips {
    stack := chips.Zero
    for _, player := range p {
        if px == player {
            continue
        }
        if player.Stack.GreaterThan(stack) {
            stack = player.Stack
        }
    }
    return stack
}


// Testing
p1 := NewPlayer(0, chips.NewFromInt(100))
p2 := NewPlayer(1, chips.NewFromInt(100))

pp := NewPlayers(p1, p2)
s, err := NewState(pp, chips.NewFromInt(1), 0)
require.NoError(t, err)

d := card.NewDeck(card.All())

DealHoleCard(s, d)
err = MakeInitialBets(s)
require.NoError(t, err)

err = MakeAction(s, DCall)
require.NoError(t, err)
err = Move(s, d)
require.NoError(t, err)

err = MakeAction(s, DCheck)
require.NoError(t, err)
err = Move(s, d)
require.NoError(t, err)

.....

require.Equal(t, s.Finished(), true)
```


Signature of new function:

```go

type MatchResult struct {
    Player *tree.Player
    // Somhow legal actions that are in tree and match
    // table state.
}
func Match(node tree.Node, s *table.State) (p *tree.Player, err error)
```

Basically given tree.Node (root) and state, which is state at any given moment for the game, we want to find a player that should play in tree.

Tree node id is associated with strategy, so we are looking to find ID
of the node to obtain strategy.

WRITE GOLANG CODE. 
WRITE GOLANG CODE. 
WRITE GOLANG CODE.


ALL FUNCTIONS MUST BE REUSABLE AND EXPORTED. IT ESSENTIAL TO USE EXISTING PACKAGES OUTLINED ABOVE.

Use packages table and tree. BE 100% sure to follow the pseudo code but only as guide not as real implementation.