import * as React from "react";

import {
  IconButton,
  Tooltip,
  TooltipProps,
  tooltipClasses,
  styled,
} from "@mui/material";
import { useTheme } from "@mui/material/styles";

import FiberManualRecordIcon from "@mui/icons-material/FiberManualRecord";

import { StatusState, StatusOnWebsitesResponse } from "@/types/dtos/status";

export interface WebsiteStatusCellProps {
  status: StatusOnWebsitesResponse;
}

type StatusSettings = {
  [state in StatusState]: {
    color: string;
    tooltipText: string;
  };
};

const HtmlTooltip = styled(({ className, ...props }: TooltipProps) => (
  <Tooltip {...props} classes={{ popper: className }} />
))(({ theme }) => ({
  [`& .${tooltipClasses.tooltip}`]: {
    backgroundColor: theme.palette.background.default,
    color: theme.palette.primary.main,
    maxWidth: 220,
    fontSize: theme.typography.pxToRem(12),
    border: "1px solid #dadde9",
  },
}));

function WebsiteStatusCell({ status }: WebsiteStatusCellProps) {
  const theme = useTheme();

  const statusSettings: StatusSettings = {
    [StatusState.NeedsUpdate]: {
      color: theme.palette.warning.main,
      tooltipText: "Website Needs to be Updated",
    },
    [StatusState.UpToDated]: {
      color: theme.palette.success.main,
      tooltipText: "Website is Up To Date",
    },
    [StatusState.ScanningFailed]: {
      color: theme.palette.error.main,
      tooltipText: `Scanning Failed: ${status.info}`,
    },
  };

  return (
    <HtmlTooltip
      title={statusSettings[status.state].tooltipText ?? ""}
      leaveDelay={100}
    >
      <IconButton>
        <FiberManualRecordIcon
          fontSize="small"
          sx={{
            color: statusSettings[status.state].color,
          }}
        />
      </IconButton>
    </HtmlTooltip>
  );
}

export default WebsiteStatusCell;
