import * as React from "react";
import { Grid, Paper } from "@mui/material";
import Websites from "@/components/Websites";
import { PageAccess } from "@/config/page-accesss";
import Typography from "@mui/material/Typography";

function Home() {
  return (
    <Grid item xs={12}>
      <Paper sx={{ p: 2, display: "flex", flexDirection: "column" }}>
        <Typography component="h2" variant="h6" color="primary" gutterBottom>
          Websites
        </Typography>

        <Websites />
      </Paper>
    </Grid>
  );
}

Home.Access = PageAccess.SECURED;

export default Home;
