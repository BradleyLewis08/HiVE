/*
  Warnings:

  - You are about to drop the column `role` on the `CourseRole` table. All the data in the column will be lost.

*/
-- AlterTable
ALTER TABLE "CourseRole" DROP COLUMN "role",
ADD COLUMN     "courseRole" "Role" NOT NULL DEFAULT 'STUDENT';
