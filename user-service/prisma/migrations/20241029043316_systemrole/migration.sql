/*
  Warnings:

  - You are about to drop the column `isAdmin` on the `User` table. All the data in the column will be lost.
  - You are about to drop the column `isFaculty` on the `User` table. All the data in the column will be lost.

*/
-- CreateEnum
CREATE TYPE "SystemRole" AS ENUM ('ADMIN', 'FACULTY', 'STUDENT');

-- AlterTable
ALTER TABLE "User" DROP COLUMN "isAdmin",
DROP COLUMN "isFaculty",
ADD COLUMN     "systemRole" "SystemRole" NOT NULL DEFAULT 'STUDENT';
