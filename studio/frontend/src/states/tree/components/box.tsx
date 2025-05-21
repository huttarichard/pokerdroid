import { useState } from "react";
import { Paper, Box, Stack, Button, Typography } from "@mui/material";
import CardSelectDrawer from "./drawer";
import { TreeResponse, Action, discreteActionToString } from "~/services/tree";
import { formatCard } from "./details";

interface StateBoxProps {
  node: TreeResponse;
  action: Action;
  isSelected: boolean;
  onClick: () => void;
  onAction: (index?: number) => void;
  onChanceSelect: (cards: string[]) => void;
  allSelectedCards: string[];
}

function StateBox({
  node,
  action,
  isSelected,
  onClick,
  onAction,
  onChanceSelect,
  allSelectedCards,
}: StateBoxProps) {
  const [drawerOpen, setDrawerOpen] = useState(false);

  const getContent = () => {
    // Show loading state when node is null and box is selected
    if (isSelected && !node) {
      return (
        <Box
          sx={{
            fontSize: 11,
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            height: "100%",
          }}
        >
          Loading...
        </Box>
      );
    }

    switch (node.kind) {
      case "root":
        return <Box sx={{ fontSize: 11 }}>Game Tree Info</Box>;
      case "chance":
        return (
          <Stack spacing={0.5}>
            <Box
              sx={{
                fontSize: 10,
                cursor: "pointer",
                "&:hover": { bgcolor: "action.hover" },
                p: 0.5,
                borderRadius: 1,
                border: 1,
                borderColor: "divider",
                textAlign: "center",
              }}
              onClick={(e) => {
                e.stopPropagation();
                setDrawerOpen(true);
              }}
            >
              Select {node.street} cards
            </Box>

            {/* Display selected cards if available */}
            {(action.cards?.length || 0) > 0 && (
              <Box
                sx={{
                  display: "flex",
                  flexWrap: "wrap",
                  gap: 0.5,
                  justifyContent: "center",
                  minHeight: 20,
                }}
              >
                {action.cards!.map((card, idx) => (
                  <Typography
                    key={idx}
                    sx={{
                      fontSize: 10,
                      fontWeight: "bold",
                      color:
                        card[1] === "h" || card[1] === "d"
                          ? "error.main"
                          : "text.primary",
                    }}
                  >
                    {formatCard(card)}
                  </Typography>
                ))}
              </Box>
            )}
          </Stack>
        );
      case "player": {
        // Group actions by type
        const standardActions: [number, number][] = []; // [action value, index]
        const raiseActions: [number, number][] = []; // [action value, index]

        node.actions?.forEach((action, index) => {
          if (action <= 0) {
            // All-in (-4), Fold (-3), Call (-2), Check (-1) go to standard actions
            standardActions.push([action, index]);
          } else {
            // All raise actions (> 0) go to raise actions
            raiseActions.push([action, index]);
          }
        });

        // Sort standard actions: Check, Call, Fold, All-in
        standardActions.sort((a, b) => b[0] - a[0]);

        // Sort raise actions: smallest to largest
        raiseActions.sort((a, b) => a[0] - b[0]);

        return (
          <Stack spacing={0.5}>
            {/* Standard actions row */}
            {standardActions.length > 0 && (
              <Box
                sx={{
                  display: "flex",
                  justifyContent: "stretch",
                  width: "100%",
                  gap: 0.5,
                }}
              >
                {standardActions.map(([action, index]) => (
                  <Button
                    key={action}
                    size="small"
                    variant="outlined"
                    onClick={(e) => {
                      e.stopPropagation();
                      onAction(index);
                    }}
                    sx={{
                      fontSize: 10,
                      py: 0.1,
                      px: 0.5,
                      minWidth: 0,
                      minHeight: 0,
                      height: 18,
                      lineHeight: 1,
                      textTransform: "none",
                      width: "100%",
                    }}
                  >
                    {discreteActionToString(action)}
                  </Button>
                ))}
              </Box>
            )}

            <Box
              sx={{
                display: "flex",
                justifyContent: "stretch",
                width: "100%",
                gap: 0.5,
              }}
            >
              {/* Raise actions - one per row */}
              {raiseActions.map(([action, index]) => (
                <Button
                  key={action}
                  size="small"
                  variant="outlined"
                  onClick={(e) => {
                    e.stopPropagation();
                    onAction(index);
                  }}
                  sx={{
                    fontSize: 10,
                    py: 0.1,
                    px: 0.5,
                    minWidth: 0,
                    minHeight: 0,
                    height: 18,
                    lineHeight: 1,
                    textTransform: "none",
                    width: "100%",
                  }}
                >
                  {discreteActionToString(action)}
                </Button>
              ))}
            </Box>
          </Stack>
        );
      }
      case "terminal":
        return <Box sx={{ fontSize: 11 }}>Terminal Node</Box>;
      default:
        return null;
    }
  };

  return (
    <>
      <Paper
        elevation={isSelected ? 4 : 1}
        sx={{
          minWidth: 160,
          minHeight: 50,
          cursor: "pointer",
          bgcolor: isSelected ? "action.selected" : "background.paper",
          "&:hover": { bgcolor: "action.hover" },
          p: 0.75,
          display: "flex",
          flexDirection: "column",
        }}
        onClick={onClick}
      >
        <Box
          sx={{
            fontSize: 9,
            fontWeight: "bold",
            textTransform: "uppercase",
            borderBottom: 1,
            borderColor: "divider",
            pb: 0.5,
            textAlign: "center",
          }}
        >
          {node.kind} Node
        </Box>
        {getContent()}
      </Paper>

      <CardSelectDrawer
        open={drawerOpen}
        street={node.street}
        selectedCards={action.cards || []}
        onClose={() => setDrawerOpen(false)}
        onSelect={onChanceSelect}
        allSelectedCards={allSelectedCards}
      />
    </>
  );
}

export default StateBox;
