import * as React from "react";

import {
  IconButton,
  Tooltip,
  TooltipProps,
  tooltipClasses,
  styled,
} from "@mui/material";

import FiberManualRecordIcon from "@mui/icons-material/FiberManualRecord";

import { StatusOnWebsitesResponse } from "@/types/dtos/status";
import { useStatusSettings } from "@/lib/util/website-status";

export interface WebsiteStatusCellProps {
  status: StatusOnWebsitesResponse;
}

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
  const settings = useStatusSettings(status);

  return (
    <HtmlTooltip
      title={
        settings.title +
        (settings.additionalText ? `: ${settings.additionalText}` : "")
      }
      leaveDelay={100}
    >
      <IconButton>
        <FiberManualRecordIcon
          fontSize="small"
          sx={{
            color: settings.colorCode,
          }}
        />
      </IconButton>
    </HtmlTooltip>
  );
}

export default WebsiteStatusCell;
