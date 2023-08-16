import axios from "axios";
import getNextConfig from "next/config";
import { getSession } from "next-auth/react";

import { produceAccessTokenInterceptor } from "@/lib/http/client/interceptors";
import { CustomAxiosInstance } from "@/types/declarations/axios";

const { publicRuntimeConfig } = getNextConfig();

const DupmanAPIClient = axios.create({
  baseURL: publicRuntimeConfig.DUPMAN_API,
}) as unknown as CustomAxiosInstance;

let IS_REFRESHING = false;

DupmanAPIClient.reloadAuth = async () => {
  if (!IS_REFRESHING) {
    IS_REFRESHING = true;
    const session = await getSession();
    if (session && !session.error) {
      DupmanAPIClient.interceptors.request.handlers = [];
      DupmanAPIClient.interceptors.request.use(
        produceAccessTokenInterceptor(session.accessToken),
      );

      IS_REFRESHING = false;
      return;
    }

    IS_REFRESHING = false;
    throw new Error();
  }
};

export { DupmanAPIClient };
