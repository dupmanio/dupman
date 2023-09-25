enum StatusState {
  UpToDated = "UP_TO_DATED",
  NeedsUpdate = "NEEDS_UPDATE",
  ScanningFailed = "SCANNING_FAILED",
}

type StatusOnWebsitesResponse = {
  state: StatusState;
  info: string;
  updatedAt: string;
};

export { StatusState };
export type { StatusOnWebsitesResponse };
