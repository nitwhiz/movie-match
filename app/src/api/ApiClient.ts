import { Media, MediaType, RecommendedMedia } from '../model/Media';
import axiosStatic, { Axios } from 'axios';
import { Vote, VoteType } from '../model/Vote';
import { User } from '../model/User';
import { Match } from '../model/Match';
import jwtDecode from 'jwt-decode';
import { Login } from '../model/Login';
import Cookies from 'js-cookie';
import EventEmitter from 'eventemitter3';

interface Results<T> {
  results: T[];
}

interface Token {
  exp: number;
  orig_iat: number;
  userId: string;
}

const ACCESS_TOKEN_EXPIRY_THRESHOLD = 15 * 60 * 1000;

const ACCESS_TOKEN_COOKIE_NAME = 'jwt';
const ACCESS_TOKEN_COOKIE_EXPIRY_THRESHOLD = 15 * 60 * 1000;

export default class ApiClient extends EventEmitter<{
  unauthorized: () => void;
  logout: () => void;
}> {
  private readonly baseUrl: string;

  private readonly axios: Axios;

  private accessToken: string;

  private accessTokenExpiry: number;

  private tokenRefreshPromise: Promise<boolean> | null;

  constructor(baseUrl: string) {
    super();

    this.baseUrl = baseUrl;

    this.axios = axiosStatic.create({
      baseURL: baseUrl,
    });

    this.accessToken = '';
    this.accessTokenExpiry = 0;
    this.tokenRefreshPromise = null;
  }

  public setAccessToken(token: string) {
    this.accessToken = token;
    this.accessTokenExpiry = jwtDecode<Token>(token).exp * 1000;

    this.axios.defaults.headers.common.Authorization = `Bearer ${this.accessToken}`;
    this.axios.defaults.headers.common['Content-Type'] = 'application/json';

    Cookies.set(ACCESS_TOKEN_COOKIE_NAME, token, {
      expires: new Date(
        this.accessTokenExpiry + ACCESS_TOKEN_COOKIE_EXPIRY_THRESHOLD
      ),
      sameSite: 'Strict',
    });
  }

  public loadAccessTokenFromCookie(): ApiClient {
    const accessToken = Cookies.get(ACCESS_TOKEN_COOKIE_NAME);

    if (accessToken) {
      this.setAccessToken(accessToken);
    }

    return this;
  }

  private async checkAccessToken(): Promise<boolean> {
    if (this.tokenRefreshPromise) {
      return this.tokenRefreshPromise;
    }

    if (this.accessToken === '' || this.accessTokenExpiry === 0) {
      this.emit('unauthorized');
      return true;
    }

    // renew token 15 min before it's invalid
    if (Date.now() >= this.accessTokenExpiry - ACCESS_TOKEN_EXPIRY_THRESHOLD) {
      if (!this.tokenRefreshPromise) {
        this.tokenRefreshPromise = this.refreshToken().catch(() =>
          this.emit('unauthorized')
        );
      }

      return this.tokenRefreshPromise.finally(
        () => (this.tokenRefreshPromise = null)
      );
    }

    return true;
  }

  public login(login: Login): Promise<void> {
    return this.axios
      .post<{ token: string }>('/auth/login', login, {
        transformRequest: [
          (data, headers) => {
            delete headers['Authorization'];

            return data;
          },
          // @ts-ignore todo: wtf?
          ...this.axios.defaults.transformRequest,
        ],
      })
      .then(({ data }) => data)
      .then(({ token }) => {
        if (token) {
          this.setAccessToken(token);
          return;
        }

        throw new Error('no token found');
      });
  }

  public logout(): Promise<void> {
    return this.axios
      .post('/auth/logout')
      .then(() => Cookies.remove(ACCESS_TOKEN_COOKIE_NAME))
      .then(() => {
        this.emit('logout');
      });
  }

  public refreshToken(): Promise<boolean> {
    return this.axios
      .post<{ token: string }>('/auth/refresh_token')
      .then(({ data }) => data)
      .then(({ token }) => {
        if (token) {
          this.setAccessToken(token);

          return true;
        }

        return false;
      });
  }

  public async me(): Promise<User> {
    await this.checkAccessToken();

    return this.axios
      .get<{ user: User }>('/me')
      .then(({ data }) => data)
      .then(({ user }) => user);
  }

  public async getMedia(mediaId: string): Promise<Media> {
    await this.checkAccessToken();

    return this.axios.get<Media>(`/media/${mediaId}`).then(({ data }) => data);
  }

  public async getPosterBlobUrl(mediaId: string): Promise<string | null> {
    await this.checkAccessToken();

    return this.axios
      .get<Blob>(`/media/${mediaId}/poster`, {
        responseType: 'blob',
      })
      .then(({ data }) => data)
      .then((blob) => URL.createObjectURL(blob))
      .catch(() => null);
  }

  public async getRecommendedMedia(
    belowScore: string = '100'
  ): Promise<RecommendedMedia[]> {
    await this.checkAccessToken();

    return this.axios
      .get<Results<RecommendedMedia>>(
        `/me/recommended?belowScore=${belowScore}`
      )
      .then(({ data }) => data)
      .then(({ results }) => results);
  }

  public async getMediaVote(mediaId: string): Promise<VoteType> {
    await this.checkAccessToken();

    return this.axios
      .get<{ voteType: VoteType }>(`/me/vote/${mediaId}`)
      .then(({ data: { voteType } }) => voteType);
  }

  /**
   * returns true if there was a match
   */
  public async voteMedia(
    mediaId: string,
    voteType: VoteType
  ): Promise<boolean> {
    await this.checkAccessToken();

    return this.axios
      .put<{ isMatch: boolean }>(`/media/${mediaId}/vote`, {
        voteType,
      })
      .then(({ data: { isMatch } }) => isMatch);
  }

  public async setMediaSeen(mediaId: string): Promise<void> {
    await this.checkAccessToken();

    return this.axios.put(`/media/${mediaId}/seen`);
  }

  public async getUsers(): Promise<User[]> {
    await this.checkAccessToken();

    return this.axios
      .get<Results<User>>('/users')
      .then(({ data }) => data)
      .then(({ results }) => results);
  }

  public async getMatches(mediaType: MediaType | null): Promise<Match[]> {
    await this.checkAccessToken();

    return this.axios
      .get<Results<Match>>(
        `/matches${mediaType !== null ? `?type=${mediaType}` : ''}`
      )
      .then(({ data }) => data)
      .then(({ results }) => results);
  }

  public async getVotes(): Promise<Vote[]> {
    await this.checkAccessToken();

    return this.axios
      .get<Results<Vote>>('/me/votes')
      .then(({ data }) => data)
      .then(({ results }) => results);
  }

  public async searchMedia(query: string): Promise<Media[]> {
    await this.checkAccessToken();

    return this.axios
      .get<Results<Media>>('/search/media', {
        params: {
          query,
        },
      })
      .then(({ data }) => data)
      .then(({ results }) => results);
  }
}
