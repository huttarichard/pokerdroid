import { call } from "./rpc";

export type NodeKind = "root" | "player" | "chance" | "terminal";

export interface Action {
  kind: NodeKind;
  actionIdx?: number; // For player nodes
  cards?: string[]; // For chance nodes
  solution?: number;
}

// Discrete actions mapping based on table/action.go
export enum DiscreteAction {
  AllIn = -4,
  Fold = -3,
  Call = -2,
  Check = -1,
  NoAction = 0,
  // > 0 = Raise/Bet with pot multiplier
}

export function discreteActionToString(action: number): string {
  switch (action) {
    case DiscreteAction.AllIn:
      return "All In";
    case DiscreteAction.Fold:
      return "Fold";
    case DiscreteAction.Call:
      return "Call";
    case DiscreteAction.Check:
      return "Check";
    case DiscreteAction.NoAction:
      return "No Action";
    default:
      return action > 0 ? `Bet ${action}x` : "Unknown";
  }
}

export interface Policy {
  iteration: number;
  strategy: number[];
  regret_sum: number[];
  strategy_sum: number[];
  baseline: number[];
}

export interface MatrixCell {
  cards: string[][];
  policies: Policy[];
  policy: Policy;
  avg_strat: number[];
  reach_probs: number[];
  reach: number;
  clusters: number[];
}

// Add a new type for Range based on equity.Range in inspect.go
export interface Range {
  win: number;
  lose: number;
  tie: number;
}

export interface TreeResponse {
  kind: NodeKind;
  street?: string;
  actions?: number[];
  matrix?: MatrixCell[][];
  state?: any; // New field for table state from inspect.go Result.State
  pot?: number; // New field for chips.Chips from inspect.go Result.Pot
  tree_history?: string;
}

export async function getTreeState(actions: Action[]): Promise<TreeResponse> {
  return await call("tree_get_state", actions);
}

export interface SolutionInfo {
  states: number;
  nodes: number;
  params: {
    num_players: number;
    max_actions_per_round: number;
    btn_pos: number;
    sb_amount: number;
    bet_sizes: number[];
    initial_stacks: number[];
    terminal_street: number;
  };
  iteration: number;
}

export async function getSolutions(): Promise<SolutionInfo[]> {
  return await call("tree_solutions");
}
