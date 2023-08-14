import { PageAccess } from "@/config/page-accesss";
import { Grid, Paper } from "@mui/material";
import * as React from "react";
import Typography from "@mui/material/Typography";

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
