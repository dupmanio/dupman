import { useState, MouseEvent } from "react";
import getNextConfig from "next/config";
import { signOut, useSession } from "next-auth/react";
import axios from "axios";

import {
  Avatar,
  Box,
  Divider,
  IconButton,
  ListItemIcon,
  Menu,
  MenuItem,
  Tooltip,
} from "@mui/material";

import { Logout } from "@mui/icons-material";

function AccountMenu() {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const { publicRuntimeConfig } = getNextConfig();
  const { data: session } = useSession();

  const open = Boolean(anchorEl);

  const logOut = async () => {
    try {
      await axios.put("/api/auth/logout");
      await signOut();

      return new Promise(() => {
        window.location.assign(publicRuntimeConfig.DUPMAN_LANDING);
      });
    } catch (error) {
      throw error;
    }
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  return (
    <>
      <Box sx={{ display: "flex", alignItems: "center", textAlign: "center" }}>
        <Tooltip title="Account settings">
          <IconButton
            onClick={(e: MouseEvent<HTMLElement>) =>
              setAnchorEl(e.currentTarget)
            }
            sx={{ ml: 2 }}
            aria-controls={open ? "account-menu" : undefined}
            aria-expanded={open ? "true" : undefined}
            aria-haspopup="true"
          >
            <Avatar sx={{ width: 32, height: 32 }}>
              {session?.user?.name?.charAt(0) ?? "D"}
            </Avatar>
          </IconButton>
        </Tooltip>
      </Box>

      <Menu
        anchorEl={anchorEl}
        id="account-menu"
        open={open}
        onClose={handleClose}
        onClick={handleClose}
        slotProps={{
          paper: {
            elevation: 0,
            sx: {
              overflow: "visible",
              filter: "drop-shadow(0px 2px 8px rgba(0,0,0,0.32))",
              mt: 1.5,
              "& .MuiAvatar-root": {
                width: 32,
                height: 32,
                ml: -0.5,
                mr: 1,
              },
              "&:before": {
                content: '""',
                display: "block",
                position: "absolute",
                top: 0,
                right: 14,
                width: 10,
                height: 10,
                bgcolor: "background.paper",
                transform: "translateY(-50%) rotate(45deg)",
                zIndex: 0,
              },
            },
          },
        }}
        transformOrigin={{ horizontal: "right", vertical: "top" }}
        anchorOrigin={{ horizontal: "right", vertical: "bottom" }}
      >
        <MenuItem
          onClick={handleClose}
          component={"a"}
          href={process.env.NEXT_PUBLIC_OIDC_ACCOUNT_PAGE}
        >
          <Avatar /> My account
        </MenuItem>

        <Divider />

        <MenuItem onClick={logOut}>
          <ListItemIcon>
            <Logout fontSize="small" />
          </ListItemIcon>
          Logout
        </MenuItem>
      </Menu>
    </>
  );
}

export default AccountMenu;
