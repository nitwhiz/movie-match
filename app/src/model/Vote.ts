export const enum VoteType {
  NONE = 'none',
  POSITIVE = 'positive',
  NEGATIVE = 'negative',
  NEUTRAL = 'neutral',
}

export interface Vote {
  id: string;
  userId: string;
  mediaId: string;
  type: VoteType;
  createdAt: string;
  updatedAt: string;
}
