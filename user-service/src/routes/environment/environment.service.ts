import type { Request, Response } from "express";
import ResponseFactory from "../../services/ResponseFactory";

export async function createEnvironmentHandler(req: Request, res: Response) {
  new ResponseFactory(res).success({ message: "Environment created" });
}
