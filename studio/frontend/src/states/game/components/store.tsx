import { create } from "zustand";
import {
  newGame,
  getGameState,
  action,
  Action,
  Game,
  DiscreteAction,
} from "~/services/game";

// Define the game manager state structure
interface GameStoreState {
  // State
  gameId: number | null;
  gameState: Game | null;
  loading: boolean;
  actionInProgress: Action | null; // Track which action is in progress
  error: string | null;
  // Actions
  startGame: () => Promise<void>;
  refreshState: (id?: number) => Promise<void>;
  makeAction: (actionType: Action) => Promise<void>;
  resetError: () => void;
}

/**
 * Global game store powered by Zustand
 * This store can be imported and used directly in any component
 * without needing to pass props down the component tree
 *
 * Example:
 * ```
 * const { gameState, makeAction } = useGameStore();
 * ```
 */
export const useGameStore = create<GameStoreState>((set, get) => ({
  // Initial state
  gameId: null,
  gameState: null,
  loading: false,
  actionInProgress: null,
  error: null,

  // Reset error state
  resetError: () => set({ error: null }),

  // Start a new game
  startGame: async () => {
    try {
      set({ loading: true, error: null });
      console.log("Starting new game...");
      const id = await newGame();
      console.log("Game created with ID:", id);
      set({ gameId: id });
      await get().refreshState(id);
    } catch (err) {
      console.error("Error starting game:", err);
      set({
        error: err instanceof Error ? err.message : "Failed to start game",
      });
    } finally {
      set({ loading: false });
    }
  },

  // Refresh the game state
  refreshState: async (id?: number) => {
    try {
      const gameId = id ?? get().gameId;
      if (gameId === null) return;

      // Don't set loading during auto-refresh to prevent UI flicker
      const isManualRefresh = id !== undefined;
      if (isManualRefresh) {
        set({ loading: true });
      }

      console.log("Fetching game state for ID:", gameId);
      const gameState = await getGameState(gameId);

      // Debug log the response
      console.log(
        "Game state raw response:",
        JSON.stringify(gameState, null, 2)
      );

      if (!gameState) {
        console.error("Empty game state received");
        return;
      }

      // Update state with the raw response - no defaults needed
      set({
        gameState,
        loading: false,
        // Clear action in progress when we get new state
        actionInProgress: null,
      });
    } catch (err) {
      console.error("Error refreshing game state:", err);
      set({
        error: err instanceof Error ? err.message : "Failed to get game state",
        loading: false,
        actionInProgress: null,
      });
    }
  },

  // Make an action in the game
  makeAction: async (actionType: Action) => {
    try {
      set({
        loading: true,
        error: null,
        actionInProgress: actionType,
      });

      const gameId = get().gameId;
      if (gameId === null) {
        throw new Error("No active game");
      }

      // Get action name for logging
      let actionName = "Unknown";
      switch (actionType) {
        case DiscreteAction.FOLD:
          actionName = "Fold";
          break;
        case DiscreteAction.CHECK:
          actionName = "Check";
          break;
        case DiscreteAction.CALL:
          actionName = "Call";
          break;
        case DiscreteAction.ALL_IN:
          actionName = "All-in";
          break;
        default:
          if (actionType > 0) {
            actionName = `Raise/Bet ${actionType}`;
          }
      }

      console.log(`Making action: ${actionName} (${actionType})`);
      await action(gameId, actionType);
      await get().refreshState();
    } catch (err) {
      console.error("Error making action:", err);
      set({
        error: err instanceof Error ? err.message : "Failed to make action",
        loading: false,
        actionInProgress: null,
      });
    }
  },
}));
