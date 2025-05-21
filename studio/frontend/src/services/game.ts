import { call } from "./rpc";

export async function newGame(): Promise<number> {
  return await call("game_new_mc", {
    Players: 2,
    Stack: 100,
    SB: 1,
  });
}

export type GameParams = {
  num_players: number;
  max_actions_per_round: number;
  btn_pos: number;
  sb_amount: number;
  bet_sizes: number[];
  initial_stacks: number[];
  terminal_street: number;
  min_bet: boolean;
  disable_v?: boolean;
};

export enum PlayerStatus {
  UNKNOWN = 0,
  FOLDED = 1,
  ACTIVE = 2,
  ALL_IN = 3,
}

export type Player = {
  paid: number;
  status: PlayerStatus;
};

export enum ActionKind {
  NO_ACTION = 0,
  BET = 1,
  SMALL_BLIND = 2,
  BIG_BLIND = 3,
  FOLD = 4,
  CHECK = 5,
  CALL = 6,
  RAISE = 7,
  ALL_IN = 8,
}

export type BSC = {
  amount: number;
  addition: number;
  action: ActionKind;
};

export enum Street {
  NO_STREET = 0,
  PREFLOP = 1,
  FLOP = 2,
  TURN = 3,
  RIVER = 4,
  FINISHED = 5,
}

export type State = {
  players: Player[];
  street: Street;
  turn_pos: number;
  btn_pos: number;
  street_action: number;
  call_amount: number;
  bsc: BSC;
  psc: number[];
  psac: number[];
  psla: ActionKind[];
};

export enum DiscreteAction {
  ALL_IN = -4,
  FOLD = -3,
  CALL = -2,
  CHECK = -1,
  NO_ACTION = 0,
  // Values greater than 0 represent raise/bet sizes as pot multiples
}

export type Action = number;

// Format from Go: Array of [discreteAction, amount] pairs from DiscreteLegalActions
export type LegalActions = Array<[Action, number]>;

export type Game = {
  params: GameParams;
  state: State;
  hole: string[];
  community: string[] | null;
  legal: LegalActions;
  pot: number;
  round: number;
  winner: number[] | null;
};

export async function getGameState(gameId: number): Promise<Game> {
  return await call("game_get_state", gameId);
}

export async function action(gameId: number, action: Action): Promise<void> {
  return await call("game_action", gameId, action);
}

const w: any = window;
w.newGame = newGame;
w.getGameState = getGameState;
w.action = action;
