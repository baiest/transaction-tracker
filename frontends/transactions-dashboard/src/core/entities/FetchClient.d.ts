export type FetchClient = <T>(
  path: string,
  options?: RequestInit
) => Promise<T>;
