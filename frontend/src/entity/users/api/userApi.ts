import { getApplicationServer } from '../../../constants';
import RequestOptions from '../../../shared/api/RequestOptions';
import { accessTokenHelper } from '../../../shared/api/accessTokenHelper';
import { apiHelper } from '../../../shared/api/apiHelper';
import type { SignInRequest } from '../model/SignInRequest';
import type { SignInResponse } from '../model/SignInResponse';
import type { SignUpRequest } from '../model/SignUpRequest';

const listeners: (() => void)[] = [];

const saveAuthorizedData = (accessToken: string, userId: string) => {
  accessTokenHelper.saveAccessToken(accessToken);
  accessTokenHelper.saveUserId(userId);
};

const notifyAuthListeners = () => {
  for (const listener of listeners) {
    listener();
  }
};

export const userApi = {
  async signUp(signUpRequest: SignUpRequest) {
    const requestOptions: RequestOptions = new RequestOptions();
    requestOptions.setBody(JSON.stringify(signUpRequest));
    return apiHelper.fetchPostRaw(`${getApplicationServer()}/api/v1/users/signup`, requestOptions);
  },

  async signIn(signInRequest: SignInRequest): Promise<SignInResponse> {
    const requestOptions: RequestOptions = new RequestOptions();
    requestOptions.setBody(JSON.stringify(signInRequest));

    return apiHelper
      .fetchPostJson(`${getApplicationServer()}/api/v1/users/signin`, requestOptions)
      .then((response: unknown): SignInResponse => {
        const typedResponse = response as SignInResponse;
        saveAuthorizedData(typedResponse.token, typedResponse.userId);
        notifyAuthListeners();
        return typedResponse;
      });
  },

  async isAnyUserExists(): Promise<boolean> {
    const requestOptions: RequestOptions = new RequestOptions();
    return apiHelper
      .fetchGetJson(
        `${getApplicationServer()}/api/v1/users/is-any-user-exist`,
        requestOptions,
        true,
      )
      .then((response: unknown) => {
        const typedResponse = response as { isExist: boolean };
        return typedResponse.isExist;
      });
  },

  isAuthorized: (): boolean => !!accessTokenHelper.getAccessToken(),

  logout: () => {
    accessTokenHelper.cleanAccessToken();
  },

  // listeners

  addAuthListener: (listener: () => void) => {
    listeners.push(listener);
  },

  removeAuthListener: (listener: () => void) => {
    listeners.splice(listeners.indexOf(listener), 1);
  },

  notifyAuthListeners: (): void => {
    for (const listener of listeners) {
      listener();
    }
  },

  saveAuthorizedData,
};
