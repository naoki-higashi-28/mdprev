import { useEffect } from "react";
import useSWR, { useSWRConfig } from "swr";
import { getWatchApiService } from "../../../shared/api/watch-api-service";
import { treeApiService } from "../api/api.service";
import type { Tree } from "../model/tree.model";

const watchApiService = getWatchApiService();

interface UseTreeResult {
  tree: Tree | null;
  error: string | null;
  loading: boolean;
}

export function useTree(): UseTreeResult {
  const { mutate } = useSWRConfig();

  const { data, error, isLoading } = useSWR(["tree", ""], ([, p]) =>
    treeApiService.fetchTree(p),
  );

  useEffect(() => {
    const unsubscribe = watchApiService.subscribe((message) => {
      if (message.type === "tree_change") {
        const dir = message.path;
        mutate(["tree", dir]);
        if (dir !== "") {
          mutate(["tree", ""]);
        }
      }
    });

    return unsubscribe;
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
