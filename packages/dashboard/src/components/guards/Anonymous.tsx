import { ReactNode, useEffect } from "react";
import { useSession } from "next-auth/react";
import { useRouter } from "next/router";

import PageLoader from "@/components/PageLoader";

type IProps = {
  children: ReactNode;
};

function AnonymousGuard({ children }: IProps) {
  const { data, status } = useSession();
  const router = useRouter();
  const hasNoAccess = !!data?.user;

  useEffect(() => {
    if (status === "loading") {
      return;
    }

    if (hasNoAccess) {
      void router.push("/");
    }
  }, [status, hasNoAccess, router]);

  if (!hasNoAccess) {
    return <>{children}</>;
  }

  return <PageLoader sx={{ width: "100%", height: "100vh" }} />;
}

export default AnonymousGuard;
