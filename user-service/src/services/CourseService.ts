import type { Role } from "@prisma/client";
import PrismaService from "./PrismaService";
import { environmentService } from "../app";

export default class CourseService extends PrismaService {
  public async getCourse(courseCode: string) {
    return this.prisma.course.findUnique({
      where: {
        courseCode,
      },
    });
  }

  public async createCourse(
    courseCode: string,
    name: string,
    description: string
  ) {
    return this.prisma.course.create({
      data: {
        courseCode,
        name,
        description,
      },
    });
  }

  public async assignUserToCourse(
    netId: string,
    courseCode: string,
    role: Role
  ) {
    try {
      // First, find the user by netId
      const user = await this.prisma.user.findUnique({
        where: { netId },
      });

      if (!user) {
        throw new Error(`User with netId ${netId} not found`);
      }

      // Find the course by courseCode
      const course = await this.prisma.course.findUnique({
        where: { courseCode },
      });

      if (!course) {
        throw new Error(`Course with code ${courseCode} not found`);
      }

      // Create or update the course role using upsert
      // This will either create a new role or update an existing one
      const courseRole = await this.prisma.courseRole.upsert({
        where: {
          userId_courseId: {
            userId: user.id,
            courseId: course.id,
          },
        },
        update: {
          courseRole: role,
        },
        create: {
          userId: user.id,
          courseId: course.id,
          courseRole: role,
        },
      });

      return courseRole;
    } catch (error) {
      // Log the error for debugging
      console.error("Error in assignUserToCourse:", error);

      // Re-throw the error with a more specific message
      if (error instanceof Error) {
        throw error;
      }
      throw new Error("Failed to assign user to course");
    }
  }

  public async getCourses(netId: string) {
    try {
      const user = await this.prisma.user.findUnique({
        where: { netId },
        include: {
          courseRoles: {
            include: {
              course: true,
            },
          },
        },
      });

      if (!user) {
        throw new Error(`User with netId ${netId} not found`);
      }

      const courses = user.courseRoles.map((courseRole) => ({
        id: courseRole.course.id,
        courseCode: courseRole.course.courseCode,
        name: courseRole.course.name,
        description: courseRole.course.description,
        role: courseRole.courseRole,
        createdAt: courseRole.course.createdAt,
        updatedAt: courseRole.course.updatedAt,
      }));

      return courses;
    } catch (error) {
      console.error("Error in getCourses:", error);
      if (error instanceof Error) {
        throw error;
      }
      throw new Error("Failed to get courses");
    }
  }

  public async createAssignment(
    courseCode: string,
    assignmentName: string,
    imageURL: string,
    netIDs: string[],
    description?: string
  ) {
    try {
      // First verify the course exists
      const course = await this.prisma.course.findUnique({
        where: { courseCode },
      });

      if (!course) {
        throw new Error(`Course with id ${courseCode} not found`);
      }

      // Create the assignment
      const assignment = await this.prisma.assignment.create({
        data: {
          title: assignmentName,
          description: description,
          courseId: course.id,
        },
      });
      // Now, we create the environment using the environment service
      try {
        const envResponse = await environmentService.createEnvironment(
          assignmentName,
          course.name,
          imageURL,
          netIDs
        );

        const environment = await this.prisma.environment.create({
          data: {
            name: `${course.name}-${assignmentName} Environment`,
            assignmentId: assignment.id,
            baseUrl: "http://example.com",
          },
        });

        const response = {
          assignment,
          environment,
        };

        return response;
      } catch (error) {
        // If there is an error creating the environment, delete the assignment
        await this.prisma.assignment.delete({
          where: {
            id: assignment.id,
          },
        });
        throw new Error(`Failed to create environment: ${error}`);
      }
    } catch (error) {
      console.error("Error in createAssignment:", error);
      if (error instanceof Error) {
        throw error;
      }
      throw new Error("Failed to create assignment");
    }
  }

  public async getAssignments(courseCode: string) {
    try {
      const course = await this.prisma.course.findUnique({
        where: { courseCode },
      });

      if (!course) {
        throw new Error(`Course with code ${courseCode} not found`);
      }

      const assignments = await this.prisma.assignment.findMany({
        where: {
          courseId: course.id,
        },
        include: {
          Environment: true,
        },
      });

      return assignments;
    } catch (error) {
      if (error instanceof Error) {
        throw error;
      }
      throw new Error("Failed to get assignments");
    }
  }

  public async deleteAssignment(assignmentId: string) {
    try {
      const assignment = await this.prisma.assignment.delete({
        where: {
          id: assignmentId,
        },
      });
      return assignment;
    } catch (error) {
      console.error("Error in deleteAssignment:", error);
      if (error instanceof Error) {
        throw error;
      }
      throw new Error("Failed to delete assignment");
    }
  }
}
