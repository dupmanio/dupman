import * as React from "react";
import { useRouter } from "next/router";
import useSWR from "swr";

import {
  Grid,
  Paper,
  TableContainer,
  Typography,
  Divider,
  Box,
} from "@mui/material";

import { WebsiteRepository } from "@/lib/repositories/website";
import PageLoader from "@/components/PageLoader";
import WebsiteUpdates from "@/components/WebsiteUpdates";
import WebsiteInfoCard from "@/components/WebsiteInfoCard";
import { StatusState } from "@/types/dtos/status";
import WebsitePreviewImage from "@/components/WebsitePreviewImage";

function WebsitePage() {
  const router = useRouter();

  const { data, isLoading, error } = useSWR(
    ["/website/:id", router.query.id],
    ([_, id]) => WebsiteRepository.getSingle(id as string),
    {
      onError(error) {
        if (error.response.status !== 401) {
          router.replace("/404");
        }
      },
    },
  );

  return (
    <>
      <Grid item xs={12}>
        <Paper sx={{ p: 2, display: "flex", flexDirection: "column" }}>
          {(isLoading || error) && <PageLoader size={40} />}

          {!isLoading && data?.data && (
            <Grid container spacing={2}>
              <Grid item xs={4}>
                <Box>
                  <WebsitePreviewImage websiteId={data?.data.id} />
                </Box>
              </Grid>
              <Grid item xs={8}>
                <Paper>
                  <WebsiteInfoCard data={data.data} />
                </Paper>
              </Grid>

              {data.data.status.state == StatusState.NeedsUpdate && (
                <Grid item xs={12}>
                  <Divider />

                  <Typography component="h2" variant="h6" gutterBottom>
                    Updates
                  </Typography>

                  <TableContainer component={Paper}>
                    <WebsiteUpdates updates={data.data.updates} />
                  </TableContainer>
                </Grid>
              )}
            </Grid>
          )}
        </Paper>
      </Grid>
    </>
  );
}

export default WebsitePage;
