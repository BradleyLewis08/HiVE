import type { Request, Response } from "express";
import { userService } from "../../app";
import ResponseFactory from "../../services/ResponseFactory";
import type { SystemRole } from "@prisma/client";

export async function getUserHandler(req: Request, res: Response) {
  const { netId } = req.params;
  const user = await userService.getUser(netId);

  if (user === null) {
    new ResponseFactory(res).notFound(`User with netId ${netId} not found`);
  }

  new ResponseFactory(res).success(user);
}

export async function changeUserPermissions(req: Request, res: Response) {
  const { netId, systemRole } = req.body;

  // Check if role is valid
  if (!["STUDENT", "FACULTY", "ADMIN"].includes(systemRole)) {
    new ResponseFactory(res).badRequest(
      "Role must be one of STUDENT, FACULTY, or ADMIN"
    );
  }

  const user = await userService.getUser(netId);

  if (user === null) {
    new ResponseFactory(res).notFound(`User with netId ${netId} not found`);
  }

  const newUser = await userService.changeUserPermissions(
    netId,
    systemRole as SystemRole
  );

  new ResponseFactory(res).success(newUser);
}
