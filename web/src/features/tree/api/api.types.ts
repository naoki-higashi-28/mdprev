export interface ApiSchemaTreeQuery {
  path?: string;
}

export interface ApiSchemaSearchQuery {
  q: string;
}

export interface ApiSchemaTreeEntry {
  type: "dir" | "file";
  name: string;
  path: string;
  ext?: string;
  size?: number;
  mtime?: number;
}

export interface ApiSchemaTreeResponse {
  path: string;
  entries: ApiSchemaTreeEntry[];
}

export interface ApiSchemaSearchResponse {
  query: string;
  entries: ApiSchemaTreeEntry[];
}
