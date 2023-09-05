import { SvgIcon } from "@mui/material";

import LanguageIcon from "@mui/icons-material/Language";
import InfoIcon from "@mui/icons-material/Info";

import { Route } from "@/config/routes";

type MenuItem = {
  route: Route;
  icon: typeof SvgIcon;
  title: string;
};

const MenuItems: MenuItem[] = [
  {
    route: Route.HOMEPAGE,
    icon: LanguageIcon,
    title: "Websites",
  },
  {
    route: Route.ABOUT,
    icon: InfoIcon,
    title: "About",
  },
];

export default MenuItems;
