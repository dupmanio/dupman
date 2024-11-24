import * as React from "react";
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
import { Website } from "@/types/entities/website";

export interface WebsiteFormDialogProps {
  website: Website | null;
  open: boolean;
  onClose: () => void;
}

function WebsiteFormDialog({ website, open, onClose }: WebsiteFormDialogProps) {
  const validationFields = {
    url: Yup.string().required().url(),
    token: Yup.string(),
  };

  if (website === null) {
    validationFields.token = validationFields.token.required();
  }

  const validationSchema = Yup.object().shape(validationFields);

  const { enqueueSnackbar } = useSnackbar();
  const formik = useFormik({
    initialValues: {
      url: website?.url ?? "",
      token: "",
    },
    validationSchema: validationSchema,
    onSubmit: ({ url, token }, { setSubmitting }) => {
      const request =
        website === null
          ? WebsiteRepository.create({ url, token })
          : WebsiteRepository.update(website.id, { url, token });

      request
        .then(() => {
          enqueueSnackbar(
            `Website has been ${
              website === null ? "created" : "updated"
            } successfully!`,
            {
              variant: "success",
            },
          );

          onClose();
        })
        .catch((reason) => {
          enqueueSnackbar(
            `Unable to ${website === null ? "create" : "update"} website: ${
              reason.response.data.error
            }`,
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
      <DialogTitle>
        {website === null ? "Add new Website" : "Edit Website"}
      </DialogTitle>
      <DialogContent>
        <DialogContentText>
          {website === null
            ? "Add new Website for monitoring."
            : "Modify existing Website."}
        </DialogContentText>
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
          {website === null ? "Add" : "Edit"}
        </LoadingButton>
        <Button onClick={handleClose}>Cancel</Button>
      </DialogActions>
    </Dialog>
  );
}

WebsiteFormDialog.defaultProps = {
  website: null,
};

export default WebsiteFormDialog;
