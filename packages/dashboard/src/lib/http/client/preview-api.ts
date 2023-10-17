import axios from "axios";
import getNextConfig from "next/config";
import { getSession } from "next-auth/react";

import { produceAccessTokenInterceptor } from "@/lib/http/client/interceptors";
import { CustomAxiosInstance } from "@/types/declarations/axios";

const { publicRuntimeConfig } = getNextConfig();

const PreviewAPIClient = axios.create({
  baseURL: publicRuntimeConfig.PREVIEW_API,
}) as unknown as CustomAxiosInstance;

let IS_REFRESHING = false;

PreviewAPIClient.reloadAuth = async () => {
  if (!IS_REFRESHING) {
    IS_REFRESHING = true;
    const session = await getSession();
    if (session && !session.error) {
      PreviewAPIClient.interceptors.request.handlers = [];
      PreviewAPIClient.interceptors.request.use(
        produceAccessTokenInterceptor(session.accessToken),
      );

      IS_REFRESHING = false;
      return;
    }

    IS_REFRESHING = false;
    throw new Error();
  }
};

export { PreviewAPIClient };
