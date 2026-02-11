import { useEffect, useRef } from "react";
import useSWR, { useSWRConfig } from "swr";
import { getWatchApiService } from "../../../shared/api/watch-api-service";
import { previewApiService } from "../api/api.service";
import type { FileContent } from "../model/file-content.model";

const watchApiService = getWatchApiService();

interface UseFileContentResult {
  content: FileContent | null;
  error: string | null;
  loading: boolean;
}

export function useFileContent(path: string): UseFileContentResult {
  const pathRef = useRef(path);
  const { mutate } = useSWRConfig();

  const { data, error, isLoading } = useSWR(
    path ? ["file", path] : null,
    ([, p]) => previewApiService.fetchFile(p),
  );

  useEffect(() => {
    pathRef.current = path;
  }, [path]);

  useEffect(() => {
    const unsubscribe = watchApiService.subscribe((message) => {
      if (message.type === "file_change" && message.path === pathRef.current) {
        mutate(["file", pathRef.current]);
      }
    });

    return unsubscribe;
  }, [mutate]);

  return {
    content: data ?? null,
    error: error
      ? error instanceof Error
        ? error.message
        : "Failed to fetch file"
      : null,
    loading: isLoading,
  };
}
