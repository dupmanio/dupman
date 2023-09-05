export { default } from "next-auth/middleware";

export const config = {
  // @todo: Investigate why this does not work with enum values.
  matcher: ["/", "/about"],
};
