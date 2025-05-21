import { Box, Stack, Paper, Grid } from "@mui/material";
import { TreeResponse } from "~/services/tree";
import { RANKS } from "./constants";
import Cell from "./cell";
import CellDetails from "./details";
import { useState, useEffect } from "react";

interface PokerMatrixProps {
  node: TreeResponse;
}

function PokerMatrix({ node }: PokerMatrixProps) {
  const [selectedCell, setSelectedCell] = useState<{
    row: number;
    col: number;
  } | null>(null);

  // Reset selection when node changes
  useEffect(() => {
    setSelectedCell(null);
  }, [node]);

  if (!node.matrix) return null;

  const handleCellClick = (row: number, col: number) => {
    setSelectedCell((prev) =>
      prev?.row === row && prev?.col === col ? null : { row, col }
    );
  };

  const selectedCellData =
    selectedCell && node.matrix[selectedCell.row][selectedCell.col];

  return (
    <Grid container spacing={2}>
      {/* Left side - Matrix */}
      <Grid item xs={7}>
        <Box sx={{ position: "relative" }}>
          {/* Matrix with row headers */}
          <Stack spacing={0.25}>
            {Array.from({ length: 13 }, (_, row) => {
              const rowRank = RANKS[row];
              return (
                <Stack key={row} direction="row" spacing={0.25}>
                  {Array.from({ length: 13 }, (_, col) => {
                    const colRank = RANKS[col];
                    const cell = node.matrix?.[row]?.[col];
                    if (!cell) return null;

                    const isSelected =
                      selectedCell?.row === row && selectedCell?.col === col;

                    // Since RANKS array is already ordered from A to 2,
                    // we can use array indices for comparison
                    const [firstRank, secondRank] =
                      row <= col ? [rowRank, colRank] : [colRank, rowRank];

                    const suited = row < col ? "s" : row > col ? "o" : "";

                    return (
                      <Paper
                        key={`${row}-${col}`}
                        sx={{
                          width: 70,
                          height: 40,
                          cursor: "pointer",
                          overflow: "hidden",
                          "&:hover": {
                            bgcolor: "action.hover",
                          },
                          bgcolor: isSelected
                            ? "action.selected"
                            : "background.paper",
                          p: 0,
                          border: 0,
                          borderRadius: "3px",
                        }}
                        onClick={() => handleCellClick(row, col)}
                      >
                        <Cell
                          hand={`${firstRank.toUpperCase()}${secondRank.toUpperCase()}${suited}`}
                          distribution={cell.avg_strat}
                          reach={cell.reach}
                          actions={node.actions || []}
                        />
                      </Paper>
                    );
                  })}
                </Stack>
              );
            })}
          </Stack>
        </Box>
      </Grid>

      {/* Right side - Cell Details */}
      <Grid item xs={5}>
        {selectedCellData && node.actions && (
          <CellDetails
            cards={selectedCellData.cards || []}
            policies={selectedCellData.policies || []}
            policy={selectedCellData.policy}
            avgStrat={selectedCellData.avg_strat || []}
            reachProbs={selectedCellData.reach_probs || []}
            clusters={selectedCellData.clusters || []}
            actions={node.actions}
          />
        )}
      </Grid>
    </Grid>
  );
}

export default PokerMatrix;
