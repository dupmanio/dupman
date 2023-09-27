import { StatusOnWebsitesResponse } from "@/types/dtos/status";

type Website = {
  id: string;
  createdAt: string;
  updatedAt: string;
  url: string;
  status: StatusOnWebsitesResponse;
};

export type { Website };
