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
    const controller = new AbortController();

    fetchEventSource(
      `${publicRuntimeConfig.DUPMAN_API}/notify/notification/realtime`,
      {
        method: "GET",
        headers: {
          Authorization: `Bearer ${session.data?.accessToken}`,
        },
        signal: controller.signal,
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

    return () => controller.abort();
  }, []);
}

export { useRealtimeNotifications };
