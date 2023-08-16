import * as React from "react";
import type { AppProps } from "next/app";
import Head from "next/head";
import type { Session } from "next-auth";
import { SessionProvider } from "next-auth/react";
import { SnackbarProvider } from "notistack";

import { ThemeProvider, CssBaseline } from "@mui/material";

import AccessChecker from "@/components/guards/AccessChecker";
import { PageAccess } from "@/config/page-accesss";
import theme from "@/themes/main";

export interface MyAppProps extends AppProps {
  session: Session;
}

export default function MyApp({
  Component,
  pageProps: { session, ...pageProps },
}: MyAppProps & {
  Component: { Access: PageAccess };
}) {
  return (
    <SessionProvider session={session}>
      <AccessChecker access={Component.Access}>
        <Head>
          <title>dupman</title>
          <meta name="viewport" content="initial-scale=1, width=device-width" />
        </Head>

        <ThemeProvider theme={theme}>
          <SnackbarProvider
            anchorOrigin={{
              vertical: "top",
              horizontal: "right",
            }}
          >
            <CssBaseline />
            <Component {...pageProps} />
          </SnackbarProvider>
        </ThemeProvider>
      </AccessChecker>
    </SessionProvider>
  );
}
