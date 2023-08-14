import { AxiosInstance } from "axios";

declare module "axios" {
  class CustomAxiosInstance extends AxiosInstance {
    interceptors: {
      request: {
        handlers?: InternalAxiosRequestConfig<unknown>[];
      } & import("axios").AxiosInterceptorManager<
        import("axios").InternalAxiosRequestConfig<unknown>
      >;
      response: import("axios").AxiosInterceptorManager<
        import("axios").AxiosResponse<unknown, unknown>
      >;
    };

    reloadAuth(): Promise<void>;
  }
}
