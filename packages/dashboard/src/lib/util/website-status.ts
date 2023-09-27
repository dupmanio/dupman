import { useTheme } from "@mui/material/styles";

import { StatusOnWebsitesResponse, StatusState } from "@/types/dtos/status";

type StatusSettingsItem = {
  colorCode: string;
  color:
    | "default"
    | "primary"
    | "secondary"
    | "error"
    | "info"
    | "success"
    | "warning";
  title: string;
  additionalText?: string;
};

type StatusSettings = {
  [state in StatusState]: StatusSettingsItem;
};

function useStatusSettings(
  status: StatusOnWebsitesResponse,
): StatusSettingsItem {
  const theme = useTheme();

  const defaultSettings: StatusSettingsItem = {
    colorCode: "",
    color: "default",
    title: "",
  };
  const settings: StatusSettings = {
    [StatusState.NeedsUpdate]: {
      colorCode: theme.palette.warning.main,
      color: "warning",
      title: "Website Needs to be Updated",
    },
    [StatusState.UpToDated]: {
      colorCode: theme.palette.success.main,
      color: "success",
      title: "Website is Up To Date",
    },
    [StatusState.ScanningFailed]: {
      colorCode: theme.palette.error.main,
      color: "error",
      title: "Unable to scan website",
    },
  };

  if (!status.state) {
    return defaultSettings;
  }

  const currentSettings = settings[status.state];
  currentSettings.additionalText = status.info;

  return currentSettings;
}

export type { StatusSettings, StatusSettingsItem };
export { useStatusSettings };
