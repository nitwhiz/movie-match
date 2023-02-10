import { Media, MediaType } from '../model/Media';
import axiosStatic, { Axios } from 'axios';
import { VoteType } from '../model/Vote';
import { User } from '../model/User';
import { Match } from '../model/Match';

interface Results<T> {
  results: T[];
}

export default class ApiClient {
  private readonly baseUrl: string;

  private readonly axios: Axios;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;

    this.axios = axiosStatic.create({
      baseURL: baseUrl,
    });
  }

  public getAllMedia(): Promise<Media[]> {
    return this.axios
      .get<Results<Media>>('/media')
      .then(({ data }) => data)
      .then(({ results }) => results);
  }

  public getMedia(mediaId: string): Promise<Media> {
    return this.axios.get<Media>(`/media/${mediaId}`).then(({ data }) => data);
  }

  public getPosterUrl(mediaId: string): string {
    return `${this.baseUrl}/media/${mediaId}/poster`;
  }

  public getRecommendedMedia(
    userId: string,
    page: number = 0
  ): Promise<Media[]> {
    return this.axios
      .get<Results<Media>>(`/user/${userId}/media/recommended?page=${page}`)
      .then(({ data }) => data)
      .then(({ results }) => results);
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

  public setMediaSeen(userId: string, mediaId: string): Promise<void> {
    return this.axios.put(`/user/${userId}/media/${mediaId}/seen`);
  }

  public getUsers(): Promise<User[]> {
    return this.axios
      .get<Results<User>>('/user')
      .then(({ data }) => data)
      .then(({ results }) => results);
  }

  public getMatches(
    userId: string,
    mediaType: MediaType | null
  ): Promise<Match[]> {
    return this.axios
      .get<Results<Match>>(
        `/user/${userId}/match${mediaType !== null ? `?type=${mediaType}` : ''}`
      )
      .then(({ data }) => data)
      .then(({ results }) => results);
  }
}
