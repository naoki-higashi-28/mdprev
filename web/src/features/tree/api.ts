import type { SearchResponse, TreeResponse } from "../../shared/types";

export async function fetchTree(path: string): Promise<TreeResponse> {
  const params = new URLSearchParams();
  if (path) {
    params.set("path", path);
  }
  const url = `/api/tree${params.toString() ? `?${params.toString()}` : ""}`;
  const res = await fetch(url);
  if (!res.ok) {
    throw new Error(`Failed to fetch tree: ${res.status}`);
  }
  return res.json() as Promise<TreeResponse>;
}

export async function searchFiles(query: string): Promise<SearchResponse> {
  const params = new URLSearchParams({ q: query });
  const res = await fetch(`/api/search?${params.toString()}`);
  if (!res.ok) {
    throw new Error(`Failed to search files: ${res.status}`);
  }
  return res.json() as Promise<SearchResponse>;
}
