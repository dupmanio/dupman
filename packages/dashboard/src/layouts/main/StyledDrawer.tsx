import { Drawer, DrawerProps, styled } from "@mui/material";

interface StyledDrawerProps extends DrawerProps {
  open?: boolean;
  width: number;
}

const StyledDrawer = styled(Drawer, {
  shouldForwardProp: (prop) => prop !== "open" && prop !== "width",
})<StyledDrawerProps>(({ theme, open, width }) => ({
  "& .MuiDrawer-paper": {
    position: "relative",
    whiteSpace: "nowrap",
    width: width,
    transition: theme.transitions.create("width", {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
    boxSizing: "border-box",
    ...(!open && {
      overflowX: "hidden",
      transition: theme.transitions.create("width", {
        easing: theme.transitions.easing.sharp,
        duration: theme.transitions.duration.leavingScreen,
      }),
      width: theme.spacing(7),
      [theme.breakpoints.up("sm")]: {
        width: theme.spacing(9),
      },
    }),
  },
}));

export default StyledDrawer;
