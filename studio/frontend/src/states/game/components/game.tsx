import { Box, Grid } from "@mui/material";
import Table from "./table";
import Controls from "./controls";
import PlayerInfo from "./player";
import StartScreen from "./start";
import { useGameStore } from "./store";
import { Action, Game, Street, LegalActions } from "~/services/game";

// Game wrapper component - displays either start screen or active game
export default function GamePage() {
  const { gameId, gameState, loading, error, startGame, makeAction } =
    useGameStore();

  // We only consider the game active if we have both gameId and gameState
  const gameActive = gameId !== null && gameState !== null;

  if (!gameActive) {
    return (
      <GameStartContainer loading={loading} error={error} onStart={startGame} />
    );
  }

  return <ActiveGame gameState={gameState} onAction={makeAction} />;
}

// Type for GameStartContainer props
interface GameStartContainerProps {
  loading: boolean;
  error: string | null;
  onStart: () => void;
}

// Container for the start screen
function GameStartContainer({
  loading,
  error,
  onStart,
}: GameStartContainerProps) {
  return (
    <Box sx={{ p: 1, height: "100vh", bgcolor: "#f5f5f5" }}>
      <StartScreen loading={loading} error={error} onStart={onStart} />
    </Box>
  );
}

// Type for ActiveGame props
interface ActiveGameProps {
  gameState: Game;
  onAction: (action: Action) => void;
}

// Main game container - only shown when game is active
function ActiveGame({ gameState, onAction }: ActiveGameProps) {
  // Process game state to extract all necessary data
  const processedGameData = processGameState(gameState, onAction);

  return (
    <Box sx={{ p: 1, height: "100vh", bgcolor: "#f5f5f5" }}>
      <Box
        sx={{
          p: 2,
          height: "calc(100% - 16px)",
          borderRadius: 2,
          display: "flex",
          flexDirection: "column",
        }}
      >
        <GameLayout
          tableProps={processedGameData.tableProps}
          controlsProps={processedGameData.controlsProps}
          playerInfoProps={processedGameData.playerInfoProps}
        />
      </Box>
    </Box>
  );
}

// Type for GameLayout props
interface GameLayoutProps {
  tableProps: TableProps;
  controlsProps: ControlsProps;
  playerInfoProps: PlayerInfoProps;
}

// Component to handle the layout of the game UI
function GameLayout({
  tableProps,
  controlsProps,
  playerInfoProps,
}: GameLayoutProps) {
  return (
    <Grid container spacing={1} sx={{ flexGrow: 1 }}>
      <Grid item xs={12} md={10}>
        <Box sx={{ height: "100%", display: "flex", flexDirection: "column" }}>
          <Box
            sx={{
              flexGrow: 1,
              position: "relative",
              minHeight: "400px",
              height: { xs: "60vh", md: "65vh" },
              maxHeight: "600px",
            }}
          >
            <Table {...tableProps} />
          </Box>

          <Box sx={{ mt: 0.5 }}>
            <Controls {...controlsProps} />
          </Box>
        </Box>
      </Grid>

      <Grid item xs={12} md={2}>
        <PlayerInfo {...playerInfoProps} />
      </Grid>
    </Grid>
  );
}

// Extract all props used by the Table component
interface TableProps {
  pot: number;
  community: string[];
  playerCards: string[];
  playerTurn: boolean;
  street: Street;
}

// Extract all props used by the Controls component
interface ControlsProps {
  isPlayerTurn: boolean;
  legal: LegalActions;
  bigBlind: number;
  onAction: (action: Action) => void;
}

// Extract all props used by the PlayerInfo component
interface PlayerInfoProps {
  playerStack: number;
  opponentStack: number;
  playerPos: number;
  btnPos: number;
  round: number;
  lastWinners: number[];
  street: Street;
}

// Process the game state to get all required props for child components
function processGameState(
  gameState: Game,
  onAction: (action: Action) => void
): {
  tableProps: TableProps;
  controlsProps: ControlsProps;
  playerInfoProps: PlayerInfoProps;
} {
  // Extract data from gameState - we know it exists because we've already checked
  const isPlayerTurn = gameState.state?.turn_pos === 0;
  const street = gameState.state?.street ?? Street.NO_STREET;

  // Get stacks - with minimal fallback for missing properties
  const initialStacks = gameState.params?.initial_stacks || [100, 100];
  const players = gameState.state?.players || [{ paid: 0 }, { paid: 0 }];

  // Calculate derived values
  const playerStack = Math.max(0, initialStacks[0] - players[0].paid);
  const opponentStack = Math.max(0, initialStacks[1] - players[1].paid);
  const bigBlind = Math.max(1, (gameState.params?.sb_amount || 1) * 2);

  // Get pot amount
  const pot = gameState.pot ?? players[0].paid + players[1].paid;

  // Get cards
  const communityCards = gameState.community ?? [];
  const playerCards = gameState.hole ?? [];

  return {
    tableProps: {
      pot,
      community: communityCards,
      playerCards,
      playerTurn: isPlayerTurn,
      street,
    },
    controlsProps: {
      isPlayerTurn,
      legal: gameState.legal || [],
      bigBlind,
      onAction,
    },
    playerInfoProps: {
      playerStack,
      opponentStack,
      playerPos: 0,
      btnPos: gameState.state?.btn_pos || 0,
      round: gameState.round || 0,
      lastWinners: gameState.winner ?? [],
      street,
    },
  };
}
