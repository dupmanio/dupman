import * as React from "react";
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
import { useState } from "react";
import { useSnackbar } from "notistack";
import { WebsiteRepository } from "@/lib/repositories/website";
import { useForm } from "react-hook-form";
import { yupResolver } from "@hookform/resolvers/yup";
import * as Yup from "yup";

export interface WebsiteFormDialogProps {
  open: boolean;
  onClose: () => void;
}

function WebsiteFormDialog({ open, onClose }: WebsiteFormDialogProps) {
  const [url, setUrl] = useState<string>("");
  const [token, setToken] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);

  const validationSchema = Yup.object().shape({
    url: Yup.string().required().url(),
    token: Yup.string().required(),
  });

  const { enqueueSnackbar } = useSnackbar();
  const { register, handleSubmit, formState } = useForm({
    resolver: yupResolver(validationSchema),
  });

  const { errors } = formState;

  const submitWebsite = async () => {
    setLoading(true);

    WebsiteRepository.create({ url, token })
      .then(() => {
        setLoading(false);

        enqueueSnackbar("Website has been created successfully!", {
          variant: "success",
        });

        onClose();
      })
      .catch((reason) => {
        setLoading(false);

        enqueueSnackbar(
          `Unable to create website: ${reason.response.data.error}`,
          {
            variant: "error",
          },
        );
      });
  };

  return (
    <Dialog open={open} onClose={onClose}>
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
          value={url}
          onInput={(e) => setUrl(e.target.value)}
          error={!!errors.url}
          helperText={errors.url?.message}
          {...register("url")}
        />
        <TextField
          required
          margin="dense"
          id="token"
          label="Token"
          type="password"
          fullWidth
          variant="standard"
          value={token}
          onInput={(e) => setToken(e.target.value)}
          error={!!errors.token}
          helperText={errors.token?.message}
          {...register("token")}
        />
      </DialogContent>
      <DialogActions>
        <LoadingButton onClick={handleSubmit(submitWebsite)} loading={loading}>
          Add
        </LoadingButton>
        <Button onClick={onClose}>Cancel</Button>
      </DialogActions>
    </Dialog>
  );
}

export default WebsiteFormDialog;
