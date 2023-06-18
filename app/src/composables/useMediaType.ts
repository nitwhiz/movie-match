import { MediaType } from '../model/Media';

const mediaTypeLabelsSingular: { [K in MediaType]: string } = {
  [MediaType.ALL]: 'Alle',
  [MediaType.MOVIE]: 'Film',
  [MediaType.TV]: 'Serie',
};

const mediaTypeLabelsPlural: { [K in MediaType]: string } = {
  [MediaType.ALL]: 'Alle',
  [MediaType.MOVIE]: 'Filme',
  [MediaType.TV]: 'Serien',
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
