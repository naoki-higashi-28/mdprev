import type { Method } from "axios";
import axios from "axios";

interface ApiRequestConfig<TParams extends object = object, TData = unknown> {
  method: Method;
  url: string;
  params?: TParams;
  data?: TData;
  responseType?: "json";
}

interface ApiRequestTextConfig<
  TParams extends object = object,
  TData = unknown,
> {
  method: Method;
  url: string;
  params?: TParams;
  data?: TData;
  responseType: "text";
}

export interface ApiClient {
  request<TRes, TParams extends object = object, TData = unknown>(
    config: ApiRequestConfig<TParams, TData>,
  ): Promise<TRes>;
  request<TParams extends object = object, TData = unknown>(
    config: ApiRequestTextConfig<TParams, TData>,
  ): Promise<string>;
}

const axiosInstance = axios.create();

export function getApiClient(): ApiClient {
  return {
    async request<TRes, TParams extends object = object, TData = unknown>(
      config:
        | ApiRequestConfig<TParams, TData>
        | ApiRequestTextConfig<TParams, TData>,
    ): Promise<TRes | string> {
      const response = await axiosInstance.request<TRes | string>({
        method: config.method,
        url: config.url,
        params: config.params,
        data: config.data,
        responseType: config.responseType ?? "json",
      });

      return response.data;
    },
  };
}
