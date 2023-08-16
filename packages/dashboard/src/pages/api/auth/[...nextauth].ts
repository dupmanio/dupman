import NextAuth from "next-auth";
import { type NextAuthOptions } from "next-auth";
import axios from "axios";
import KeycloakProvider from "next-auth/providers/keycloak";
import { JWT } from "next-auth/jwt";

async function refreshAccessToken(token: JWT): Promise<JWT> {
  try {
    const { data: newToken, status } = await axios.post(
      `${process.env.OIDC_ISSUER}/protocol/openid-connect/token`,
      {
        grant_type: "refresh_token",
        client_id: process.env.OIDC_CLIENT_ID ?? "",
        client_secret: process.env.OIDC_CLIENT_SECRET ?? "",
        refresh_token: token.refreshToken,
      },
      {
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
      },
    );

    if (status !== 200) {
      return {
        ...token,
        error: "RefreshAccessTokenError",
      };
    }

    return {
      ...token,
      accessToken: newToken.access_token,
      accessTokenExpires: Date.now() + newToken.expires_in * 1000,
      refreshToken: newToken.refresh_token,
      idToken: newToken.id_token,
    };
  } catch (error) {
    return {
      ...token,
      error: "RefreshAccessTokenError",
    };
  }
}

export const authOptions: NextAuthOptions = {
  providers: [
    KeycloakProvider({
      issuer: process.env.OIDC_ISSUER,
      clientId: process.env.OIDC_CLIENT_ID ?? "",
      clientSecret: process.env.OIDC_CLIENT_SECRET ?? "",
      authorization: { params: { scope: "openid email profile website" } },
    }),
  ],
  pages: {
    signIn: "/login",
    error: "/error",
  },
  session: {
    maxAge: 7 * 24 * 60 * 60, // 7 days.
  },
  callbacks: {
    jwt: async ({ token, account }) => {
      // Initial sign in.
      if (
        account &&
        account.access_token &&
        account.refresh_token &&
        account.id_token &&
        account.expires_at
      ) {
        token.accessToken = account.access_token;
        token.refreshToken = account.refresh_token;
        token.idToken = account.id_token;
        token.accessTokenExpires = account.expires_at * 1000;

        return token;
      }

      // Return previous token if the access token has not expired yet.
      if (Date.now() < (token.accessTokenExpires as number)) {
        return token;
      }

      // Access token has expired, try to update it.
      return refreshAccessToken(token);
    },
    session: async ({ session, token }) => {
      session.accessToken = token.accessToken as string;
      session.refreshToken = token.refreshToken as string;
      session.idToken = token.idToken as string;
      session.error = !!token?.error;

      return session;
    },
    redirect: async ({ url, baseUrl }) => {
      // Redirect to Keycloak logout page.
      if (url.startsWith("/logout")) {
        const url = new URL(
          `${process.env.OIDC_ISSUER}/protocol/openid-connect/logout`,
        );

        url.searchParams.append(
          "post_logout_redirect_uri",
          process.env.DUPMAN_LANDING ?? baseUrl,
        );
        url.searchParams.append("client_id", process.env.OIDC_CLIENT_ID ?? "");

        return url.toString();
      }

      if (url.startsWith(baseUrl)) {
        return url;
      }

      if (url.startsWith("/")) {
        return `${baseUrl}${url}`;
      }

      return baseUrl;
    },
  },
};

export default NextAuth(authOptions);
