import axios from "axios";
import getNextConfig from "next/config";
import { getSession } from "next-auth/react";

import { produceAccessTokenInterceptor } from "@/lib/http/client/interceptors";
import { CustomAxiosInstance } from "@/types/declarations/axios";

const { publicRuntimeConfig } = getNextConfig();

const NotifyClient = axios.create({
  baseURL: publicRuntimeConfig.NOTIFY_URL,
}) as unknown as CustomAxiosInstance;

let IS_REFRESHING = false;

NotifyClient.reloadAuth = async () => {
  if (!IS_REFRESHING) {
    IS_REFRESHING = true;
    const session = await getSession();
    if (session && !session.error) {
      NotifyClient.interceptors.request.handlers = [];
      NotifyClient.interceptors.request.use(
        produceAccessTokenInterceptor(session.accessToken),
      );

      IS_REFRESHING = false;
      return;
    }

    IS_REFRESHING = false;
    throw new Error();
  }
};

export { NotifyClient };
