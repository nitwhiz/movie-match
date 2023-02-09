import { MediaType } from '../model/Media';

const mediaTypeLabelsSingular: { [K in MediaType]: string } = {
  [MediaType.MOVIE]: 'Movie',
  [MediaType.TV]: 'TV-Show',
};

const mediaTypeLabelsPlural: { [K in MediaType]: string } = {
  [MediaType.MOVIE]: 'Movies',
  [MediaType.TV]: 'TV-Shows',
};

const getMediaTypeLabelSingular = (mediaType: MediaType) => {
  return mediaTypeLabelsSingular[mediaType];
};

const getMediaTypeLabelPlural = (mediaType: MediaType) => {
  return mediaTypeLabelsPlural[mediaType];
};

export const useMediaType = () => ({
  getMediaTypeLabelSingular,
  getMediaTypeLabelPlural,
});
