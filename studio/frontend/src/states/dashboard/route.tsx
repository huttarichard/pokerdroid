import { RouteObject } from "react-router-dom";
import {
  Box,
  Grid,
  Typography,
  Paper,
  Button,
  Card,
  CardContent,
  CardActions,
  List,
  ListItem,
  ListItemText,
  Divider,
} from "@mui/material";
import {
  LocalLibraryOutlined,
  SportsEsportsOutlined,
  SearchOutlined,
} from "@mui/icons-material";

const route: RouteObject = {
  index: true,
  Component,
};

export default route;

function Component() {
  // Sample release log data
  const releaseLog = [
    {
      version: "1.0.0",
      date: "2025-03-27",
      description: "Initial release",
    },
  ];

  return (
    <Box sx={{ flexGrow: 1, p: 3 }}>
      <Grid container spacing={4}>
        {/* Release Log Column */}
        <Grid item xs={12} md={5}>
          <Paper elevation={3} sx={{ p: 2, height: "100%" }}>
            <Typography variant="h6" component="h2" gutterBottom>
              Release Log
            </Typography>
            <Divider sx={{ mb: 2 }} />
            <List>
              {releaseLog.map((release, index) => (
                <ListItem key={index} divider={index < releaseLog.length - 1}>
                  <ListItemText
                    primary={`v${release.version} - ${release.date}`}
                    secondary={release.description}
                  />
                </ListItem>
              ))}
            </List>
          </Paper>
        </Grid>

        {/* Actions Column */}
        <Grid item xs={12} md={7}>
          <Grid container spacing={3} direction="column">
            {/* Row 1: Solution Library */}
            <Grid item>
              <Card elevation={3}>
                <CardContent>
                  <Grid container spacing={2} alignItems="center">
                    <Grid item>
                      <LocalLibraryOutlined fontSize="large" color="primary" />
                    </Grid>
                    <Grid item xs>
                      <Typography variant="h6">Solution Library</Typography>
                      <Typography variant="body2" color="text.secondary">
                        Browse through our collection of pre-computed optimal
                        poker strategies
                      </Typography>
                    </Grid>
                  </Grid>
                </CardContent>
                <CardActions>
                  <Button href="/tree" variant="contained" color="primary">
                    Explore Library
                  </Button>
                </CardActions>
              </Card>
            </Grid>

            {/* Row 2: Play Against Agent */}
            <Grid item>
              <Card elevation={3}>
                <CardContent>
                  <Grid container spacing={2} alignItems="center">
                    <Grid item>
                      <SportsEsportsOutlined
                        fontSize="large"
                        color="secondary"
                      />
                    </Grid>
                    <Grid item xs>
                      <Typography variant="h6">Play Against Agent</Typography>
                      <Typography variant="body2" color="text.secondary">
                        Test your skills against our AI poker agent
                      </Typography>
                    </Grid>
                  </Grid>
                </CardContent>
                <CardActions>
                  <Button href="/game" variant="contained" color="secondary">
                    Play Now
                  </Button>
                </CardActions>
              </Card>
            </Grid>

            {/* Row 3: Play Against Agent with Search */}
            <Grid item>
              <Card elevation={3}>
                <CardContent>
                  <Grid container spacing={2} alignItems="center">
                    <Grid item>
                      <SearchOutlined fontSize="large" color="info" />
                    </Grid>
                    <Grid item xs>
                      <Typography variant="h6">Advanced Play Mode</Typography>
                      <Typography variant="body2" color="text.secondary">
                        Challenge our enhanced AI that uses real-time search for
                        decision making
                      </Typography>
                    </Grid>
                  </Grid>
                </CardContent>
                <CardActions>
                  <Button
                    href="/game?mode=search"
                    variant="contained"
                    color="info"
                  >
                    Play with Search
                  </Button>
                </CardActions>
              </Card>
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </Box>
  );
}
