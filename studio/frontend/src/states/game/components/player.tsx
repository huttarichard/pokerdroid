import React from "react";
import {
  Box,
  Paper,
  Typography,
  Divider,
  List,
  ListItem,
  ListItemText,
  Chip,
} from "@mui/material";
import { Street } from "~/services/game";

interface PlayerInfoProps {
  playerStack: number;
  opponentStack: number;
  playerPos: number;
  btnPos: number;
  round: number;
  lastWinners: number[];
  street: Street;
}

// Street names based on Go enum: Street.NO_STREET(0), PREFLOP(1), FLOP(2), TURN(3), RIVER(4), FINISHED(5)
const streetNames = ["Unknown", "Preflop", "Flop", "Turn", "River", "Showdown"];

export default function PlayerInfo({
  playerStack,
  opponentStack,
  playerPos,
  btnPos,
  round,
  lastWinners,
  street,
}: PlayerInfoProps) {
  const playerIsBtn = playerPos === btnPos;
  const playerPosition = playerIsBtn ? "Button (Dealer)" : "Big Blind";
  const opponentPosition = playerIsBtn ? "Big Blind" : "Button (Dealer)";

  // Determine winner text
  let winnerText = "";
  if (lastWinners.length > 0) {
    if (lastWinners.includes(playerPos)) {
      winnerText = "You won!";
    } else {
      winnerText = "Opponent won";
    }
  }

  return (
    <Paper
      elevation={3}
      sx={{
        height: "100%",
        borderRadius: 2,
        display: "flex",
        flexDirection: "column",
        fontSize: "0.85rem",
      }}
    >
      <Box sx={{ p: 0.5 }}>
        <Typography
          variant="caption"
          sx={{ fontWeight: "bold", fontSize: "0.9rem" }}
        >
          Game Info
        </Typography>

        <Divider sx={{ my: 0.5, bgcolor: "rgba(255,255,255,0.1)" }} />

        <List disablePadding dense>
          <ListItem dense disablePadding sx={{ py: 0.25 }}>
            <ListItemText
              primary="Round"
              secondary={round}
              primaryTypographyProps={{
                variant: "caption",
                color: "text.secondary",
              }}
              secondaryTypographyProps={{ variant: "caption" }}
              sx={{ my: 0 }}
            />
          </ListItem>

          <ListItem dense disablePadding sx={{ py: 0.25 }}>
            <ListItemText
              primary="Street"
              secondary={streetNames[street] || "Unknown"}
              primaryTypographyProps={{
                variant: "caption",
                color: "text.secondary",
              }}
              secondaryTypographyProps={{ variant: "caption" }}
              sx={{ my: 0 }}
            />
          </ListItem>

          <ListItem dense disablePadding sx={{ py: 0.25 }}>
            <ListItemText
              primary="Your Stack"
              secondary={`$${playerStack}`}
              primaryTypographyProps={{
                variant: "caption",
                color: "text.secondary",
              }}
              secondaryTypographyProps={{ variant: "caption" }}
              sx={{ my: 0 }}
            />
          </ListItem>

          <ListItem dense disablePadding sx={{ py: 0.25 }}>
            <ListItemText
              primary="Opponent Stack"
              secondary={`$${opponentStack}`}
              primaryTypographyProps={{
                variant: "caption",
                color: "text.secondary",
              }}
              secondaryTypographyProps={{ variant: "caption" }}
              sx={{ my: 0 }}
            />
          </ListItem>

          <ListItem dense disablePadding sx={{ py: 0.25 }}>
            <ListItemText
              primary="Your Position"
              secondary={playerPosition}
              primaryTypographyProps={{
                variant: "caption",
                color: "text.secondary",
              }}
              secondaryTypographyProps={{ variant: "caption" }}
              sx={{ my: 0 }}
            />
          </ListItem>

          <ListItem dense disablePadding sx={{ py: 0.25 }}>
            <ListItemText
              primary="Opponent Position"
              secondary={opponentPosition}
              primaryTypographyProps={{
                variant: "caption",
                color: "text.secondary",
              }}
              secondaryTypographyProps={{ variant: "caption" }}
              sx={{ my: 0 }}
            />
          </ListItem>
        </List>

        {winnerText && (
          <Box sx={{ mt: 0.5, textAlign: "center" }}>
            <Chip
              label={winnerText}
              color={winnerText === "You won!" ? "success" : "error"}
              size="small"
              sx={{ px: 0.5, fontSize: "0.65rem", height: "16px" }}
            />
          </Box>
        )}
      </Box>
    </Paper>
  );
}
