import React, { ReactNode, useState } from "react";
import Image from "next/image";

import {
  Box,
  Toolbar,
  Typography,
  Divider,
  IconButton,
  Container,
  Grid,
  Link,
} from "@mui/material";

import MenuIcon from "@mui/icons-material/Menu";
import ChevronLeftIcon from "@mui/icons-material/ChevronLeft";

import AccountMenu from "@/layouts/main/AccountMenu";
import Navbar from "@/layouts/main/Navbar";
import Copyright from "@/layouts/main/Copyright";
import StyledAppBar from "@/layouts/main/StyledAppBar";
import StyledDrawer from "@/layouts/main/StyledDrawer";
import Notifications from "@/layouts/main/Notifications";

interface MainLayoutProps {
  children: ReactNode;
}

function MainLayout({ children }: MainLayoutProps) {
  const [drawerOpen, setDrawerOpen] = useState<boolean>(true);

  return (
    <>
      <Box sx={{ display: "flex" }}>
        <StyledAppBar position="absolute" drawerWidth={240} open={drawerOpen}>
          <Toolbar sx={{ pr: "24px" }}>
            <IconButton
              edge="start"
              color="inherit"
              aria-label="open drawer"
              onClick={() => setDrawerOpen(!drawerOpen)}
              sx={{
                marginRight: "36px",
                ...(drawerOpen && { display: "none" }),
              }}
            >
              <MenuIcon />
            </IconButton>
            <Typography
              component="h1"
              variant="h6"
              color="inherit"
              noWrap
              sx={{ flexGrow: 1 }}
            >
              Dashboard
            </Typography>

            <Notifications />

            <AccountMenu />
          </Toolbar>
        </StyledAppBar>
        <StyledDrawer variant="permanent" width={240} open={drawerOpen}>
          <Toolbar
            sx={{
              display: "flex",
              alignItems: "center",
              justifyContent: "flex-end",
              px: [1],
            }}
          >
            <Link>
              <Image
                src="/assets/dupman.png"
                alt="dupman"
                priority={true}
                width={160}
                height={40}
              />
            </Link>

            <IconButton onClick={() => setDrawerOpen(!drawerOpen)}>
              <ChevronLeftIcon />
            </IconButton>
          </Toolbar>
          <Divider />
          <Navbar />
        </StyledDrawer>
        <Box
          component="main"
          sx={{
            flexGrow: 1,
            height: "100vh",
            overflow: "auto",
          }}
        >
          <Toolbar />
          <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
            <Grid container spacing={3}>
              {children}
            </Grid>
            <Copyright />
          </Container>
        </Box>
      </Box>
    </>
  );
}

export default MainLayout;
