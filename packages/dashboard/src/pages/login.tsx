import { signIn, useSession } from "next-auth/react";
import { useEffect } from "react";
import { useRouter } from "next/router";
import { PageAccess } from "@/config/page-accesss";

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

Login.Access = PageAccess.ANONYMOUS;

export default Login;
