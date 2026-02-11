import { useEffect, useState } from "react";
import { treeApiService } from "../api/api.service";
import type { SearchResult } from "../model/tree.model";

export function useSearch(query: string): SearchResult | null {
  const [result, setResult] = useState<SearchResult | null>(null);

  useEffect(() => {
    if (!query) {
      setResult(null);
      return;
    }

    let cancelled = false;

    const timer = setTimeout(async () => {
      try {
        const data = await treeApiService.searchFiles(query);
        if (!cancelled) {
          setResult(data);
        }
      } catch {
        // Silently handle error
      }
    }, 300);

    return () => {
      cancelled = true;
      clearTimeout(timer);
    };
  }, [query]);

  return result;
}
