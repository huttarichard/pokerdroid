import { Box } from "@mui/material";

interface CellProps {
  hand: string;
  distribution: number[] | undefined;
  actions: number[];
  reach: number;
}

function Cell({ hand, distribution, actions, reach }: CellProps) {
  const actionColors: { [key: number]: string } = {
    [-4]: "rgb(255, 126, 103)", // All In - Dark Red
    [-3]: "rgb(0, 138, 197)", // Fold - Light Gray
    [-2]: "rgb(68, 176, 0)", // Call - Green
    [-1]: "rgb(68, 176, 0)", // Check - Same Green
    [0]: "rgba(128, 128, 128, 0.4)", // NoAction - Gray
  };
  // Helper function to get color for an action
  const getActionColor = (action: number) => {
    if (action > 0) {
      // Scale from darker to brighter red based on action 1-8
      const brightness = 35 + Math.floor((action / 8) * 50); // Goes from 30% to 80%
      return `hsl(0, 100%, ${brightness}%)`; // Consistent red hue, varying brightness
    }
    return actionColors[action] || "rgba(128, 128, 128, 0.4)";
  };

  // If no distribution, show just the hand label
  if (!distribution) {
    return (
      <Box
        sx={{
          width: "100%",
          height: "100%",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          fontSize: 12,
          fontWeight: "bold",
        }}
      >
        {hand}
      </Box>
    );
  }

  return (
    <Box
      sx={{
        width: "100%",
        height: "100%",
        position: "relative",
        p: 0,
        overflow: "hidden",
      }}
    >
      {/* Action bars */}
      <Box
        sx={{
          position: "absolute",
          width: "100%",
          height: `${Math.round(reach * 100)}%`,
          display: "flex",
          flexDirection: "row",
          overflow: "hidden",
        }}
      >
        {distribution.map((prob, idx) => (
          <Box
            key={idx}
            sx={{
              height: "100%",
              width: `${prob * 100}%`,
              bgcolor: getActionColor(actions[idx]),
            }}
          />
        ))}
      </Box>

      {/* Hand label */}
      <Box
        sx={{
          position: "absolute",
          width: "100%",
          height: "100%",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          fontSize: 11,
          fontWeight: "normal",
          color: "black",
          textShadow: "0 0 8px white",
          zIndex: 1,
        }}
      >
        {hand}
      </Box>
    </Box>
  );
}

export default Cell;
