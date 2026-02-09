export interface TreeEntry {
  type: "dir" | "file";
  name: string;
  path: string;
  ext?: string;
  size?: number;
  mtime?: number;
}

export interface TreeResponse {
  path: string;
  entries: TreeEntry[];
}

export interface SearchResponse {
  query: string;
  entries: TreeEntry[];
}
