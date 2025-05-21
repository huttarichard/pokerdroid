# Action Translation Example with Pseudo-Harmonic Mapping

Let me demonstrate with a concrete poker scenario:

## Game Setup
- No-limit Texas Hold'em with $100 stacks
- P1 has A♥K♠, P2 has J♦J♣
- Current pot: $20
- Board: Q♠7♥2♦

## P1's Abstraction
P1's AI only considers these bet sizes in its abstraction:
- $0 (check)
- $20 (pot-sized bet)
- $60 (3x pot)
- $100 (all-in)

## The Action Translation Problem
1. P1 checks
2. P2 bets $35 (not in P1's abstraction)
3. P1 needs to translate this non-modeled bet to decide how to respond

## Solution Using Pseudo-Harmonic Mapping

### Step 1: Identify boundaries
- Lower boundary (A): $20
- Upper boundary (B): $60
- Opponent action (x): $35

### Step 2: Calculate mapping probability
Using the formula: fA,B(x) = (B-x)(1+A)/((B-A)(1+x))

f($20,$60)($35) = (60-35)(1+20)/((60-20)(1+35))
= 25 × 21 / (40 × 36)
= 525 / 1440
≈ 0.365

### Step 3: Apply the mapping
- With 36.5% probability, treat the $35 bet as if it were $20
- With 63.5% probability, treat it as if it were $60

### Step 4: Choose action based on strategy
If using randomized mapping:
- P1 randomly selects which abstraction point to use based on the probabilities
- If $20 is selected, P1 might call easily (only $20 to call)
- If $60 is selected, P1 might fold more often (would need to call $60)

If using deterministic mapping:
- Calculate threshold x* = (A+B+2AB)/(A+B+2) ≈ (20+60+22060)/(20+60+2) = $30.24
- Since $35 > $30.24, always map to $60

This approach maintains consistent strategy and is much less exploitable than simply mapping to the closest value ($20), which would cause P1 to call too often with marginal hands against a significant bet.

