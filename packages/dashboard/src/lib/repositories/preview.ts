import { AxiosInstance } from "axios";

import { DupmanAPIClient } from "@/lib/http/client/dupman-api";
import { HTTPResponse } from "@/types/dtos/http";
import { Preview } from "@/types/dtos/preview";

interface IPreviewRepository {
  get: (websiteId: string) => Promise<HTTPResponse<Preview>>;
}

const servicePrefix = "/preview-api";

function UseRepositoryFactory(http: AxiosInstance): IPreviewRepository {
  return {
    get: async (websiteId) => {
      const response = await http.get(`${servicePrefix}/preview/${websiteId}`);
      return response.data;
    },
  };
}

const PreviewRepository = UseRepositoryFactory(DupmanAPIClient);

export { UseRepositoryFactory, PreviewRepository };
