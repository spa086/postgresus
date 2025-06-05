import { accessTokenHelper } from '.';
import RequestOptions from './RequestOptions';

const handleOrThrowMessageIfResponseError = async (
  url: string,
  response: Response,
  handleNotAuthorizedError = true,
) => {
  if (handleNotAuthorizedError && response.status === 401) {
    accessTokenHelper?.cleanAccessToken();
    window.location.reload();
  }

  if (response.status === 502 || response.status === 504) {
    throw new Error('failed to fetch');
  }

  if (response.status >= 400 && response.status <= 600) {
    let errorMessage: string | undefined;

    try {
      const json = (await response.json()) as { message?: string; error?: string };
      errorMessage = json.message || json.error;
    } catch {
      try {
        errorMessage = await response.text();
      } catch {
        /* ignore */
      }
    }

    throw new Error(errorMessage ?? `${url}: ${await response.text()}`);
  }
};

const makeRequest = async (url: string, optionsWrapper: RequestOptions): Promise<Response> => {
  const response = await fetch(url, optionsWrapper.toRequestInit());
  await handleOrThrowMessageIfResponseError(url, response);
  return response;
};

export const apiHelper = {
  fetchPostJson: async <T>(url: string, requestOptions?: RequestOptions): Promise<T> => {
    const optionsWrapper = (requestOptions ?? new RequestOptions())
      .setMethod('POST')
      .addHeader('Content-Type', 'application/json')
      .addHeader('Access-Control-Allow-Methods', 'POST')
      .addHeader('Accept', 'application/json')
      .addHeader('Authorization', accessTokenHelper.getAccessToken());

    const response = await makeRequest(url, optionsWrapper);

    return response.json();
  },

  fetchPostRaw: async (url: string, requestOptions?: RequestOptions): Promise<string> => {
    const optionsWrapper = (requestOptions ?? new RequestOptions())
      .setMethod('POST')
      .addHeader('Content-Type', 'application/json')
      .addHeader('Access-Control-Allow-Methods', 'POST')
      .addHeader('Accept', 'application/json')
      .addHeader('Authorization', accessTokenHelper.getAccessToken());

    const response = await makeRequest(url, optionsWrapper);

    return response.text();
  },

  fetchPostBlob: async (url: string, requestOptions?: RequestOptions): Promise<Blob> => {
    const optionsWrapper = (requestOptions ?? new RequestOptions())
      .setMethod('POST')
      .addHeader('Content-Type', 'application/json')
      .addHeader('Access-Control-Allow-Methods', 'POST')
      .addHeader('Authorization', accessTokenHelper.getAccessToken());

    const response = await makeRequest(url, optionsWrapper);

    return response.blob();
  },

  fetchGetJson: async <T>(url: string, requestOptions?: RequestOptions): Promise<T> => {
    const optionsWrapper = (requestOptions ?? new RequestOptions())
      .addHeader('Content-Type', 'application/json')
      .addHeader('Access-Control-Allow-Methods', 'GET')
      .addHeader('Accept', 'application/json')
      .addHeader('Authorization', accessTokenHelper.getAccessToken());

    const response = await makeRequest(url, optionsWrapper);
    return response.json();
  },

  fetchGetRaw: async (url: string, requestOptions?: RequestOptions): Promise<string> => {
    const optionsWrapper = (requestOptions ?? new RequestOptions())
      .addHeader('Content-Type', 'application/json')
      .addHeader('Access-Control-Allow-Methods', 'GET')
      .addHeader('Accept', 'application/json')
      .addHeader('Authorization', accessTokenHelper.getAccessToken());

    const response = await makeRequest(url, optionsWrapper);
    return response.text();
  },

  fetchGetBlob: async (url: string, requestOptions?: RequestOptions): Promise<Blob> => {
    const optionsWrapper = (requestOptions ?? new RequestOptions())
      .addHeader('Content-Type', 'application/json')
      .addHeader('Access-Control-Allow-Methods', 'GET')
      .addHeader('Authorization', accessTokenHelper.getAccessToken());

    const response = await makeRequest(url, optionsWrapper);
    return response.blob();
  },

  fetchPutJson: async <T>(url: string, requestOptions?: RequestOptions): Promise<T> => {
    const optionsWrapper = (requestOptions ?? new RequestOptions())
      .setMethod('PUT')
      .addHeader('Content-Type', 'application/json')
      .addHeader('Access-Control-Allow-Methods', 'PUT')
      .addHeader('Accept', 'application/json')
      .addHeader('Authorization', accessTokenHelper.getAccessToken());

    const response = await makeRequest(url, optionsWrapper);
    return response.json();
  },

  fetchDeleteJson: async <T>(url: string, requestOptions?: RequestOptions): Promise<T> => {
    const optionsWrapper = (requestOptions ?? new RequestOptions())
      .setMethod('DELETE')
      .addHeader('Access-Control-Allow-Methods', 'DELETE')
      .addHeader('Accept', 'application/json')
      .addHeader('Authorization', accessTokenHelper.getAccessToken());

    const response = await makeRequest(url, optionsWrapper);
    return response.json();
  },

  fetchDeleteRaw: async (url: string, requestOptions?: RequestOptions): Promise<string> => {
    const optionsWrapper = (requestOptions ?? new RequestOptions())
      .setMethod('DELETE')
      .addHeader('Access-Control-Allow-Methods', 'DELETE')
      .addHeader('Accept', 'application/json')
      .addHeader('Authorization', accessTokenHelper.getAccessToken());

    const response = await makeRequest(url, optionsWrapper);
    return response.text();
  },
};
