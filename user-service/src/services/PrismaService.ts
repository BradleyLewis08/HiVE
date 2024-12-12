import { PrismaClient } from "@prisma/client";

export default class PrismaService {
  prisma: PrismaClient;

  constructor(prisma: PrismaClient) {
    this.prisma = prisma;
  }
}
