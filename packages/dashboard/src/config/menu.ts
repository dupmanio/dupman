import { SvgIcon } from "@mui/material";

import LanguageIcon from "@mui/icons-material/Language";
import InfoIcon from "@mui/icons-material/Info";

type MenuItem = {
  route: string;
  icon: typeof SvgIcon;
  title: string;
};

const MenuItems: MenuItem[] = [
  {
    route: "/",
    icon: LanguageIcon,
    title: "Websites",
  },
  {
    route: "/about",
    icon: InfoIcon,
    title: "About",
  },
];

export default MenuItems;
