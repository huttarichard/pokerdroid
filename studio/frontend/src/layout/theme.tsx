import { createTheme, ThemeOptions, alpha } from "@mui/material/styles";
import {
  Link as RouterLink,
  LinkProps as RouterLinkProps,
} from "react-router-dom";
import { LinkProps } from "@mui/material/Link";
import { forwardRef } from "react";

declare module "@mui/material/styles/createPalette" {
  interface ColorRange {
    50: string;
    100: string;
    200: string;
    300: string;
    400: string;
    500: string;
    600: string;
    700: string;
    800: string;
    900: string;
  }

  interface PaletteColor extends ColorRange {}
}

export const brand = {
  50: "hsl(220, 100%, 97%)",
  100: "hsl(220, 100%, 90%)",
  200: "hsl(220, 100%, 80%)",
  300: "hsl(220, 100%, 75%)",
  400: "hsl(220, 100%, 65%)",
  500: "hsl(220, 100%, 50%)",
  600: "hsl(220, 100%, 40%)",
  700: "hsl(220, 100%, 35%)",
  800: "hsl(220, 100%, 15%)",
  900: "hsl(220, 100%, 5%)",
};

export const gray = {
  50: "hsl(0, 0%, 96%)",
  100: "hsl(0, 0%, 94%)",
  200: "hsl(0, 0%, 88%)",
  300: "hsl(0, 0%, 80%)",
  400: "hsl(0, 0%, 65%)",
  500: "hsl(0, 0%, 42%)",
  600: "hsl(0, 0%, 35%)",
  700: "hsl(0, 0%, 25%)",
  800: "hsl(0, 0%, 10%)",
  900: "hsl(0, 0%, 5%)",
};

export const green = {
  50: "hsl(120, 80%, 98%)",
  100: "hsl(120, 75%, 94%)",
  200: "hsl(120, 75%, 87%)",
  300: "hsl(120, 61%, 77%)",
  400: "hsl(120, 44%, 53%)",
  500: "hsl(120, 59%, 30%)",
  600: "hsl(120, 70%, 25%)",
  700: "hsl(120, 75%, 16%)",
  800: "hsl(120, 84%, 10%)",
  900: "hsl(120, 87%, 6%)",
};

export const orange = {
  50: "hsl(45, 100%, 97%)",
  100: "hsl(45, 92%, 90%)",
  200: "hsl(45, 94%, 80%)",
  300: "hsl(45, 90%, 65%)",
  400: "hsl(45, 90%, 40%)",
  500: "hsl(45, 90%, 35%)",
  600: "hsl(45, 91%, 25%)",
  700: "hsl(45, 94%, 20%)",
  800: "hsl(45, 95%, 16%)",
  900: "hsl(45, 93%, 12%)",
};

export const red = {
  50: "hsl(0, 100%, 97%)",
  100: "hsl(0, 92%, 90%)",
  200: "hsl(0, 94%, 80%)",
  300: "hsl(0, 90%, 65%)",
  400: "hsl(0, 90%, 40%)",
  500: "hsl(0, 90%, 30%)",
  600: "hsl(0, 91%, 25%)",
  700: "hsl(0, 94%, 20%)",
  800: "hsl(0, 95%, 16%)",
  900: "hsl(0, 93%, 12%)",
};

const LinkBehavior = forwardRef<
  HTMLAnchorElement,
  Omit<RouterLinkProps, "to"> & { href: RouterLinkProps["to"] }
>((props, ref) => {
  const { href, ...other } = props;
  // Map href (Material UI) -> to (react-router)
  return <RouterLink ref={ref} to={href} {...other} />;
});

