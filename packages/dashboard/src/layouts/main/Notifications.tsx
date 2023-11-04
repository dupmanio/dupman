import React, { useState } from "react";
import useSWR from "swr";
import useSWRInfinite from "swr/infinite";
import { useSnackbar } from "notistack";

import {
  Box,
  Badge,
  Divider,
  IconButton,
  Menu,
  MenuItem,
  Tooltip,
  useTheme,
  Typography,
} from "@mui/material";

import NotificationsIcon from "@mui/icons-material/Notifications";

import NotificationItem from "@/components/NotificationItem";
import { NotifyRepository } from "@/lib/repositories/notify";
import { NotificationOnResponse } from "@/types/dtos/notification";
import { useRealtimeNotifications } from "@/lib/http/client/notify-realtime";

const PAGE_SIZE = 5;

function Notifications() {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [totalPages, setTotalPages] = useState<number>(0);
  const [totalNotificationsCount, setTotalNotificationsCount] =
    useState<number>(0);
  const [unreadNotificationsCount, setUnreadNotificationsCount] =
    useState<number>(0);
  const [notifications, setNotifications] = useState<NotificationOnResponse[]>(
    [],
  );

  const open = Boolean(anchorEl);

  const theme = useTheme();
  const { enqueueSnackbar } = useSnackbar();

  useSWR(["/notification/count"], NotifyRepository.getCount, {
    onSuccess(data) {
      setUnreadNotificationsCount(data.data ?? 0);
    },
  });

  const { data, size, setSize, isLoading } = useSWRInfinite(
    (page, prevData) => {
      if (prevData && !prevData.pagination.totalItems) return null;

      if (page === 0) return ["/notification", 1, PAGE_SIZE];

      return ["/notification", page + 1, PAGE_SIZE];
    },
    ([_, page, limit]) => NotifyRepository.getAll(page, limit),
    {
      onSuccess(data) {
        if (data && data.length > 0) {
          const lastElement = data[data.length - 1];
          const allNotifications = data
            .flat()
            .filter((page) => page.data != null && page.data != undefined)
            .map((page) => page.data)
            .flat();

          setNotifications(allNotifications as NotificationOnResponse[]);
          setTotalPages(lastElement?.pagination?.totalPages ?? 0);
          setTotalNotificationsCount(lastElement?.pagination?.totalItems ?? 0);
        }
      },
    },
  );

  useRealtimeNotifications((notification) => {
    setNotifications((prevNotifications) => [
      notification,
      ...prevNotifications,
    ]);
    setTotalNotificationsCount(
      (prevTotalNotificationsCount) => prevTotalNotificationsCount + 1,
    );
    setUnreadNotificationsCount(
      (prevUnreadNotificationsCount) => prevUnreadNotificationsCount + 1,
    );

    // TODO: handle action.
    enqueueSnackbar(notification.message, {
      variant: "info",
    });
  });

  const isLoadingMore =
    isLoading || (size > 0 && data && typeof data[size - 1] === "undefined");

  const scrollHandler = (e: React.UIEvent<HTMLElement>) => {
    const scrollRemaining =
      e.target.scrollHeight - (e.target.clientHeight + e.target.scrollTop);

    if (scrollRemaining <= 10 && size < totalPages && !isLoadingMore) {
      setSize(size + 1);
    }
  };

  const readAllCallback = () => {
    NotifyRepository.markAllAsRead().then(() => {
      const updatedNotifications = notifications.map((notification) => {
        notification.seen = true;

        return notification;
      });

      setNotifications(updatedNotifications);
      setUnreadNotificationsCount(0);
    });
  };

  const deleteAllCallback = () => {
    NotifyRepository.deleteAll().then(() => {
      setNotifications([]);
      setTotalNotificationsCount(0);
      setUnreadNotificationsCount(0);
      setTotalPages(0);
    });
  };

  const processNotificationCallback = (
    notification: NotificationOnResponse,
  ) => {
    // @todo: implement action processing.

    if (!notification.seen) {
      NotifyRepository.markAsRead(notification.id).then(() => {
        notification.seen = true;

        setNotifications([...notifications]);
        setUnreadNotificationsCount(unreadNotificationsCount - 1);
      });
    }
  };

  return (
    <>
      <Tooltip title="Notifications">
        <IconButton
          color="inherit"
          onClick={(e) => {
            setAnchorEl(e.currentTarget);
          }}
          sx={{ ml: 2 }}
          aria-controls={open ? "notification-menu" : undefined}
          aria-expanded={open ? "true" : undefined}
          aria-haspopup="true"
        >
          <Badge badgeContent={unreadNotificationsCount} color="secondary">
            <NotificationsIcon />
          </Badge>
        </IconButton>
      </Tooltip>

      <Menu
        anchorEl={anchorEl}
        id="notification-menu"
        open={open}
        onClose={() => {
          setAnchorEl(null);
        }}
        slotProps={{
          paper: {
            onScroll: scrollHandler,
            elevation: 0,
            sx: {
              width: theme.typography.pxToRem(400),
              maxHeight: "40vh",
              overflowY: "auto",
              filter: "drop-shadow(0px 2px 8px rgba(0,0,0,0.32))",
              mt: 1.5,
              "&:before": {
                content: '""',
                display: "block",
                position: "absolute",
                top: 0,
                right: 14,
                width: 10,
                height: 10,
                bgcolor: "background.paper",
                transform: "translateY(-50%) rotate(45deg)",
                zIndex: 0,
              },
            },
          },
        }}
        transformOrigin={{ horizontal: "right", vertical: "top" }}
        anchorOrigin={{ horizontal: "right", vertical: "bottom" }}
      >
        <Box
          sx={{
            display: "flex",
            justifyContent: "space-between",
          }}
        >
          <MenuItem
            disabled={unreadNotificationsCount == 0}
            component={"a"}
            href="#"
            sx={{ color: theme.palette.primary.main }}
            onClick={readAllCallback}
          >
            Mark all as read
          </MenuItem>

          <MenuItem
            disabled={totalNotificationsCount == 0}
            component={"a"}
            href="#"
            sx={{ color: theme.palette.error.main }}
            onClick={deleteAllCallback}
          >
            Remove all
          </MenuItem>
        </Box>

        <Divider />

        {notifications.length == 0 && (
          <Box
            sx={{
              display: "flex",
              justifyContent: "center",
              padding: 2,
            }}
          >
            <Typography variant="body1">
              Notifications? Nowhere to be found.
            </Typography>
          </Box>
        )}

        {notifications.length > 0 &&
          notifications?.map((notification) => (
            <NotificationItem
              key={notification.id.toString()}
              data={notification}
              onClick={() => processNotificationCallback(notification)}
            />
          ))}

        <Divider />

        <Box
          sx={{
            display: "flex",
            justifyContent: "center",
          }}
        >
          <MenuItem
            disabled={size >= totalPages || isLoadingMore}
            component={"a"}
            href="#"
            onClick={() => {
              setSize(size + 1);
            }}
          >
            {isLoadingMore ? "Loading..." : "Load More"}
          </MenuItem>
        </Box>
      </Menu>
    </>
  );
}

export default Notifications;
