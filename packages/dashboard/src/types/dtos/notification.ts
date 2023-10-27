import { v4 as uuid } from "uuid";

type NotificationOnResponse = {
  id: typeof uuid;
  createdAt: string;
  userID: typeof uuid;
  type: string;
  title: string;
  message: string;
  seen: boolean;
};

export type { NotificationOnResponse };
