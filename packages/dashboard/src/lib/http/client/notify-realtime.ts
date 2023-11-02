import { useEffect } from "react";
import getNextConfig from "next/config";
import { useSession } from "next-auth/react";
import { fetchEventSource } from "@microsoft/fetch-event-source";

import { NotificationOnResponse } from "@/types/dtos/notification";

function useRealtimeNotifications(
  onNotification: (notification: NotificationOnResponse) => void,
) {
  const { publicRuntimeConfig } = getNextConfig();
  const session = useSession();

  useEffect(() => {
    fetchEventSource(
      `${publicRuntimeConfig.NOTIFY_URL}/notification/realtime`,
      {
        method: "GET",
        headers: {
          Authorization: `Bearer ${session.data?.accessToken}`,
        },
        onmessage(message) {
          if (message.event === "notification") {
            const notification = JSON.parse(
              message.data,
            ) as NotificationOnResponse;

            onNotification(notification);
          }
        },
      },
    );
  }, []);
}

export { useRealtimeNotifications };
