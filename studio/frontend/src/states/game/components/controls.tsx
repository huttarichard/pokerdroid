import { useState } from "react";
import {
  Box,
  Button,
  Slider,
  Typography,
  Paper,
  CircularProgress,
  Stack,
} from "@mui/material";
import { Action, DiscreteAction, LegalActions } from "~/services/game";
import { useGameStore } from "./store";

interface ControlsProps {
  isPlayerTurn: boolean;
  legal: LegalActions;
  bigBlind: number;
  onAction: (action: Action) => void;
}

export default function Controls({
  isPlayerTurn,
  legal,
  bigBlind,
  onAction,
}: ControlsProps) {
  const [betSize, setBetSize] = useState<number>(1);
  const { actionInProgress } = useGameStore();

  // When no legal actions are available, show loading state
  if (legal.length === 0) {
    return (
      <Paper sx={{ p: 1, borderRadius: 2 }}>
        <Box
          sx={{
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            gap: 1,
          }}
        >
          <CircularProgress size={16} color="primary" />
          <Typography variant="caption">Loading actions</Typography>
        </Box>
      </Paper>
    );
  }

  // When it's not the player's turn, show opponent's turn message
  if (!isPlayerTurn) {
    return (
      <Paper sx={{ p: 1, borderRadius: 2 }}>
        <Typography variant="subtitle2" align="center">
          Opponent's Turn
        </Typography>
      </Paper>
    );
  }

  // Find specific actions in the legal actions array
  const findAction = (action: DiscreteAction): [Action, number] | undefined => {
    return legal.find(([act]) => act === action);
  };

  // Determine available actions
  const checkAction = findAction(DiscreteAction.CHECK);
  const callAction = findAction(DiscreteAction.CALL);
  const foldAction = findAction(DiscreteAction.FOLD);
  const allInAction = findAction(DiscreteAction.ALL_IN);

  const hasCheck = Boolean(checkAction);
  const hasCall = Boolean(callAction);
  const hasFold = Boolean(foldAction);
  const hasAllIn = Boolean(allInAction);

  // Get amounts
  const callAmount = hasCall ? callAction![1] : 0;
  const allInAmount = hasAllIn ? allInAction![1] : 0;

  // Find raise actions (positive numbers)
  const raiseActions = legal.filter(([act]) => act > 0);
  const hasRaise = raiseActions.length > 0;

  // Find minimum raise action
  const minRaiseAction = hasRaise
    ? raiseActions.reduce(
        (min, curr) => (min[0] < curr[0] ? min : curr),
        raiseActions[0]
      )
    : [0, 0];

  const minRaise = minRaiseAction[1];

  // Calculate bet sizes in big blinds
  const betSizesInBB = [0.5, 0.75, 1, 1.5, 2, 3, 4];
  const maxBetInBB = Math.floor(allInAmount / bigBlind);

  // Adjust bet size if needed
  const currentBetInBB =
    betSize * bigBlind > allInAmount ? maxBetInBB : betSize;

  const handleBetSizeChange = (_: Event, newValue: number | number[]) => {
    setBetSize(newValue as number);
  };

  const handleAction = (action: Action) => {
    onAction(action);
  };

  // Find the exact betAction that matches our chosen betSize, or use a custom amount
  const getBetAction = () => {
    // Find an existing bet action that matches our chosen size
    const exactBetAction = raiseActions.find(
      ([, amount]) => Math.abs(amount - currentBetInBB * bigBlind) < 0.01
    );

    // If we found an exact match, use that action
    if (exactBetAction) {
      return exactBetAction[0];
    }

    // Otherwise, use the bet size as the action
    return currentBetInBB;
  };

  // Create a action for a specific bet size
  const getBetActionForSize = (size: number) => {
    // Find an existing bet action that matches this size
    const exactBetAction = raiseActions.find(
      ([, amount]) => Math.abs(amount - size * bigBlind) < 0.01
    );

    // If we found an exact match, use that action
    if (exactBetAction) {
      return exactBetAction[0];
    }

    // Otherwise, use the bet size as the action
    return size;
  };

  // Check if a specific action is in progress
  const isActionInProgress = (action: Action) => actionInProgress === action;

  // Filter bet sizes to show only those within valid range
  const validBetSizes = betSizesInBB.filter(
    (size) => size <= maxBetInBB && size >= minRaise / bigBlind
  );

  // Limit the number of displayed bet sizes to keep the UI manageable
  const displayedBetSizes = validBetSizes.slice(0, 4);

  // Common button style
  const buttonStyle = {
    fontSize: "0.75rem",
    minWidth: 0,
    padding: "4px 8px",
  };

  return (
    <Paper sx={{ p: 1, borderRadius: 2 }}>
      {/* Slider (only shown when raise is available) */}
      {hasRaise && (
        <Box sx={{ mb: 1 }}>
          <Slider
            size="small"
            value={currentBetInBB}
            onChange={handleBetSizeChange}
            step={0.25}
            min={Math.max(minRaise / bigBlind, 0.5)}
            max={maxBetInBB}
            valueLabelDisplay="auto"
            valueLabelFormat={(value) => `$${(value * bigBlind).toFixed(0)}`}
            disabled={Boolean(actionInProgress)}
          />
        </Box>
      )}

      {/* All actions in one row */}
      <Stack
        direction="row"
        spacing={0.5}
        sx={{
          flexWrap: { xs: "wrap", sm: "nowrap" },
          "& > *": { flex: 1, minWidth: 0 },
        }}
      >
        {/* Basic actions */}
        {hasFold && (
          <Button
            variant="contained"
            color="error"
            size="small"
            onClick={() => handleAction(DiscreteAction.FOLD)}
            disabled={Boolean(actionInProgress)}
            startIcon={
              isActionInProgress(DiscreteAction.FOLD) && (
                <CircularProgress size={14} color="inherit" />
              )
            }
            sx={buttonStyle}
          >
            Fold
          </Button>
        )}

        {hasCheck && (
          <Button
            variant="contained"
            color="primary"
            size="small"
            onClick={() => handleAction(DiscreteAction.CHECK)}
            disabled={Boolean(actionInProgress)}
            startIcon={
              isActionInProgress(DiscreteAction.CHECK) && (
                <CircularProgress size={14} color="inherit" />
              )
            }
            sx={buttonStyle}
          >
            Check
          </Button>
        )}

        {hasCall && (
          <Button
            variant="contained"
            color="primary"
            size="small"
            onClick={() => handleAction(DiscreteAction.CALL)}
            disabled={Boolean(actionInProgress)}
            startIcon={
              isActionInProgress(DiscreteAction.CALL) && (
                <CircularProgress size={14} color="inherit" />
              )
            }
            sx={buttonStyle}
          >
            Call ${callAmount}
          </Button>
        )}

        {/* Quick bet buttons */}
        {hasRaise &&
          displayedBetSizes.map((size) => (
            <Button
              key={size}
              variant="contained"
              color="warning"
              size="small"
              onClick={() => handleAction(getBetActionForSize(size))}
              disabled={Boolean(actionInProgress)}
              sx={buttonStyle}
            >
              {size}x
            </Button>
          ))}

        {/* Custom bet button */}
        {hasRaise && (
          <Button
            variant="contained"
            color="warning"
            size="small"
            onClick={() => handleAction(getBetAction())}
            disabled={Boolean(actionInProgress)}
            sx={{ ...buttonStyle, fontWeight: "bold" }}
          >
            {currentBetInBB.toFixed(1)}x
          </Button>
        )}

        {hasAllIn && (
          <Button
            variant="contained"
            color="secondary"
            size="small"
            onClick={() => handleAction(DiscreteAction.ALL_IN)}
            disabled={Boolean(actionInProgress)}
            startIcon={
              isActionInProgress(DiscreteAction.ALL_IN) && (
                <CircularProgress size={14} color="inherit" />
              )
            }
            sx={buttonStyle}
          >
            All-in
          </Button>
        )}
      </Stack>
    </Paper>
  );
}
