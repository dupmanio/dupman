import { AxiosInstance } from "axios";

import { DupmanAPIClient } from "@/lib/http/client/dupman-api";
import { Website } from "@/types/entities/website";
import { WebsiteOnCreate } from "@/types/dtos/website";
import { HTTPResponse } from "@/types/dtos/http";

interface IWebsiteRepository {
  getSingle: (id: string) => Promise<HTTPResponse<Website>>;
  getAll: (page: number, limit: number) => Promise<HTTPResponse<Website[]>>;
  create: (payload: WebsiteOnCreate) => Promise<HTTPResponse<Website>>;
  delete: (id: string) => Promise<HTTPResponse<null>>;
}

const servicePrefix = "/api";

function UseRepositoryFactory(http: AxiosInstance): IWebsiteRepository {
  return {
    getSingle: async (id) => {
      const response = await http.get(`${servicePrefix}/website/${id}`);
      return response.data;
    },
    create: async (payload: WebsiteOnCreate) => {
      const response = await http.post(`${servicePrefix}/website/`, payload, {
        headers: {
          "Content-Type": "application/json",
        },
      });
      return response.data;
    },
    getAll: async (page, limit) => {
      const response = await http.get(`${servicePrefix}/website/`, {
        params: { page, limit },
      });
      return response.data;
    },
    delete: async (id) => {
      const response = await http.delete(`${servicePrefix}/website/${id}`);
      return response.data;
    },
  };
}

const WebsiteRepository = UseRepositoryFactory(DupmanAPIClient);

export { UseRepositoryFactory, WebsiteRepository };
