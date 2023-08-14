import NextAuth from "next-auth";
import { JWT } from "next-auth/src/jwt";
import { AdapterUser } from "next-auth/src/adapters";
import { Account, Profile, Session, User } from "next-auth/src/core/types";
import KeycloakProvider from "next-auth/providers/keycloak";

type JWTParams = {
  token: JWT;
  user: User | AdapterUser;
  account: Account | null;
  profile?: Profile;
};

type SessionParams = {
  session: Session;
  token: JWT;
  user: AdapterUser;
};

async function refreshAccessToken(token: JWT): Promise<JWT> {
  try {
    const params = new URLSearchParams();
    params.append("grant_type", "refresh_token");
    params.append("client_id", process.env.OIDC_CLIENT_ID ?? "");
    params.append("client_secret", process.env.OIDC_CLIENT_SECRET ?? "");
    params.append("refresh_token", token.refreshToken);

    const response = await fetch(
      `${process.env.OIDC_ISSUER}/protocol/openid-connect/token`,
      {
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
        method: "POST",
        body: params,
      },
    );

    if (!response.ok) {
      return {
        ...token,
        error: "RefreshAccessTokenError",
      };
    }

    const newToken = await response.json();

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

export const authOptions = {
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
    jwt: async function ({ token, account }: JWTParams) {
      // Initial sign in/
      if (
        account &&
        account.access_token &&
        account.refresh_token &&
        account.id_token
      ) {
        token.accessToken = account.access_token;
        token.refreshToken = account.refresh_token;
        token.idToken = account.id_token;
        token.accessTokenExpires = account.expires_at * 1000;

        return token;
      }

      if (Date.now() < (token.accessTokenExpires as number)) {
        return token;
      }

      return refreshAccessToken(token);
    },
    session: async function ({ session, token }: SessionParams) {
      session.accessToken = token.accessToken as string;
      session.refreshToken = token.refreshToken as string;
      session.idToken = token.idToken as string;
      session.error = !!token?.error;

      return session;
    },
  },
};

export default NextAuth(authOptions);
