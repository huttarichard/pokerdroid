import { Outlet } from "react-router-dom";
import Box from "@mui/material/Box";
import List from "@mui/material/List";
import Divider from "@mui/material/Divider";
// import IconButton from "@mui/material/IconButton";
import Paper from "@mui/material/Paper";
import { ReactLenis } from "lenis/react";

import CssBaseline from "@mui/material/CssBaseline";
import { ThemeProvider } from "@mui/material/styles";
import theme from "./theme";
import Inter from "~/components/inter";

import { mainListItems } from "./items";
import { GlobalStyles } from "@mui/material";

export default function Layout() {
  // const lenis = useLenis(({ scroll }) => {
  //   // called every scroll
  // });

  return (
    <ThemeProvider theme={theme()}>
      <CssBaseline />
      <GlobalStyles
        styles={`
          html.lenis, html.lenis body {
            height: auto;
          }

          .lenis.lenis-smooth {
            scroll-behavior: auto !important;
          }

          .lenis.lenis-smooth [data-lenis-prevent] {
            overscroll-behavior: contain;
          }

          .lenis.lenis-stopped {
            overflow: hidden;
          }

          .lenis.lenis-smooth iframe {
            pointer-events: none;
          }
      `}
      />

      <Inter />

      <Box>
        <Box
          component="aside"
          sx={{
            position: "fixed",
            width: "56px",
            height: "100vh",
            backgroundColor: (theme) => theme.palette.grey[100],
            borderRightColor: (theme) => theme.palette.grey[200],
          }}
        >
          {/* <Box
            sx={{
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              px: [1],
            }}
            component="div"
          >
            <Logo />
          </Box> */}

          <Divider />

          <List component="nav">{mainListItems}</List>
        </Box>

        <Box
          component="main"
          sx={{
            backgroundColor: "white",
            width: "calc(100% - 56px)",
            marginLeft: "56px",
          }}
        >
          <ReactLenis
            root
            options={{
              lerp: 0.5,
            }}
          >
            <Paper
              sx={{
                display: "flex",
                padding: 0,
                flexDirection: "column",
                backgroundColor: (theme) => theme.palette.grey[50],
                minHeight: "100vh",
                borderRadius: 0,
              }}
            >
              <Outlet />
            </Paper>
          </ReactLenis>
        </Box>
      </Box>
    </ThemeProvider>
  );
}

// function Logo() {
//   return (
//     <svg
//       xmlns="http://www.w3.org/2000/svg"
//       viewBox="0 0 61.66 61.2"
//       style={{
//         minWidth: "36px",
//       }}
//     ></svg>
//   );
// }
