import * as React from "react";
import { useState } from "react";
import { useSnackbar } from "notistack";
import * as Yup from "yup";
import { useFormik } from "formik";

import { LoadingButton } from "@mui/lab";
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  TextField,
} from "@mui/material";

import { WebsiteRepository } from "@/lib/repositories/website";

export interface WebsiteFormDialogProps {
  open: boolean;
  onClose: () => void;
}

function WebsiteFormDialog({ open, onClose }: WebsiteFormDialogProps) {
  const validationSchema = Yup.object().shape({
    url: Yup.string().required().url(),
    token: Yup.string().required(),
  });

  const { enqueueSnackbar } = useSnackbar();
  const formik = useFormik({
    initialValues: {
      url: "",
      token: "",
    },
    validationSchema: validationSchema,
    onSubmit: ({ url, token }, { setSubmitting }) => {
      WebsiteRepository.create({ url, token })
        .then(() => {
          enqueueSnackbar("Website has been created successfully!", {
            variant: "success",
          });

          onClose();
        })
        .catch((reason) => {
          enqueueSnackbar(
            `Unable to create website: ${reason.response.data.error}`,
            {
              variant: "error",
            },
          );
        })
        .finally(() => {
          setSubmitting(false);
        });
    },
  });

  const handleClose = () => {
    formik.resetForm();
    onClose();
  };

  return (
    <Dialog open={open} onClose={handleClose}>
      <DialogTitle>Add new Website</DialogTitle>
      <DialogContent>
        <DialogContentText>Add new Website for monitoring.</DialogContentText>
        <TextField
          autoFocus
          required
          margin="dense"
          id="url"
          label="URL"
          type="url"
          fullWidth
          variant="standard"
          value={formik.values.url}
          onChange={formik.handleChange}
          onBlur={formik.handleBlur}
          error={formik.touched.url && Boolean(formik.errors.url)}
          helperText={formik.touched.url && formik.errors.url}
        />
        <TextField
          required
          margin="dense"
          id="token"
          label="Token"
          type="password"
          fullWidth
          variant="standard"
          value={formik.values.token}
          onChange={formik.handleChange}
          onBlur={formik.handleBlur}
          error={formik.touched.token && Boolean(formik.errors.token)}
          helperText={formik.touched.url && formik.errors.token}
        />
      </DialogContent>
      <DialogActions>
        <LoadingButton
          onClick={formik.submitForm}
          loading={formik.isSubmitting}
        >
          Add
        </LoadingButton>
        <Button onClick={handleClose}>Cancel</Button>
      </DialogActions>
    </Dialog>
  );
}

export default WebsiteFormDialog;
