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

function WebsiteUpdates() {
  const rows = [
    {
      name: "drupal",
      title: "Drupal core",
      link: "https://www.drupal.org/project/drupal",
      type: "core",
      currentVersion: "9.2.1",
      latestVersion: "9.5.10",
      recommendedVersion: "9.5.10",
      installType: "official",
      status: 1,
    },
  ];

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
        {rows.map((row) => (
          <TableRow
            key={row.name}
            sx={{
              "&:last-child td, &:last-child th": { border: 0 },
            }}
          >
            <TableCell component="th" scope="row">
              <Link href={row.link} target="_blank">
                {row.title} ({row.name})
              </Link>
            </TableCell>
            <TableCell>{row.currentVersion}</TableCell>
            <TableCell>{row.latestVersion}</TableCell>
            <TableCell>{row.recommendedVersion}</TableCell>
            <TableCell>
              <Stack direction="row" spacing={0.5}>
                <Chip
                  label={`Install Type: ${row.installType}`}
                  color={row.installType === "official" ? "success" : "error"}
                />
                <Chip
                  label={`Type: ${row.type}`}
                  color={row.type === "core" ? "primary" : "default"}
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
