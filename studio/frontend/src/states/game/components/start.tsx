import React from "react";
import {
  Box,
  Button,
  Typography,
  CircularProgress,
  Paper,
} from "@mui/material";

interface StartScreenProps {
  loading: boolean;
  error: string | null;
  onStart: () => void;
}

export default function StartScreen({
  loading,
  error,
  onStart,
}: StartScreenProps) {
  return (
    <Box
      sx={{
        display: "flex",
        flexDirection: "column",
        justifyContent: "center",
        alignItems: "center",
        height: "100%",
        width: "100%",
      }}
    >
      <Paper
        elevation={2}
        sx={{
          p: 4,
          maxWidth: 500,
          width: "100%",
          borderRadius: 2,
          textAlign: "center",
        }}
      >
        <Typography variant="h4" gutterBottom>
          Poker Droid
        </Typography>

        <Typography variant="h6" sx={{ mb: 3 }}>
          Heads-Up No-Limit Hold'em
        </Typography>

        <Typography variant="body1" sx={{ mb: 4 }}>
          Play heads-up poker against our AI opponent. Test your skills against
          a bot trained with Counterfactual Regret Minimization.
        </Typography>

        {error && (
          <Typography variant="body2" color="error" sx={{ mb: 2 }}>
            Error: {error}
          </Typography>
        )}

        <Button
          variant="contained"
          color="primary"
          size="large"
          onClick={onStart}
          disabled={loading}
          sx={{
            px: 4,
            py: 1.5,
            borderRadius: 2,
            fontSize: "1.1rem",
          }}
        >
          {loading ? (
            <CircularProgress size={24} color="inherit" />
          ) : (
            "Start Game"
          )}
        </Button>
      </Paper>
    </Box>
  );
}
