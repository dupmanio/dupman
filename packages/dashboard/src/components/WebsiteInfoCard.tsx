import * as React from "react";

import {
  Chip,
  Link,
  List,
  ListItem,
  ListItemText,
  Typography,
} from "@mui/material";
import { useTheme } from "@mui/material/styles";

import { formatISO } from "@/lib/util/time";
import { useStatusSettings } from "@/lib/util/website-status";
import { Website } from "@/types/entities/website";

export interface WebsiteInfoCardProps {
  data: Website;
}

function WebsiteInfoCard({ data }: WebsiteInfoCardProps) {
  const theme = useTheme();
  const statusSettings = useStatusSettings(data.status);

  const listItemTypography = {
    primary: {
      fontSize: theme.typography.pxToRem(14),
      color: theme.palette.text.secondary,
    },
    secondary: {
      fontSize: theme.typography.pxToRem(16),
      color: theme.palette.text.primary,
    },
  };

  return (
    <List>
      <ListItem>
        <ListItemText
          inset
          primary="ID"
          secondary={data.id}
          primaryTypographyProps={listItemTypography.primary}
          secondaryTypographyProps={listItemTypography.secondary}
        />
      </ListItem>
      <ListItem>
        <ListItemText
          inset
          primary="URL"
          secondary={
            <Link href={data.url} target="_blank">
              {data.url}
            </Link>
          }
          primaryTypographyProps={listItemTypography.primary}
          secondaryTypographyProps={listItemTypography.secondary}
        />
      </ListItem>
      <ListItem>
        <ListItemText
          inset
          primary="Created At"
          secondary={formatISO(data.createdAt)}
          primaryTypographyProps={listItemTypography.primary}
          secondaryTypographyProps={listItemTypography.secondary}
        />
        <ListItemText
          inset
          primary="Updated At"
          secondary={formatISO(data.updatedAt)}
          primaryTypographyProps={listItemTypography.primary}
          secondaryTypographyProps={listItemTypography.secondary}
        />
        {data.status.state && (
          <ListItemText
            inset
            primary="Last Scanned At"
            secondary={formatISO(data.status.updatedAt)}
            primaryTypographyProps={listItemTypography.primary}
            secondaryTypographyProps={listItemTypography.secondary}
          />
        )}
      </ListItem>

      {data.status.state && (
        <ListItem>
          <ListItemText
            inset
            primary="Status"
            secondary={
              <>
                <Chip
                  label={statusSettings?.title}
                  color={statusSettings?.color}
                />
                {statusSettings?.additionalText && (
                  <Typography
                    variant="body2"
                    sx={{
                      mt: 1,
                      p: 1,
                      backgroundColor: theme.palette.background.default,
                      borderRadius: 1,
                    }}
                  >
                    {statusSettings?.additionalText}
                  </Typography>
                )}
              </>
            }
            primaryTypographyProps={listItemTypography.primary}
            secondaryTypographyProps={listItemTypography.secondary}
          />
        </ListItem>
      )}
    </List>
  );
}

export default WebsiteInfoCard;
