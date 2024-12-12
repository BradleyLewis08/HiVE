import PrismaService from "./PrismaService";

import type { Role, User } from "@prisma/client";
import { SystemRole } from "@prisma/client";
import type Yalie from "Yalie";

const YALIES_ENDPOINT = "https://yalies.io/api/people";

export default class UserService extends PrismaService {
  public async getUserFromYalies(netId: string): Promise<Yalie> {
    const body = {
      filters: {
        netId: netId,
      },
    };

    try {
      const response = await fetch(YALIES_ENDPOINT, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${process.env.YALIES_API_KEY}`,
        },
        body: JSON.stringify(body),
      });

      console.log(response);

      if (!response.ok) {
        throw new Error(
          `"Failed to fetch display name from Yalies: ${response.status} ${response.statusText}`
        );
      }

      const data = (await response.json()) as Yalie[];

      if (data.length === 0) {
        throw new Error(`No user found with netID ${netId}`);
      }

      return data[0];
    } catch (error) {
      throw new Error(`Failed to fetch display name from Yalies: ${error}`);
    }
  }

  public async createUser(netId: string): Promise<User> {
    let yaliesData: Yalie;
    try {
      yaliesData = await this.getUserFromYalies(netId);
    } catch (err) {
      console.log(err);
      yaliesData = {
        first_name: "Unknown",
        last_name: "Unknown",
        netId: netId,
        email: "Unknown",
      };
    }

    const name = `${yaliesData.first_name} ${yaliesData.last_name}`;
    const email = yaliesData.email;

    const user = await this.prisma.user.create({
      data: {
        netId,
        name,
        email,
      },
    });

    return user;
  }

  public async getUser(netId: string): Promise<User | null> {
    const user = await this.prisma.user.findUnique({
      where: {
        netId,
      },
    });

    if (user !== null) {
      return user;
    }

    return null;
  }

  public async changeUserPermissions(
    netId: string,
    systemRole: SystemRole
  ): Promise<User> {
    const user = await this.prisma.user.update({
      where: {
        netId,
      },
      data: {
        systemRole,
      },
    });

    return user;
  }
}
