import { getApiClient } from "../../../shared/api/get-api-client";
import { RootDirectoryInfo } from "../model/root-directory-info.model";
import { pathBarApiEndpoints } from "./api.endpoints";
import type { ApiSchemaRootDirectoryInfoResponse } from "./api.types";

const apiClient = getApiClient();

function toRootDirectoryInfoModel(
  schema: ApiSchemaRootDirectoryInfoResponse,
): RootDirectoryInfo {
  return new RootDirectoryInfo({
    rootDirectoryName: schema.name,
  });
}

export const pathBarApiService = {
  fetchRootDirectoryInfo: async (): Promise<RootDirectoryInfo> => {
    const schema = await apiClient.request<ApiSchemaRootDirectoryInfoResponse>({
      method: "GET",
      url: pathBarApiEndpoints.info,
    });
    return toRootDirectoryInfoModel(schema);
  },
} as const;
