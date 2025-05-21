import { Box, Paper, Stack, Typography } from "@mui/material";
import { Policy } from "~/services/tree";
import { discreteActionToString } from "~/services/tree";
import { SUITS } from "./constants";

interface CellDetailsProps {
  cards: string[][];
  policies: Policy[];
  policy: Policy;
  avgStrat: number[];
  reachProbs: number[];
  clusters: number[];
  actions: number[];
}

export function formatCard(card: string): string {
  const rank = card[0].toUpperCase();
  const suitCode = card[1];
  const suitSymbol =
    Object.entries(SUITS).find(([, code]) => code === suitCode)?.[0] ||
    suitCode;
  return `${rank}${suitSymbol}`;
}

function CellDetails({
  cards,
  policies,
  avgStrat,
  reachProbs,
  clusters,
  actions,
}: CellDetailsProps) {
  if (
    !avgStrat ||
    !policies ||
    !cards ||
    !reachProbs ||
    !clusters ||
    !actions
  ) {
    return null;
  }

  const getNormalizedStrategy = (policy: Policy) => {
    if (!policy.strategy_sum) return [];
    const total = policy.strategy_sum.reduce((a, b) => a + b, 0);
    return total > 0 ? policy.strategy_sum.map((v) => v / total) : [];
  };

  return (
    <Paper sx={{ p: 2, height: "100%", overflow: "auto" }}>
      <Stack spacing={2}>
        {/* Average Strategy Section */}
        <Box>
          <Typography variant="subtitle1" gutterBottom>
            Average Strategy
          </Typography>
          <Stack
            direction="row"
            sx={{
              width: "100%",
              gap: 1,
            }}
          >
            {avgStrat.map((prob, idx) => (
              <Box
                key={idx}
                sx={{
                  flex: 1,
                  minWidth: 0,
                  px: 1,
                }}
              >
                <Typography variant="caption" display="block" noWrap>
                  {discreteActionToString(actions[idx] || 0)}
                </Typography>
                <Typography variant="body2" noWrap>
                  {(prob * 100).toFixed(1)}%
                </Typography>
              </Box>
            ))}
          </Stack>
        </Box>

        {/* Individual Combinations Section */}
        <Box>
          <Typography variant="subtitle1" gutterBottom>
            Card Combinations
          </Typography>
          <Stack spacing={1}>
            {cards.map((combo, idx) => {
              const policy = policies[idx];
              if (!policy) return null;

              return (
                <Paper key={idx} variant="outlined" sx={{ p: 1 }}>
                  <Stack
                    direction="row"
                    justifyContent="space-between"
                    alignItems="center"
                    sx={{ borderBottom: "1px solid #f7f7f7" }}
                  >
                    <Typography variant="subtitle2">
                      {combo.map(formatCard).join(" ")}
                    </Typography>
                    <Typography variant="caption">
                      Cluster: {clusters[idx]} | Iter: {policy.iteration} |
                      Reach: {(reachProbs[idx] * 100).toFixed(2)}%
                    </Typography>
                  </Stack>
                  <Stack
                    direction="row"
                    sx={{
                      width: "100%",
                      mt: 0.5,
                    }}
                  >
                    {getNormalizedStrategy(policy).map((value, actionIdx) => (
                      <Box
                        key={actionIdx}
                        sx={{
                          flex: 1,
                          minWidth: 0,
                          px: 1,
                          borderRight:
                            actionIdx < getNormalizedStrategy(policy).length - 1
                              ? 1
                              : 0,
                          borderColor: "divider",
                        }}
                      >
                        <Typography variant="caption" display="block" noWrap>
                          {discreteActionToString(actions[actionIdx] || 0)}
                        </Typography>
                        <Stack
                          direction="row"
                          justifyContent="space-between"
                          alignItems="baseline"
                          spacing={0.5}
                        >
                          <Typography variant="body2" noWrap>
                            {(value * 100).toFixed(1)}%
                          </Typography>
                          <Typography
                            variant="caption"
                            color={
                              policy.baseline[actionIdx] >= 0
                                ? "success.main"
                                : "error.main"
                            }
                            noWrap
                          >
                            {policy.baseline[actionIdx].toFixed(2)}
                          </Typography>
                        </Stack>
                      </Box>
                    ))}
                  </Stack>
                </Paper>
              );
            })}
          </Stack>
        </Box>
      </Stack>
    </Paper>
  );
}

export default CellDetails;
