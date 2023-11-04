type WebsiteOnCreate = {
  url: string;
  token: string;
};

type WebsiteOnUpdate = {
  id: string;
  url: string;
  token: string;
};

export type { WebsiteOnCreate, WebsiteOnUpdate };
