import { useState, useEffect } from "react";
import { Drawer, Stack, Box, Typography, Grid, Paper } from "@mui/material";
import { RANKS, SUITS } from "./constants";

interface CardSelectDrawerProps {
  open: boolean;
  street?: string;
  selectedCards?: string[];
  onClose: () => void;
  onSelect: (cards: string[]) => void;
  allSelectedCards?: string[];
}

function CardSelectDrawer({
  open,
  street,
  selectedCards = [],
  onClose,
  onSelect,
  allSelectedCards = [],
}: CardSelectDrawerProps) {
  const [selectedCardState, setSelectedCardState] =
    useState<string[]>(selectedCards);

  // Update selected cards when the prop changes or drawer opens
  useEffect(() => {
    if (open) {
      setSelectedCardState(selectedCards);
    }
  }, [open, selectedCards]);

  const getRequiredCards = () => {
    switch (street?.toLowerCase()) {
      case "flop":
        return 3;
      case "turn":
        return 1;
      case "river":
        return 1;
      default:
        return 0;
    }
  };

  // Simple check if a card is already used somewhere else in the game tree
  const isCardUsedElsewhere = (card: string): boolean => {
    return allSelectedCards.includes(card);
  };

  const handleCardSelect = (rank: string, suit: keyof typeof SUITS) => {
    const card = `${rank}${SUITS[suit]}`; // Format as "as", "kh", etc.

    // Don't allow selection of cards that are already used elsewhere
    if (isCardUsedElsewhere(card)) {
      return;
    }

    setSelectedCardState((prev) => {
      if (prev.includes(card)) {
        return prev.filter((c) => c !== card);
      }
      if (prev.length < getRequiredCards()) {
        const newSelection = [...prev, card];

        // Auto-submit when we reach the required number of cards
        if (newSelection.length === getRequiredCards()) {
          setTimeout(() => {
            onSelect(newSelection);
            onClose();
          }, 300); // Short delay for visual feedback
        }

        return newSelection;
      }
      return prev;
    });
  };

  return (
    <Drawer
      anchor="right"
      open={open}
      onClose={onClose}
      PaperProps={{ sx: { width: "300px", p: 2 } }}
    >
      <Stack spacing={2}>
        <Box
          sx={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
          }}
        >
          <Typography variant="h6">
            Select {street} Cards
            <Typography variant="body1" component="span" sx={{ ml: 1 }}>
              ({selectedCardState.length}/{getRequiredCards()})
            </Typography>
          </Typography>
        </Box>

        <Grid container spacing={1}>
          {/* Each row is one rank */}
          {RANKS.map((rank) => (
            <Grid item xs={12} key={rank} sx={{ py: 0.2 }}>
              <Stack direction="row" spacing={1} justifyContent="flex-start">
                {/* Each column is one suit */}
                {(Object.keys(SUITS) as Array<keyof typeof SUITS>).map(
                  (suit) => {
                    const card = `${rank}${SUITS[suit]}`;
                    const isSelected = selectedCardState.includes(card);
                    const isUsedElsewhere = isCardUsedElsewhere(card);
                    const displayCard = `${rank.toUpperCase()}${suit}`;
                    return (
                      <Paper
                        key={card}
                        sx={{
                          width: 35,
                          height: 38,
                          display: "flex",
                          p: 0.2,
                          alignItems: "center",
                          justifyContent: "center",
                          cursor: isUsedElsewhere ? "not-allowed" : "pointer",
                          bgcolor: isUsedElsewhere
                            ? "action.disabledBackground"
                            : isSelected
                              ? "primary.main"
                              : "background.paper",
                          color: isUsedElsewhere
                            ? "text.disabled"
                            : isSelected
                              ? "primary.contrastText"
                              : suit === "♥" || suit === "♦"
                                ? "error.main"
                                : "text.primary",
                          "&:hover": {
                            bgcolor: isUsedElsewhere
                              ? "action.disabledBackground"
                              : isSelected
                                ? "primary.dark"
                                : "action.hover",
                          },
                          opacity: isUsedElsewhere ? 0.5 : 1,
                          userSelect: "none",
                          fontSize: "0.8rem",
                        }}
                        onClick={() => {
                          if (!isUsedElsewhere) {
                            handleCardSelect(rank, suit);
                          }
                        }}
                      >
                        {displayCard}
                      </Paper>
                    );
                  }
                )}
              </Stack>
            </Grid>
          ))}
        </Grid>
      </Stack>
    </Drawer>
  );
}

export default CardSelectDrawer;
