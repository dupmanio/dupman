import { ReactNode, useCallback, useEffect, useRef, useState } from "react";
import { useRouter } from "next/router";
import { signOut, useSession } from "next-auth/react";

import {
  produceAccessTokenInterceptor,
  produceLogoutInterceptor,
} from "@/lib/http/client/interceptors";
import { DupmanAPIClient } from "@/lib/http/client/dupman-api";
import PageLoader from "@/components/PageLoader";
import Layout from "@/layouts/main";

type IProps = {
  children: ReactNode;
};

function AuthGuard({ children }: IProps) {
  const [interceptor, setInterceptor] = useState(false);
  const axiosClientRef = useRef<number>(0);
  const axiosClientReqRef = useRef<number>(0);

  const { data, status } = useSession();
  const router = useRouter();
  const logOutCallback = useCallback(async () => {
    await signOut({ redirect: false });
    void router.push("/login");
  }, [router]);
  const hasAccess = !!data?.user;

  useEffect(() => {
    if (status === "loading") {
      return;
    }

    if (!hasAccess) {
      void router.push("/login");
    }
  }, [status, hasAccess, router]);

  useEffect(() => {
    if (status === "authenticated") {
      axiosClientReqRef.current = DupmanAPIClient.interceptors.request.use(
        produceAccessTokenInterceptor(data?.accessToken),
      );

      setInterceptor(true);
    }

    return () => {
      DupmanAPIClient.interceptors.request.eject(axiosClientReqRef.current);
    };
  }, [data, status]);

  useEffect(() => {
    if (interceptor) {
      axiosClientRef.current = DupmanAPIClient.interceptors.response.use(
        (req) => req,
        produceLogoutInterceptor(logOutCallback),
      );
    }

    return () => {
      DupmanAPIClient.interceptors.response.eject(axiosClientRef.current);
    };
  }, [interceptor, logOutCallback, router]);

  if (hasAccess && interceptor) {
    return <Layout>{children}</Layout>;
  }

  return <PageLoader sx={{ width: "100%", height: "100vh" }} />;
}

export default AuthGuard;
