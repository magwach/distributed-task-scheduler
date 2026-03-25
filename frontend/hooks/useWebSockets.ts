"use client";

import { Task } from "@/services/types";
import { useEffect, useRef, useState } from "react";

export const useWebSocket = ({
  onMessage,
}: {
  onMessage: (taskUpdate: Task) => void;
}) => {
  const socketRef = useRef<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);

  useEffect(() => {
    let isMounted = true;

    const connect = () => {
      const protocol = window.location.protocol === "https:" ? "wss" : "ws";
      const host =
        process.env.NEXT_PUBLIC_WEB_SOCKET_URL ?? window.location.host;
      const webSocketUrl = `${protocol}://${host}/ws`;

      const socket = new WebSocket(webSocketUrl);
      socketRef.current = socket;

      socket.onopen = () => {
        if (!isMounted) return;
        console.log("WebSocket connected");
        setConnected(true);
      };

      socket.onmessage = (event) => {
        try {
          const data: Task = JSON.parse(event.data);
          onMessage(data);
        } catch (err) {
          console.error("Invalid WS message:", err);
        }
      };

      socket.onclose = () => {
        if (!isMounted) return;
        console.log("WebSocket disconnected");
        setConnected(false);
        setTimeout(connect, 3000);
      };

      socket.onerror = (err) => {
        console.error("WebSocket error:", err);
        socket.close();
      };
    };

    connect();

    return () => {
      isMounted = false;
      socketRef.current?.close();
    };
  }, [onMessage]);
};
