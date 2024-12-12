import { userService } from "../../app";
import type { User } from "@prisma/client";

export async function handleUserLogin(netId: string): Promise<User> {
  try {
    let user;
    user = await userService.getUser(netId);

    if (user === null) {
      user = await userService.createUser(netId);
    }

    return user;
  } catch (error) {
    throw new Error(`Failed to handle user login: ${error}`);
  }
}
