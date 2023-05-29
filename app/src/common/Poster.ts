import ApiClient from '../api/ApiClient';

const GC_INTERVAL = 1000;

const GC_THRESHOLD = 1000;

export default class Poster {
  private static instancesByMediaId: Record<string, Poster> = {};

  public static getByMediaId(apiClient: ApiClient, mediaId: string): Poster {
    if (!Poster.instancesByMediaId[mediaId]) {
      Poster.instancesByMediaId[mediaId] = new Poster(
        mediaId,
        apiClient.getPosterBlobUrl(mediaId).then((url) => {
          if (!url) {
            throw new Error('unable to load poster');
          }

          return url;
        })
      );
    }

    return Poster.instancesByMediaId[mediaId].use();
  }

  public static freeAll(): void {
    for (const poster of Object.values(Poster.instancesByMediaId)) {
      poster.free(true);
    }
  }

  public static startGC(): void {
    this.cycleGC();
  }

  private static cycleGC(): void {
    for (const [mediaId, poster] of Object.entries(Poster.instancesByMediaId)) {
      if (poster.shouldGC()) {
        delete Poster.instancesByMediaId[mediaId];
        poster.revoke().catch((e) => console.error(e));
      }
    }

    window.setTimeout(() => {
      Poster.cycleGC();
    }, GC_INTERVAL);
  }

  private refCount: number = 0;

  private lastFree: number = Infinity;

  private constructor(
    public readonly mediaId: string,
    private readonly urlPromise: Promise<string>
  ) {
    Poster.instancesByMediaId[mediaId] = this.free();
  }

  private use(): Poster {
    ++this.refCount;

    return this;
  }

  public async getUrl(): Promise<string> {
    return await this.urlPromise;
  }

  public free(unset: boolean = false): Poster {
    if (unset) {
      this.refCount = 0;
    }

    if (this.refCount > 0) {
      --this.refCount;
    }

    this.lastFree = Date.now();

    return this;
  }

  public shouldGC(): boolean {
    return this.refCount <= 0 && Date.now() - this.lastFree >= GC_THRESHOLD;
  }

  public async revoke(): Promise<void> {
    URL.revokeObjectURL(await this.urlPromise);
  }
}
