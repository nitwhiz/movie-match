import { Media } from '../model/Media';
import axiosStatic, { Axios } from 'axios';
import { VoteType } from '../model/Vote';
import { User } from '../model/User';
import { Match } from '../model/Match';
import {
  DEFAULT_VALIDITY,
  getCached,
  MINUTES,
  putStored,
} from '../cache/CacheStorage';

interface Results<T> {
  results: T[];
}

const getCacheKey = (path: string) => `api:${path}`;

export default class ApiClient {
  private readonly baseUrl: string;

  private readonly axios: Axios;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;

    this.axios = axiosStatic.create({
      baseURL: baseUrl,
    });
  }

  private requestUncached<DataType>(url: string): Promise<DataType> {
    return this.axios.get<DataType>(url).then(({ data }) => data);
  }

  private requestCached<DataType>(
    url: string,
    validity: number = DEFAULT_VALIDITY
  ): Promise<DataType> {
    return getCached<DataType>(
      getCacheKey(url),
      () => this.axios.get<DataType>(url).then(({ data }) => data),
      validity
    );
  }

  public getAllMedia(): Promise<Media[]> {
    return this.requestCached<Results<Media>>('/media', 5 * MINUTES).then(
      ({ results }) => results
    );
  }

  public getMedia(mediaId: string): Promise<Media> {
    return this.requestCached<Media>(`/media/${mediaId}`);
  }

  public getPosterUrl(mediaId: string): string {
    return `${this.baseUrl}/media/${mediaId}/poster`;
  }

  public getRecommendedMedia(
    userId: string,
    page: number = 0
  ): Promise<Media[]> {
    return this.requestUncached<Results<Media>>(
      `/user/${userId}/media/recommended?page=${page}`
    ).then(({ results }) => {
      for (const media of results) {
        putStored(getCacheKey(`/media/${media.id}`), media);
      }

      return results;
    });
  }

  /**
   * returns true if there was a match
   */
  public voteMedia(
    userId: string,
    mediaId: string,
    voteType: VoteType
  ): Promise<boolean> {
    return this.axios
      .put<{ isMatch: boolean }>(`/user/${userId}/media/${mediaId}/vote`, {
        voteType,
      })
      .then(({ data: { isMatch } }) => isMatch);
  }

  public setMediaSeen(
    userId: string,
    mediaId: string,
    seen: boolean = true
  ): Promise<void> {
    if (seen) {
      return this.axios.post(`/user/${userId}/media/${mediaId}/seen`);
    }

    return this.axios.delete(`/user/${userId}/media/${mediaId}/seen`);
  }

  public getUsers(): Promise<User[]> {
    return this.requestCached<Results<User>>('/user').then(
      ({ results }) => results
    );
  }

  public getMatches(userId: string): Promise<Match[]> {
    return this.requestUncached<Results<Match>>(`/user/${userId}/match`).then(
      ({ results }) => results
    );
  }
}
