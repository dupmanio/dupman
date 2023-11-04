import { AxiosInstance } from "axios";

import { DupmanAPIClient } from "@/lib/http/client/dupman-api";
import { Website } from "@/types/entities/website";
import { WebsiteOnCreate, WebsiteOnUpdate } from "@/types/dtos/website";
import { HTTPResponse } from "@/types/dtos/http";

interface IWebsiteRepository {
  create: (payload: WebsiteOnCreate) => Promise<HTTPResponse<Website>>;
  getSingle: (id: string) => Promise<HTTPResponse<Website>>;
  getAll: (page: number, limit: number) => Promise<HTTPResponse<Website[]>>;
  update: (payload: WebsiteOnUpdate) => Promise<HTTPResponse<null>>;
  delete: (id: string) => Promise<HTTPResponse<Website>>;
}

const servicePrefix = "/api";

function UseRepositoryFactory(http: AxiosInstance): IWebsiteRepository {
  return {
    create: async (payload: WebsiteOnCreate) => {
      const response = await http.post(`${servicePrefix}/website/`, payload, {
        headers: {
          "Content-Type": "application/json",
        },
      });
      return response.data;
    },
    getSingle: async (id) => {
      const response = await http.get(`${servicePrefix}/website/${id}`);
      return response.data;
    },
    getAll: async (page, limit) => {
      const response = await http.get(`${servicePrefix}/website/`, {
        params: { page, limit },
      });
      return response.data;
    },
    update: async (payload: WebsiteOnUpdate) => {
      const response = await http.patch(`${servicePrefix}/website/`, payload, {
        headers: {
          "Content-Type": "application/json",
        },
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
