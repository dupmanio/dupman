import { Roboto } from "next/font/google";
import { createTheme } from "@mui/material/styles";
import { red } from "@mui/material/colors";

export const roboto = Roboto({
  weight: ["300", "400", "500", "700"],
  subsets: ["latin"],
  display: "swap",
});

//   colors: {
//     background: "#EDEFF4",
//
//     "primary-lighten-5": "#f0f9ff",
//     "primary-lighten-4": "#c2e7fd",
//     "primary-lighten-3": "#85cefb",
//     "primary-lighten-2": "#48b6f9",
//     "primary-lighten-1": "#0b9ef7",
//     primary: "#0678be",
//     "primary-darken-1": "#056098",
//     "primary-darken-2": "#044872",
//     "primary-darken-3": "#02304c",
//     "primary-darken-4": "#000609",
//
//     "secondary-lighten-5": "#f7fbfe",
//     "secondary-lighten-4": "#deeff9",
//     "secondary-lighten-3": "#bddef4",
//     "secondary-lighten-2": "#9cceee",
//     "secondary-lighten-1": "#7bbde9",
//     secondary: "#5aade3",
//     "secondary-darken-1": "#2592d9",
//     "secondary-darken-2": "#1c6da3",
//     "secondary-darken-3": "#12496c",
//     "secondary-darken-4": "#02090e",
//   },

const index = createTheme({
  palette: {
    background: {
      default: "#edeff4",
    },
    primary: {
      main: "#0678be",
      light: "#5aa3e3",
      dark: "#005194",
      contrastText: "#ffffff",
    },
    secondary: {
      main: "#5aade3",
      light: "#8cd0ff",
      dark: "#0078b2",
      contrastText: "#000000",
    },
    error: {
      main: red.A400,
    },
  },
  typography: {
    fontFamily: roboto.style.fontFamily,
  },
});

export default index;
