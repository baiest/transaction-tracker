type ExtendedBody = BodyInit | object | null | undefined;

export interface ExtendedRequestInit extends Omit<RequestInit, "body"> {
  body?: ExtendedBody;
}

export type FetchClient = <T>(
  path: string,
  options?: ExtendedRequestInit
) => Promise<T>;
