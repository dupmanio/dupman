import { authOptions } from "./[...nextauth]";
import { getServerSession } from "next-auth/next";
import axios from "axios";
import { NextApiRequest, NextApiResponse } from "next";

export default async function signOut(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  if (req.method === "PUT") {
    const session = await getServerSession(req, res, authOptions);

    if (session?.idToken) {
      try {
        await axios.get(
          `${process.env.OIDC_ISSUER}/protocol/openid-connect/logout`,
          { params: { id_token_hint: session.idToken } },
        );

        res.status(200).json({ message: "Logged-out successfully" });

        return;
      } catch (error) {
        res.status(500).json({ error: "Unable to log-out from provide" });

        return;
      }
    }

    res.status(200).json({});
  }

  res.status(405).json({ error: "Method Not Allowed" });
}
