import { useEffect, useState } from "react";
import { fetchFile } from "./api";
import { useFileWatcher } from "./useFileWatcher";

interface UseFileContentResult {
  content: string | null;
  error: string | null;
  loading: boolean;
}

export function useFileContent(path: string): UseFileContentResult {
  const [content, setContent] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const revision = useFileWatcher(path);

  // biome-ignore lint/correctness/useExhaustiveDependencies: revision triggers re-fetch on file change
  useEffect(() => {
    if (!path) {
      setContent(null);
      setError(null);
      setLoading(false);
      return;
    }

    let cancelled = false;
    setLoading(true);

    async function load() {
      try {
        const text = await fetchFile(path);
        if (!cancelled) {
          setContent(text);
          setError(null);
        }
      } catch (e) {
        if (!cancelled) {
          setContent(null);
          setError(e instanceof Error ? e.message : "Failed to fetch file");
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    load();
    return () => {
      cancelled = true;
    };
  }, [path, revision]);

  return { content, error, loading };
}
