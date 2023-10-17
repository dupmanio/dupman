import { AxiosInstance } from "axios";

import { PreviewAPIClient } from "@/lib/http/client/preview-api";
import { HTTPResponse } from "@/types/dtos/http";
import { Preview } from "@/types/dtos/preview";

interface IPreviewRepository {
  get: (websiteId: string) => Promise<HTTPResponse<Preview>>;
}

function UseRepositoryFactory(http: AxiosInstance): IPreviewRepository {
  return {
    get: async (websiteId) => {
      const response = await http.get(`/preview/${websiteId}`);
      return response.data;
    },
  };
}

const PreviewRepository = UseRepositoryFactory(PreviewAPIClient);

export { UseRepositoryFactory, PreviewRepository };
