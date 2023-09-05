import * as React from "react";
import type { AppProps } from "next/app";
import { NextComponentType } from "next/types";
import Head from "next/head";
import type { Session } from "next-auth";
import { SessionProvider } from "next-auth/react";

import { SnackbarProvider } from "notistack";

import { ThemeProvider, CssBaseline } from "@mui/material";
import AuthGuard from "@/components/guards/Auth";
import MainLayout from "@/layouts/main";
import theme from "@/themes/main";

export interface MyAppProps extends AppProps {
  session: Session;
  Component: NextComponentType & {
    getLayout: (children: React.ReactElement) => React.JSX.Element;
  };
}

export default function MyApp({
  Component,
  pageProps: { session, ...pageProps },
}: MyAppProps) {
  const getLayout =
    Component.getLayout ?? ((children) => <MainLayout>{children}</MainLayout>);

  return (
    <SessionProvider session={session}>
      <AuthGuard>
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
            {getLayout(<Component {...pageProps} />)}
          </SnackbarProvider>
        </ThemeProvider>
      </AuthGuard>
    </SessionProvider>
  );
}
