import { useEffect, useState } from "react";
import type { SearchResponse } from "../../shared/types";
import { searchFiles } from "./api";

export function useSearch(query: string): SearchResponse | null {
  const [result, setResult] = useState<SearchResponse | null>(null);

  useEffect(() => {
    if (!query) {
      setResult(null);
      return;
    }

    let cancelled = false;

    const timer = setTimeout(async () => {
      try {
        const data = await searchFiles(query);
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
