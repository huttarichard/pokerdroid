import { useEffect, useState } from "react";
import {
  Box,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Typography,
  Chip,
  Stack,
  CircularProgress,
  Tooltip,
} from "@mui/material";
import {
  PeopleAlt as PeopleIcon,
  AttachMoney as MoneyIcon,
  Speed as SpeedIcon,
  Memory as MemoryIcon,
} from "@mui/icons-material";
import { getSolutions, SolutionInfo } from "~/services/tree";

interface PathSelectorProps {
  onSubmit: (solution: number) => void;
}

// Group solutions by number of players
type GroupedSolutions = {
  [players: number]: {
    solutions: SolutionInfo[];
    indices: number[]; // Original indices to maintain correct onSubmit
  };
};

function SolutionRow({
  solution,
  onSelect,
}: {
  solution: SolutionInfo;
  originalIndex: number;
  onSelect: () => void;
}) {
  // Format large numbers with K, M, B, T notation
  const formatLargeNumber = (num: number): string => {
    if (num < 1000) return num.toString();
    if (num < 1_000_000) return (num / 1000).toFixed(1) + "K";
    if (num < 1_000_000_000) return (num / 1_000_000).toFixed(1) + "M";
    if (num < 1_000_000_000_000) return (num / 1_000_000_000).toFixed(1) + "B";
    return (num / 1_000_000_000_000).toFixed(1) + "T";
  };

  // Format number with commas for thousands (for more readable values)
  const formatNumber = (num: number): string => {
    return num.toLocaleString();
  };

  // Get average stack size
  const avgStack =
    solution.params.initial_stacks.reduce((a, b) => a + b, 0) /
    solution.params.initial_stacks.length;

  // Calculate big blinds
  const bigBlind = solution.params.sb_amount * 2;
  const stackDepthBB = Math.round(avgStack / bigBlind);

  return (
    <TableRow
      hover
      onClick={onSelect}
      sx={{
        cursor: "pointer",
        transition: "all 0.2s",
        "&:hover": {
          // borderBottom: (theme) => `2px solid ${theme.palette.primary.main}`,
          backgroundColor: (theme) => theme.palette.background.paper,
        },
      }}
    >
      <TableCell>
        <Box
          sx={{
            bgcolor: (theme) => theme.palette.primary.main,
            color: "white",
            borderRadius: 1,
            py: 0.5,
            px: 1,
            display: "inline-flex",
            alignItems: "center",
            gap: 1,
          }}
        >
          <MoneyIcon fontSize="small" />
          <Tooltip title={`${formatNumber(avgStack)} chips / ${bigBlind}bb`}>
            <Typography variant="body2" sx={{ fontWeight: "medium" }}>
              {stackDepthBB}bb
            </Typography>
          </Tooltip>
        </Box>
      </TableCell>
      <TableCell>
        <Stack direction="row" spacing={1} alignItems="center">
          <SpeedIcon fontSize="small" color="action" />
          <Typography variant="body2">
            {formatLargeNumber(solution.iteration)} Iterations
          </Typography>
        </Stack>
      </TableCell>
      <TableCell>
        <Stack direction="row" spacing={0.5} alignItems="center">
          <MemoryIcon fontSize="small" color="action" />
          <Typography variant="body2">
            {formatNumber(solution.nodes)} nodes
          </Typography>
        </Stack>
      </TableCell>
      <TableCell>
        <Box sx={{ display: "flex", flexWrap: "wrap", gap: 0.5 }}>
          {solution.params.bet_sizes.map((size, i) => (
            <Chip
              key={i}
              label={`${size}x POT`}
              size="small"
              color="default"
              sx={{
                height: "24px",
                "& .MuiChip-label": {
                  px: 1,
                  fontSize: "0.75rem",
                },
              }}
            />
          ))}
        </Box>
      </TableCell>
    </TableRow>
  );
}

