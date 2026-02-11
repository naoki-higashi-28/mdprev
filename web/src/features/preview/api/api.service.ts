import { getApiClient } from "../../../shared/api/get-api-client";
import { FileContent, FileNotFoundError } from "../model/file-content.model";
import { previewApiEndpoints } from "./api.endpoints";
import type { ApiSchemaFileContent, ApiSchemaFileQuery } from "./api.types";

const apiClient = getApiClient();

export const previewApiService = {
  fetchFile: async (path: string): Promise<FileContent> => {
    const req: ApiSchemaFileQuery = { path };
    try {
      const schema = await apiClient.request({
        method: "GET",
        url: previewApiEndpoints.file,
        params: req,
        responseType: "text",
      });
      return new FileContent(schema as ApiSchemaFileContent);
    } catch (error) {
      if (
        typeof error === "object" &&
        error !== null &&
        "response" in error &&
        typeof error.response === "object" &&
        error.response !== null &&
        "status" in error.response &&
        error.response.status === 404
      ) {
        throw new FileNotFoundError("File not found");
      }
      throw error;
    }
  },
} as const;
