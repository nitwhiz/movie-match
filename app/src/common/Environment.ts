export const API_SERVER_BASE_URL: string = (
  import.meta.env.VITE_API_SERVER_BASE_URL || ''
).replace(/\/+$/, '');
