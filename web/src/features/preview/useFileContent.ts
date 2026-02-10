import { useEffect, useRef } from "react";
import useSWR, { useSWRConfig } from "swr";
import { fetchFile } from "./api";

interface UseFileContentResult {
  content: string | null;
  error: string | null;
  loading: boolean;
}

export function useFileContent(path: string): UseFileContentResult {
  const pathRef = useRef(path);
  const { mutate } = useSWRConfig();

  const { data, error, isLoading } = useSWR(
    path ? ["file", path] : null,
    ([, p]) => fetchFile(p),
  );

  useEffect(() => {
    pathRef.current = path;
  }, [path]);

  useEffect(() => {
    const eventSource = new EventSource("/api/watch");

    eventSource.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        if (msg.type === "file_change" && msg.path === pathRef.current) {
          mutate(["file", pathRef.current]);
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
    content: data ?? null,
    error: error
      ? error instanceof Error
        ? error.message
        : "Failed to fetch file"
      : null,
    loading: isLoading,
  };
}
