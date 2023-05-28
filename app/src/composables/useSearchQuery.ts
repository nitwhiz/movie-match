import { ref } from 'vue';

const searchQuery = ref('');

export const useSearchQuery = () => ({
  searchQuery,
});