function PathSelector({ onSubmit }: PathSelectorProps) {
  const [solutions, setSolutions] = useState<SolutionInfo[]>([]);
  const [groupedSolutions, setGroupedSolutions] = useState<GroupedSolutions>(
    {}
  );
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadSolutions = async () => {
      try {
        const data = await getSolutions();
        setSolutions(data);

        // Group solutions by number of players
        const grouped: GroupedSolutions = {};
        data.forEach((solution, index) => {
          const playerCount = solution.params.num_players;
          if (!grouped[playerCount]) {
            grouped[playerCount] = { solutions: [], indices: [] };
          }
          grouped[playerCount].solutions.push(solution);
          grouped[playerCount].indices.push(index);
        });

        // Sort each group by stack depth (highest to lowest)
        Object.keys(grouped).forEach((playerCount) => {
          const group = grouped[Number(playerCount)];
          const sortedIndices = [...Array(group.solutions.length).keys()];

          // Calculate stack depth for each solution
          const stackDepths = group.solutions.map((solution) => {
            const avgStack =
              solution.params.initial_stacks.reduce((a, b) => a + b, 0) /
              solution.params.initial_stacks.length;
            const bigBlind = solution.params.sb_amount * 2;
            return Math.round(avgStack / bigBlind);
          });

          // Sort indices by stack depth (descending)
          sortedIndices.sort((a, b) => stackDepths[b] - stackDepths[a]);

          // Reorder solutions and indices
          const sortedSolutions = sortedIndices.map((i) => group.solutions[i]);
          const sortedOriginalIndices = sortedIndices.map(
            (i) => group.indices[i]
          );

          grouped[Number(playerCount)].solutions = sortedSolutions;
          grouped[Number(playerCount)].indices = sortedOriginalIndices;
        });

        setGroupedSolutions(grouped);
      } catch (err) {
        setError("Failed to load solutions");
        console.error(err);
      } finally {
        setLoading(false);
      }
    };
    loadSolutions();
  }, []);

  if (loading) {
    return (
      <Stack
        alignItems="center"
        justifyContent="center"
        sx={{ height: "100%" }}
      >
        <CircularProgress />
      </Stack>
    );
  }

  if (error) {
    return (
      <Stack
        alignItems="center"
        justifyContent="center"
        sx={{ height: "100%" }}
      >
        <Typography color="error">{error}</Typography>
      </Stack>
    );
  }

  if (solutions.length === 0) {
    return (
      <Box sx={{ p: 3, width: "100%" }}>
        <Typography variant="h5" sx={{ mb: 3 }}>
          Available Solutions
        </Typography>
        <Box sx={{ textAlign: "center", mt: 4 }}>
          <Typography variant="body1" color="text.secondary">
            No solutions available. Create a solution to get started.
          </Typography>
        </Box>
      </Box>
    );
  }

  // Sort player counts numerically
  const playerCounts = Object.keys(groupedSolutions)
    .map(Number)
    .sort((a, b) => a - b);

  return (
    <Box sx={{ p: 3, width: "100%" }}>
      <Typography variant="h5" sx={{ mb: 3 }}>
        Available Solutions
      </Typography>

      {playerCounts.map((playerCount) => (
        <Box key={playerCount} sx={{ mb: 4 }}>
          <Stack
            direction="row"
            spacing={1}
            alignItems="center"
            sx={{
              mb: 2,
              pb: 1,
              borderBottom: (theme) => `1px solid ${theme.palette.divider}`,
            }}
          >
            <PeopleIcon color="primary" />
            <Typography variant="h6">{playerCount} Players</Typography>
            <Typography variant="body2" color="text.secondary" sx={{ ml: 1 }}>
              ({groupedSolutions[playerCount].solutions.length} solutions)
            </Typography>
          </Stack>

          <TableContainer component={Paper} elevation={1} sx={{ mb: 3 }}>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Stack Depth</TableCell>
                  <TableCell>Iterations</TableCell>
                  <TableCell>Nodes</TableCell>
                  <TableCell>Bet Sizes</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {groupedSolutions[playerCount].solutions.map(
                  (solution, index) => (
                    <SolutionRow
                      key={index}
                      solution={solution}
                      originalIndex={
                        groupedSolutions[playerCount].indices[index]
                      }
                      onSelect={() =>
                        onSubmit(groupedSolutions[playerCount].indices[index])
                      }
                    />
                  )
                )}
              </TableBody>
            </Table>
          </TableContainer>
        </Box>
      ))}
    </Box>
  );
}

export default PathSelector;
