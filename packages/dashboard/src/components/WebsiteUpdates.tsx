import * as React from "react";

import {
  Table,
  TableHead,
  TableRow,
  TableCell,
  TableBody,
  Link,
  Stack,
  Chip,
} from "@mui/material";
import { UpdatesOnResponse } from "@/types/dtos/update";

export interface WebsiteUpdatesProps {
  updates: UpdatesOnResponse;
}

function WebsiteUpdates({ updates }: WebsiteUpdatesProps) {
  return (
    <Table aria-label="website updates">
      <TableHead>
        <TableRow>
          <TableCell>Module</TableCell>
          <TableCell>Current Version</TableCell>
          <TableCell>Latest Version</TableCell>
          <TableCell>Recommended Version</TableCell>
          <TableCell></TableCell>
        </TableRow>
      </TableHead>
      <TableBody>
        {updates.map((update) => (
          <TableRow
            key={update.name}
            sx={{
              "&:last-child td, &:last-child th": { border: 0 },
            }}
          >
            <TableCell component="th" scope="row">
              <Link href={update.link} target="_blank">
                {update.title} ({update.name})
              </Link>
            </TableCell>
            <TableCell>{update.currentVersion}</TableCell>
            <TableCell>{update.latestVersion}</TableCell>
            <TableCell>{update.recommendedVersion}</TableCell>
            <TableCell>
              <Stack direction="row" spacing={0.5}>
                <Chip
                  label={`Install Type: ${update.installType}`}
                  color={
                    update.installType === "official" ? "success" : "error"
                  }
                />
                <Chip
                  label={`Type: ${update.type}`}
                  color={update.type === "core" ? "primary" : "default"}
                />
              </Stack>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}

export default WebsiteUpdates;
