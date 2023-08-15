import { DupmanAPIClient } from "@/lib/http/client/dupman-api";
import { AxiosInstance } from "axios";
import { Website } from "@/types/entities/website";
import { WebsiteOnCreate } from "@/types/dtos/website";
import { HTTPResponse } from "@/types/dtos/http";

interface IWebsiteRepository {
  getAll: (page: number, limit: number) => Promise<HTTPResponse<Website[]>>;
  create: (payload: WebsiteOnCreate) => Promise<HTTPResponse<Website>>;
}

function UseRepositoryFactory(http: AxiosInstance): IWebsiteRepository {
  return {
    create: async (payload: WebsiteOnCreate) => {
      const response = await http.post("/website/", payload, {
        headers: {
          "Content-Type": "application/json",
        },
      });
      return response.data;
    },
    getAll: async (page, limit) => {
      const response = await http.get("/website/", { params: { page, limit } });
      return response.data;
    },
  };
}

const WebsiteRepository = UseRepositoryFactory(DupmanAPIClient);

export { UseRepositoryFactory, WebsiteRepository };
