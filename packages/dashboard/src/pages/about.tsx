import * as React from "react";
import { Grid, Paper, Typography } from "@mui/material";

import { PageAccess } from "@/config/page-accesss";

function About() {
  return (
    <Grid item xs={12}>
      <Paper sx={{ p: 2, display: "flex", flexDirection: "column" }}>
        <Typography component="h2" variant="h6" color="primary" gutterBottom>
          About
        </Typography>
      </Paper>
    </Grid>
  );
}

About.Access = PageAccess.SECURED;
export default About;
