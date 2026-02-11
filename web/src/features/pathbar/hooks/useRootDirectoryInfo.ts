import { useEffect, useState } from "react";
import { pathBarApiService } from "../api/api.service";
import type { RootDirectoryInfo } from "../model/root-directory-info.model";

export function useRootDirectoryInfo(): RootDirectoryInfo | null {
  const [rootDirectoryInfo, setRootDirectoryInfo] =
    useState<RootDirectoryInfo | null>(null);

  useEffect(() => {
    let cancelled = false;

    (async () => {
      try {
        const data = await pathBarApiService.fetchRootDirectoryInfo();
        if (!cancelled) {
          setRootDirectoryInfo(data);
        }
      } catch {
        // Silently handle error
      }
    })();

    return () => {
      cancelled = true;
    };
  }, []);

  return rootDirectoryInfo;
}
