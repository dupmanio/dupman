import { Roboto } from "next/font/google";
import { createTheme } from "@mui/material/styles";
import { red } from "@mui/material/colors";

export const roboto = Roboto({
  weight: ["300", "400", "500", "700"],
  subsets: ["latin"],
  display: "swap",
});

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
