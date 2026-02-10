import { useEffect } from "react";
import useSWR, { useSWRConfig } from "swr";
import type { TreeResponse } from "../../shared/types";
import { fetchTree } from "./api";

interface UseTreeResult {
  tree: TreeResponse | null;
  error: string | null;
  loading: boolean;
}

export function useTree(): UseTreeResult {
  const { mutate } = useSWRConfig();

  const { data, error, isLoading } = useSWR(["tree", ""], ([, p]) =>
    fetchTree(p),
  );

  useEffect(() => {
    const eventSource = new EventSource("/api/watch");

    eventSource.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        if (msg.type === "tree_change") {
          const dir: string = msg.path;
          mutate(["tree", dir]);
          if (dir !== "") {
            mutate(["tree", ""]);
          }
        }
      } catch {
        // Ignore non-JSON messages
      }
    };

    return () => {
      eventSource.close();
    };
  }, [mutate]);

  return {
    tree: data ?? null,
    error: error
      ? error instanceof Error
        ? error.message
        : "Failed to fetch tree"
      : null,
    loading: isLoading,
  };
}
