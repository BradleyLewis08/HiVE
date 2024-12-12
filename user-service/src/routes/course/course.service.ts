import type { Request, Response } from "express";
import { courseService } from "../../app";
import { Role } from "@prisma/client";
import ResponseFactory from "../../services/ResponseFactory";

export async function createCourseHandler(req: Request, res: Response) {
  try {
    const { netId, courseCode, name, description } = req.body;

    if (!netId || !courseCode || !name || !description) {
      res.status(400).json({ error: "Missing required fields" });
      return;
    }

    const course = await courseService.createCourse(
      req.body.courseCode,
      req.body.name,
      req.body.description
    );

    // Assign the user to the course as an instructor
    await courseService.assignUserToCourse(
      netId,
      course.courseCode,
      Role.INSTRUCTOR
    );

    new ResponseFactory(res).success(course);
  } catch (error) {
    new ResponseFactory(res).internalServerError("Failed to create course");
  }
}

export async function getCoursesHandler(req: Request, res: Response) {
  try {
    const { netId } = req.body;
    const courses = await courseService.getCourses(netId);
    new ResponseFactory(res).success(courses);
  } catch (error) {
    new ResponseFactory(res).internalServerError("Failed to get courses");
  }
}

export async function createAssignmentHandler(req: Request, res: Response) {
  try {
    const { courseCode } = req.params;
    const { title, imageURL, netIDs, description } = req.body;
    const response = await courseService.createAssignment(
      courseCode,
      title,
      imageURL,
      netIDs,
      description
    );
    new ResponseFactory(res).success(response);
  } catch (error) {
    new ResponseFactory(res).internalServerError("Failed to create assignment");
  }
}

export async function getAssignmentsHandler(req: Request, res: Response) {
  try {
    const { courseCode } = req.params;
    const assignments = await courseService.getAssignments(courseCode);
    new ResponseFactory(res).success(assignments);
  } catch (error) {
    new ResponseFactory(res).internalServerError("Failed to get assignments");
  }
}

export async function deleteAssignmentHandler(req: Request, res: Response) {
  try {
    const { assignmentId } = req.body;
    const deletedAssignment = await courseService.deleteAssignment(
      assignmentId
    );
    new ResponseFactory(res).success(deletedAssignment);
  } catch (error) {
    new ResponseFactory(res).internalServerError("Failed to delete assignment");
  }
}
