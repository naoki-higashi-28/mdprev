import { getApiClient } from "../../../shared/api/get-api-client";
import { SearchResult, Tree, TreeEntry } from "../model/tree.model";
import { treeApiEndpoints } from "./api.endpoints";
import type {
  ApiSchemaSearchQuery,
  ApiSchemaSearchResponse,
  ApiSchemaTreeEntry,
  ApiSchemaTreeQuery,
  ApiSchemaTreeResponse,
} from "./api.types";

function toTreeEntryModel(schema: ApiSchemaTreeEntry): TreeEntry {
  return new TreeEntry({
    type: schema.type,
    name: schema.name,
    path: schema.path,
    ext: schema.ext,
    size: schema.size,
    mtime: schema.mtime,
  });
}

function toTreeModel(schema: ApiSchemaTreeResponse): Tree {
  return new Tree({
    path: schema.path,
    entries: schema.entries.map(toTreeEntryModel),
  });
}

function toSearchResultModel(schema: ApiSchemaSearchResponse): SearchResult {
  return new SearchResult({
    query: schema.query,
    entries: schema.entries.map(toTreeEntryModel),
  });
}

const apiClient = getApiClient();

export const treeApiService = {
  fetchTree: async (path: string): Promise<Tree> => {
    const req: ApiSchemaTreeQuery = path ? { path } : {};
    const schema = await apiClient.request<ApiSchemaTreeResponse>({
      method: "GET",
      url: treeApiEndpoints.tree,
      params: req,
    });
    return toTreeModel(schema);
  },
  searchFiles: async (query: string): Promise<SearchResult> => {
    const req: ApiSchemaSearchQuery = { q: query };
    const schema = await apiClient.request<ApiSchemaSearchResponse>({
      method: "GET",
      url: treeApiEndpoints.search,
      params: req,
    });
    return toSearchResultModel(schema);
  },
} as const;
