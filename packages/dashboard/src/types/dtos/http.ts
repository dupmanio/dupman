type Pagination = {
  limit: number;
  page: number;
  totalItems: number;
  totalPages: number;
};

type HTTPResponse<T> = {
  code: number;
  data?: T;
  error?: unknown;
  pagination?: Pagination | null;
};
