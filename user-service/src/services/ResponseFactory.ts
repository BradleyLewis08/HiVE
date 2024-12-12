import type { Response } from "express";

export default class ResponseFactory {
  private res: Response;

  constructor(res: Response) {
    this.res = res;
  }

  public success<T>(payload: T) {
    this.res.status(200).json(payload);
    return this.res;
  }

  public missingRequiredParam(param: string) {
    this.res
      .status(400)
      .json({ message: `Missing required attribute: ${param}` });
    return this.res;
  }

  public badRequest(message: string) {
    this.res.status(400).json({ message });
    return this.res;
  }

  public unauthorized(message: string) {
    this.res.status(401).json({ message });
    return this.res;
  }

  public notFound(message: string) {
    this.res.status(404).json({ message });
    return this.res;
  }

  public internalServerError(message: string) {
    this.res.status(500).json({ message });
    return this.res;
  }
}