export default function makeTheme(): ThemeOptions {
  const pxToRem = (px: number) => `${px / 16}rem`;

  return createTheme({
    palette: {
      mode: "light",
      primary: {
        light: brand[200],
        main: brand[500],
        dark: brand[800],
        contrastText: brand[50],
      },
      info: {
        light: brand[100],
        main: brand[300],
        dark: brand[600],
        contrastText: gray[50],
      },
      warning: {
        light: orange[300],
        main: orange[400],
        dark: orange[800],
      },
      error: {
        light: red[300],
        main: red[400],
        dark: red[800],
      },
      success: {
        light: green[300],
        main: green[400],
        dark: green[800],
      },
      grey: {
        ...gray,
      },
      divider: alpha(gray[300], 0.5),
      background: {
        default: "hsl(0, 0%, 100%)",
        paper: gray[100],
      },
      text: {
        primary: gray[800],
        secondary: gray[600],
      },
      action: {
        selected: `${alpha(brand[200], 0.2)}`,
      },
    },
    typography: {
      fontFamily: ['"Inter", "sans-serif"'].join(","),
      h1: {
        fontSize: pxToRem(30),
        fontWeight: 600,
        lineHeight: 1.2,
        letterSpacing: -0.5,
      },
      h2: {
        fontSize: pxToRem(26),
        fontWeight: 600,
        lineHeight: 1.2,
      },
      h3: {
        fontSize: pxToRem(22),
        fontWeight: 600,
        lineHeight: 1.2,
      },
      h4: {
        fontSize: pxToRem(20),
        fontWeight: 500,
        lineHeight: 1.5,
      },
      h5: {
        fontSize: pxToRem(18),
        fontWeight: 500,
      },
      h6: {
        fontSize: pxToRem(16),
        fontWeight: 500,
      },
      subtitle1: {
        fontSize: pxToRem(18),
        fontWeight: 600,
      },
      subtitle2: {
        fontSize: pxToRem(16),
        fontWeight: 500,
      },
      body1: {
        fontSize: pxToRem(16),
        fontWeight: 400,
      },
      body2: {
        fontSize: pxToRem(14),
        fontWeight: 400,
      },
      caption: {
        fontSize: pxToRem(12),
        fontWeight: 400,
      },
    },
    shape: {
      borderRadius: 6,
    },
    components: {
      MuiAccordion: {
        defaultProps: {
          elevation: 0,
          disableGutters: true,
        },
        styleOverrides: {
          root: {
            padding: 8,
            overflow: "clip",
            backgroundColor: "hsl(0, 0%, 100%)",
            border: "1px solid",
            borderColor: gray[100],
            ":before": {
              backgroundColor: "transparent",
            },
            "&:first-of-type": {
              borderTopLeftRadius: 10,
              borderTopRightRadius: 10,
            },
            "&:last-of-type": {
              borderBottomLeftRadius: 10,
              borderBottomRightRadius: 10,
            },
          },
        },
      },
      MuiAccordionSummary: {
        styleOverrides: {
          root: {
            border: "none",
            borderRadius: 8,
            "&:hover": { backgroundColor: gray[100] },
            "&:focus-visible": { backgroundColor: "transparent" },
          },
        },
      },
      MuiAccordionDetails: {
        styleOverrides: {
          root: { mb: 20, border: "none" },
        },
      },
      MuiButtonBase: {
        defaultProps: {
          LinkComponent: LinkBehavior,
          disableTouchRipple: true,
          disableRipple: true,
        },
        styleOverrides: {
          root: {
            boxSizing: "border-box",
            transition: "all 100ms ease",
            "&:focus-visible": {
              outline: `3px solid ${alpha(brand[400], 0.5)}`,
              outlineOffset: "2px",
            },
          },
        },
      },
      MuiButton: {
        defaultProps: {
          size: "small",
        },
        styleOverrides: {
          root: {
            boxShadow: "none",
            borderRadius: 6,
            textTransform: "none",
            variants: [
              {
                props: {
                  size: "small",
                },
                style: {
                  height: "2rem", // 32px
                  padding: "0 0.5rem",
                },
              },
              {
                props: {
                  size: "medium",
                },
                style: {
                  height: "2.5rem", // 40px
                },
              },
              {
                props: {
                  color: "primary",
                  variant: "contained",
                },
                style: {
                  color: "white",
                  backgroundColor: brand[400],
                  backgroundImage: `linear-gradient(to bottom, ${alpha(brand[400], 0.6)}, ${brand[400]})`,
                  boxShadow: `inset 0 2px 0 ${alpha(brand[200], 0.2)}, inset 0 -2px 0 ${alpha(brand[600], 0.4)}`,
                  border: `1px solid ${brand[500]}`,
                  "&:hover": {
                    backgroundColor: brand[600],
                    boxShadow: "none",
                  },
                  "&:active": {
                    backgroundColor: brand[700],
                    boxShadow: `inset 0 2.5px 0 ${alpha(brand[700], 0.4)}`,
                  },
                },
              },
              {
                props: {
                  variant: "outlined",
                },
                style: {
                  color: brand[700],
                  backgroundColor: alpha(brand[300], 0.1),
                  borderColor: alpha(brand[200], 0.8),
                  "&:hover": {
                    backgroundColor: alpha(brand[300], 0.2),
                    borderColor: alpha(brand[300], 0.5),
                    boxShadow: "none",
                  },
                  "&:active": {
                    backgroundColor: alpha(brand[300], 0.3),
                    boxShadow: `inset 0 2.5px 0 ${alpha(brand[400], 0.2)}`,
                    backgroundImage: "none",
                  },
                },
              },
              {
                props: {
                  color: "secondary",
                  variant: "outlined",
                },
                style: {
                  backgroundColor: alpha(gray[300], 0.1),
                  borderColor: alpha(gray[300], 0.5),
                  color: gray[700],
                  "&:hover": {
                    backgroundColor: alpha(gray[300], 0.3),
                    borderColor: alpha(gray[300], 0.5),
                    boxShadow: "none",
                  },
                  "&:active": {
                    backgroundColor: alpha(gray[300], 0.4),
                    boxShadow: `inset 0 2.5px 0 ${alpha(gray[400], 0.2)}`,
                    backgroundImage: "none",
                  },
                },
              },
              {
                props: {
                  color: "primary",
                  variant: "text",
                },
                style: {
                  color: brand[700],
                  "&:hover": {
                    backgroundColor: alpha(brand[300], 0.3),
                  },
                },
              },
              {
                props: {
                  color: "info",
                  variant: "text",
                },
                style: {
                  color: gray[700],
                  "&:hover": {
                    backgroundColor: alpha(gray[300], 0.3),
                  },
                },
              },
            ],
          },
        },
      },
      MuiCard: {
        styleOverrides: {
          root: {
            transition: "all 100ms ease",
            backgroundColor: gray[50],
            borderRadius: 6,
            border: `1px solid ${alpha(gray[200], 0.5)}`,
            boxShadow: "none",
            variants: [
              {
                props: {
                  variant: "outlined",
                },
                style: {
                  border: `1px solid ${gray[200]}`,
                  boxShadow: "none",
                  background: `linear-gradient(to bottom, hsl(0, 0%, 100%), ${gray[50]})`,
                },
              },
            ],
          },
        },
      },
      MuiChip: {
        styleOverrides: {
          root: {
            py: 1.5,
            px: 0.5,
            border: "1px solid",
            borderColor: brand[100],
            fontWeight: 600,
            backgroundColor: brand[50],
            "&:hover": {
              backgroundColor: brand[500],
            },
            "&:focus-visible": {
              borderColor: brand[300],
              backgroundColor: brand[200],
            },
            "& .MuiChip-label": {
              color: brand[500],
            },
            "& .MuiChip-icon": {
              color: brand[500],
            },
          },
        },
      },
      MuiDivider: {
        styleOverrides: {
          root: {
            borderColor: `${alpha(gray[200], 0.8)}`,
          },
        },
      },
      MuiFormLabel: {
        styleOverrides: {
          root: ({ theme }) => ({
            typography: theme.typography.caption,
            marginBottom: 8,
          }),
        },
      },
      MuiIconButton: {
        styleOverrides: {
          root: {
            color: brand[500],
            "&:hover": {
              backgroundColor: alpha(brand[300], 0.3),
              borderColor: brand[200],
            },
            variants: [
              {
                props: {
                  size: "small",
                },
                style: {
                  height: "2rem",
                  width: "2rem",
                },
              },
              {
                props: {
                  size: "medium",
                },
                style: {
                  height: "2.5rem",
                  width: "2.5rem",
                },
              },
            ],
          },
        },
      },
      MuiInputBase: {
        styleOverrides: {
          root: {
            border: "none",
          },
        },
      },
      MuiLink: {
        defaultProps: {
          underline: "none",
          component: LinkBehavior,
        } as LinkProps,

        styleOverrides: {
          root: {
            color: brand[700],
            fontWeight: 500,
            position: "relative",
            textDecoration: "none",
            "&::before": {
              content: '""',
              position: "absolute",
              width: 0,
              height: "1px",
              bottom: 0,
              left: 0,
              backgroundColor: brand[200],
              opacity: 0.7,
              transition: "width 0.3s ease, opacity 0.3s ease",
            },
            "&:hover::before": {
              width: "100%",
              opacity: 1,
            },
            "&:focus-visible": {
              outline: `3px solid ${alpha(brand[500], 0.5)}`,
              outlineOffset: "4px",
              borderRadius: "2px",
            },
          },
        },
      },
      MuiMenuItem: {
        styleOverrides: {
          root: {
            paddingRight: 6,
            paddingLeft: 6,
            color: gray[700],
            fontSize: "0.875rem",
            fontWeight: 500,
            borderRadius: 6,
          },
        },
      },
      MuiOutlinedInput: {
        styleOverrides: {
          notchedOutline: {
            border: "none",
          },
          input: {
            paddingLeft: 10,
          },
          root: {
            "input:-webkit-autofill": {
              WebkitBoxShadow: `0 0 0 1000px ${brand[100]} inset, 0 0 0 1px ${brand[200]}`,
              maxHeight: "4px",
              borderRadius: "8px",
            },
            "& .MuiInputBase-input": {
              fontSize: "1rem",
              "&::placeholder": {
                opacity: 0.7,
                color: gray[500],
              },
            },
            boxSizing: "border-box",
            flexGrow: 1,
            height: "40px",
            borderRadius: 6,
            border: "1px solid",
            borderColor: alpha(gray[300], 0.8),
            boxShadow: "0 0 0 1.5px hsla(210, 0%, 0%, 0.02) inset",
            transition: "border-color 120ms ease-in",
            backgroundColor: alpha(gray[100], 0.4),
            "&:hover": {
              borderColor: brand[300],
            },
            "&.Mui-focused": {
              outline: `3px solid ${alpha(brand[500], 0.5)}`,
              outlineOffset: "2px",
              borderColor: brand[400],
            },
            variants: [
              {
                props: {
                  color: "error",
                },
                style: {
                  borderColor: red[200],
                  color: red[500],
                  "& + .MuiFormHelperText-root": {
                    color: red[500],
                  },
                },
              },
            ],
          },
        },
      },
      MuiPaper: {
        defaultProps: {
          elevation: 0,
        },
        styleOverrides: {
          root: {
            padding: "6px 12px",
            backgroundColor: "white",
            border: "1px solid",
            borderColor: alpha(gray[200], 0.8),
            boxShadow: "1px 1px 2px hsla(210, 0%, 0%, 0.1)",
          },
        },
      },
      MuiStack: {
        defaultProps: {
          useFlexGap: true,
        },
      },
      MuiSwitch: {
        styleOverrides: {
          root: {
            boxSizing: "border-box",
            width: 36,
            height: 24,
            padding: 0,
            transition: "background-color 100ms ease-in",
            "&:hover": {
              "& .MuiSwitch-track": {
                backgroundColor: brand[600],
              },
            },
            "& .MuiSwitch-switchBase": {
              "&.Mui-checked": {
                transform: "translateX(13px)",
              },
            },
            "& .MuiSwitch-track": {
              borderRadius: 50,
            },
            "& .MuiSwitch-thumb": {
              boxShadow: "0 0 2px 2px hsla(220, 0%, 0%, 0.2)",
              backgroundColor: "hsl(0, 0%, 100%)",
              width: 16,
              height: 16,
              margin: 2,
            },
          },
          switchBase: {
            height: 24,
            width: 24,
            padding: 0,
            color: "hsl(0, 0%, 100%)",
            "&.Mui-checked + .MuiSwitch-track": {
              opacity: 1,
            },
          },
        },
      },
      MuiToggleButtonGroup: {
        styleOverrides: {
          root: {
            borderRadius: 6,
            boxShadow: `0 1px 2px hsla(210, 0%, 0%, 0.05), 0 2px 12px ${alpha(brand[200], 0.5)}`,
            "& .Mui-selected": {
              color: brand[500],
            },
          },
        },
      },
      MuiToggleButton: {
        styleOverrides: {
          root: {
            padding: "12px 16px",
            textTransform: "none",
            borderRadius: 6,
            fontWeight: 500,
          },
        },
      },
    },
  });
}
