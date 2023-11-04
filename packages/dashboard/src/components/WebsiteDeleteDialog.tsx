import * as React from "react";
import { useState } from "react";
import { useSnackbar } from "notistack";

import { LoadingButton } from "@mui/lab";
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
} from "@mui/material";

import { WebsiteRepository } from "@/lib/repositories/website";
import { Website } from "@/types/entities/website";

export interface WebsiteDeleteDialogProps {
  website: Website | null;
  open: boolean;
  onClose: () => void;
}

function WebsiteDeleteDialog({
  website,
  open,
  onClose,
}: WebsiteDeleteDialogProps) {
  const [isDeleting, setIsDeleting] = useState<boolean>(false);

  const { enqueueSnackbar } = useSnackbar();

  const handleDelete = () => {
    setIsDeleting(true);
    WebsiteRepository.delete(website?.id ?? "")
      .then(() => {
        enqueueSnackbar("Website has been deleted successfully!", {
          variant: "success",
        });

        onClose();
      })
      .catch((reason) => {
        enqueueSnackbar(
          `Unable to delete website: ${reason.response.data.error}`,
          {
            variant: "error",
          },
        );

        onClose();
      })
      .finally(() => {
        setIsDeleting(false);
      });
  };

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogTitle>Delete website</DialogTitle>
      <DialogContent>
        <DialogContentText>
          Are you sure you want to delete this Website? This action is
          irreversible!
        </DialogContentText>
      </DialogContent>
      <DialogActions>
        <LoadingButton onClick={handleDelete} loading={isDeleting}>
          Delete
        </LoadingButton>
        <Button onClick={onClose}>Cancel</Button>
      </DialogActions>
    </Dialog>
  );
}

export default WebsiteDeleteDialog;
