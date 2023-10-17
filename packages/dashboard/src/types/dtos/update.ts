import { v4 as uuid } from "uuid";

type Update = {
  name: string;
  title: string;
  link: string;
  type: string;
  currentVersion: string;
  latestVersion: string;
  recommendedVersion: string;
  installType: string;
  status: number;
};

type Updates = Update[];

interface UpdateOnResponse extends Update {
  id: typeof uuid;
  createdAt: Date;
  updatedAt: Date;
}

type UpdatesOnResponse = UpdateOnResponse[];

export type { Update, Updates, UpdateOnResponse, UpdatesOnResponse };
