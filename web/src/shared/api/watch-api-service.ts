export interface WatchMessage {
  type: "file_change" | "tree_change";
  path: string;
}

export interface WatchApiService {
  subscribe: (onMessage: (message: WatchMessage) => void) => () => void;
}

function isWatchMessage(value: unknown): value is WatchMessage {
  if (typeof value !== "object" || value === null) {
    return false;
  }
  if (!("type" in value) || !("path" in value)) {
    return false;
  }

  const type = value.type;
  const path = value.path;
  return (
    (type === "file_change" || type === "tree_change") &&
    typeof path === "string"
  );
}

export function getWatchApiService(): WatchApiService {
  return {
    subscribe(onMessage) {
      const eventSource = new EventSource("/api/watch");

      eventSource.onmessage = (event) => {
        try {
          const raw: unknown = JSON.parse(event.data);
          if (isWatchMessage(raw)) {
            onMessage(raw);
          }
        } catch {
          // Ignore non-JSON messages.
        }
      };

      return () => {
        eventSource.close();
      };
    },
  };
}
