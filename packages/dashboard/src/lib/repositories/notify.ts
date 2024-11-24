import { AxiosInstance } from "axios";
import { v4 as uuid } from "uuid";

import { HTTPResponse } from "@/types/dtos/http";
import { DupmanAPIClient } from "@/lib/http/client/dupman-api";
import { NotificationOnResponse } from "@/types/dtos/notification";

interface INotifyRepository {
  getCount: () => Promise<HTTPResponse<number>>;
  getAll: (
    page: number,
    limit: number,
  ) => Promise<HTTPResponse<NotificationOnResponse[]>>;
  markAsRead: (id: typeof uuid) => Promise<HTTPResponse<null>>;
  markAllAsRead: () => Promise<HTTPResponse<null>>;
  deleteAll: () => Promise<HTTPResponse<null>>;
}

const servicePrefix = "/notify";

function UseRepositoryFactory(http: AxiosInstance): INotifyRepository {
  return {
    getCount: async () => {
      const response = await http.get(`${servicePrefix}/notification/count`);
      return response.data;
    },
    getAll: async (page, limit) => {
      const response = await http.get(`${servicePrefix}/notification`, {
        params: { page, limit },
      });
      return response.data;
    },
    markAsRead: async (id: typeof uuid) => {
      const response = await http.post(
        `${servicePrefix}/notification/${id}/mark-as-read`,
      );
      return response.data;
    },
    markAllAsRead: async () => {
      const response = await http.post(
        `${servicePrefix}/notification/mark-all-as-read`,
      );
      return response.data;
    },
    deleteAll: async () => {
      const response = await http.delete(`${servicePrefix}/notification`);
      return response.data;
    },
  };
}

const NotifyRepository = UseRepositoryFactory(DupmanAPIClient);

export { UseRepositoryFactory, NotifyRepository };
