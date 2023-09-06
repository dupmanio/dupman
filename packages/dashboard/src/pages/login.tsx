import { useEffect, ReactElement } from "react";
import { useRouter } from "next/router";
import { signIn, useSession } from "next-auth/react";

function Login() {
  const router = useRouter();
  const { status } = useSession();

  const authProvider = "keycloak";

  useEffect(() => {
    if (status === "authenticated" || status === "loading") {
      void router.push("/");
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
