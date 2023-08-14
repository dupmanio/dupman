import * as React from "react";
import { ReactElement, Fragment } from "react";
import {
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
} from "@mui/material";
import LanguageIcon from "@mui/icons-material/Language";
import InfoIcon from "@mui/icons-material/Info";
import Link from "next/link";
import { useRouter } from "next/router";

type IMenuItem = {
  route: string;
  icon: ReactElement;
  title: string;
};

export default function Navbar() {
  const router = useRouter();

  const menu: IMenuItem[] = [
    {
      route: "/",
      icon: <LanguageIcon />,
      title: "Websites",
    },
    {
      route: "/about",
      icon: <InfoIcon />,
      title: "About",
    },
  ];

  return (
    <List component="nav">
      {menu.map((data, idx) => (
        <Fragment key={idx}>
          <ListItemButton
            selected={router.pathname === data.route}
            component={Link}
            href={data.route}
            passHref
          >
            <ListItemIcon>{data.icon}</ListItemIcon>
            <ListItemText primary={data.title} />
          </ListItemButton>
        </Fragment>
      ))}
    </List>
  );
}
