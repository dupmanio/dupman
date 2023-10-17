import { StatusOnWebsitesResponse } from "@/types/dtos/status";
import { UpdatesOnResponse } from "@/types/dtos/update";

type Website = {
  id: string;
  createdAt: string;
  updatedAt: string;
  url: string;
  status: StatusOnWebsitesResponse;
  updates: UpdatesOnResponse;
};

export type { Website };
