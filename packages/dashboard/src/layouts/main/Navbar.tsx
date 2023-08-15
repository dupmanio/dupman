import * as React from "react";
import Link from "next/link";
import { useRouter } from "next/router";

import {
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
} from "@mui/material";

import MenuItems from "@/config/menu";

export default function Navbar() {
  const router = useRouter();

  return (
    <List component="nav">
      {MenuItems.map((data, idx) => (
        <ListItemButton
          key={idx}
          selected={router.pathname === data.route}
          component={Link}
          href={data.route}
          passHref
        >
          <ListItemIcon>
            <data.icon />
          </ListItemIcon>
          <ListItemText primary={data.title} />
        </ListItemButton>
      ))}
    </List>
  );
}
