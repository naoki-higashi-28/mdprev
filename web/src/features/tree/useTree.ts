import { useEffect, useState } from "react";
import type { TreeResponse } from "../../shared/types";
import { fetchTree } from "./api";

interface UseTreeResult {
  tree: TreeResponse | null;
  error: string | null;
  loading: boolean;
}

export function useTree(): UseTreeResult {
  const [tree, setTree] = useState<TreeResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;

    async function load() {
      try {
        const data = await fetchTree("");
        if (!cancelled) {
          setTree(data);
          setError(null);
        }
      } catch (e) {
        if (!cancelled) {
          setError(e instanceof Error ? e.message : "Failed to fetch tree");
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
  }, []);

  return { tree, error, loading };
}
