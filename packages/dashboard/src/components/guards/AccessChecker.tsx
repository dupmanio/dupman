import * as React from "react";
import { ReactNode } from "react";
import AuthGuard from "@/components/guards/Auth";
import AnonymousGuard from "@/components/guards/Anonymous";
import { PageAccess } from "@/config/page-accesss";

type AccessCheckerProps = {
  children: ReactNode;
  access?: PageAccess;
};

function AccessChecker({ children, access }: AccessCheckerProps) {
  if (!access) {
    return <h1>Page not found!</h1>;
  }

  return (
    <>
      {access === PageAccess.SECURED && <AuthGuard>{children}</AuthGuard>}
      {access === PageAccess.ANONYMOUS && (
        <AnonymousGuard>{children}</AnonymousGuard>
      )}
      {access == PageAccess.SHARED && children}
    </>
  );
}

export default AccessChecker;
