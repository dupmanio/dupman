import { useEffect, ReactElement } from "react";
import { useRouter } from "next/router";
import { signIn, useSession } from "next-auth/react";

function Login() {
  const router = useRouter();
  const { status } = useSession();

  const authProvider = "dupman";

  useEffect(() => {
    if (status === "authenticated" || status === "loading") {
      const callbackUrl = Array.isArray(router.query.callbackUrl)
        ? router.query.callbackUrl[0]
        : router.query.callbackUrl;
      void router.push(callbackUrl || "/");
    }

    if (status === "unauthenticated") {
      void signIn(authProvider);
    }
  }, [status, router]);

  return <></>;
}

Login.getLayout = function getLayout(children: ReactElement) {
  return children;
};

export default Login;
