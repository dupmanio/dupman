import { useEffect } from "react";
import NProgress from "nprogress";

import { Box, CircularProgress } from "@mui/material";
import { SxProps } from "@mui/system/styleFunctionSx";
import { Theme } from "@mui/system/createTheme";

export interface PageLoaderProps {
  sx?: SxProps<Theme>;
  size?: number;
}

function PageLoader({ sx, size }: PageLoaderProps) {
  useEffect(() => {
    NProgress.start();

    return () => {
      NProgress.done();
    };
  }, []);

  return (
    <Box sx={sx} display="flex" alignItems="center" justifyContent="center">
      <CircularProgress size={size} disableShrink thickness={3} />
    </Box>
  );
}

PageLoader.defaultProps = {
  size: 128,
};

export default PageLoader;
