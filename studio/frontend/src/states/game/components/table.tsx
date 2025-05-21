import { Box, Typography, Paper } from "@mui/material";
import Card from "../../../components/cards/card";
import { Street } from "~/services/game";

interface TableProps {
  pot: number;
  community: string[];
  playerCards: string[];
  playerTurn: boolean;
  street: Street;
}

// Street names based on Go enum: Street.NO_STREET(0), PREFLOP(1), FLOP(2), TURN(3), RIVER(4), FINISHED(5)
const streetNames = ["Unknown", "Preflop", "Flop", "Turn", "River", "Showdown"];

export default function Table({
  pot,
  community,
  playerCards,
  playerTurn,
  street,
}: TableProps) {
  return (
    <Box
      sx={{
        position: "relative",
        height: "100%",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      {/* Poker table */}
      <Paper
        elevation={3}
        sx={{
          backgroundColor: "#036b3c",
          borderRadius: "50%",
          border: "12px solid #70432c",
          width: "90%",
          height: "85%",
          maxHeight: "500px",
          position: "relative",
          display: "flex",
          flexDirection: "column",
          justifyContent: "space-between",
          alignItems: "center",
          padding: 2,
          m: "auto",
          overflow: "hidden",
        }}
      >
        {/* Opponent */}
        <Box
          sx={{
            width: "100%",
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            mt: 1,
            zIndex: 2,
          }}
        >
          <Typography variant="subtitle2" sx={{ color: "white", mb: 0.5 }}>
            Opponent
          </Typography>
          <Box sx={{ display: "flex", gap: 0.5 }}>
            <Card card="back" height="70px" />
            <Card card="back" height="70px" />
          </Box>
        </Box>

        {/* Center - Community cards and pot */}
        <Box
          sx={{
            position: "absolute",
            top: "50%",
            left: "50%",
            transform: "translate(-50%, -50%)",
            width: "100%",
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            justifyContent: "center",
            zIndex: 1,
          }}
        >
          <Typography variant="subtitle2" sx={{ color: "white", mb: 0.5 }}>
            {streetNames[street] || "Unknown"} - Pot: ${pot}
          </Typography>
          <Box
            sx={{
              display: "flex",
              gap: 0.5,
              justifyContent: "center",
              flexWrap: "wrap",
              maxWidth: "100%",
            }}
          >
            {community.length > 0
              ? community.map((card, i) => (
                  <Card key={i} card={card} height="70px" />
                ))
              : Array(5)
                  .fill(null)
                  .map((_, i) => (
                    <Box
                      key={i}
                      sx={{
                        height: "70px",
                        width: "50px",
                        border: "1px dashed rgba(255,255,255,0.3)",
                        borderRadius: "4px",
                      }}
                    />
                  ))}
          </Box>
        </Box>

        {/* Player */}
        <Box
          sx={{
            width: "100%",
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            mb: 1,
            zIndex: 2,
          }}
        >
          <Box sx={{ display: "flex", gap: 0.5 }}>
            {playerCards.map((card, i) => (
              <Card key={i} card={card} height="70px" />
            ))}
          </Box>
          <Typography variant="subtitle2" sx={{ color: "white", mt: 0.5 }}>
            {playerTurn ? "Your Turn" : "Your Hand"}
          </Typography>
        </Box>
      </Paper>
    </Box>
  );
}
