import { ReactNode, useCallback, useEffect, useRef, useState } from "react";
import { useRouter } from "next/router";
import { signOut, useSession } from "next-auth/react";
import { AxiosRequestConfig } from "axios";

import {
  produceAccessTokenInterceptor,
  produceLogoutInterceptor,
} from "@/lib/http/client/interceptors";
import { DupmanAPIClient } from "@/lib/http/client/dupman-api";
import { PreviewAPIClient } from "@/lib/http/client/preview-api";
import PageLoader from "@/components/PageLoader";
import { Route } from "@/config/routes";

interface AuthGuardProps {
  children: ReactNode;
}

function AuthGuard({ children }: AuthGuardProps) {
  const [interceptor, setInterceptor] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(true);
  const dupmanAPIRequestInterceptorRef = useRef<number>(0);
  const dupmanAPIResponseInterceptorRef = useRef<number>(0);
  const previewAPIRequestInterceptorRef = useRef<number>(0);
  const previewAPIResponseInterceptorRef = useRef<number>(0);

  const { data, status } = useSession();
  const router = useRouter();
  const logOutCallback = useCallback(async () => {
    await signOut({ redirect: false });
    void router.push(Route.LOGIN);
  }, [router]);

  useEffect(() => {
    if (status === "authenticated") {
      dupmanAPIRequestInterceptorRef.current =
        DupmanAPIClient.interceptors.request.use(
          produceAccessTokenInterceptor(data?.accessToken),
        );
      previewAPIRequestInterceptorRef.current =
        PreviewAPIClient.interceptors.request.use(
          produceAccessTokenInterceptor(data?.accessToken),
        );

      setInterceptor(true);
      setLoading(false);
    }

    if (status === "unauthenticated") {
      setLoading(false);
    }

    return () => {
      DupmanAPIClient.interceptors.request.eject(
        dupmanAPIRequestInterceptorRef.current,
      );
      PreviewAPIClient.interceptors.request.eject(
        previewAPIRequestInterceptorRef.current,
      );
    };
  }, [data, status]);

  useEffect(() => {
    if (interceptor) {
      dupmanAPIResponseInterceptorRef.current =
        DupmanAPIClient.interceptors.response.use(
          (req: AxiosRequestConfig) => req,
          produceLogoutInterceptor(logOutCallback),
        );
      previewAPIResponseInterceptorRef.current =
        PreviewAPIClient.interceptors.response.use(
          (req: AxiosRequestConfig) => req,
          produceLogoutInterceptor(logOutCallback),
        );
    }

    return () => {
      DupmanAPIClient.interceptors.response.eject(
        dupmanAPIResponseInterceptorRef.current,
      );
      PreviewAPIClient.interceptors.response.eject(
        previewAPIResponseInterceptorRef.current,
      );
    };
  }, [interceptor, logOutCallback, router]);

  if (loading) {
    return <PageLoader sx={{ width: "100%", height: "100vh" }} />;
  }

  return children;
}

export default AuthGuard;
