import { Router } from "express";
import {
  createAssignmentHandler,
  createCourseHandler,
  deleteAssignmentHandler,
  getAssignmentsHandler,
  getCoursesHandler,
} from "./course.service";

const router = Router();

router.post("/create", createCourseHandler);

router.post("/getcourses", getCoursesHandler);

router.post("/:courseCode/assignment", createAssignmentHandler);

router.get("/:courseCode/assignments", getAssignmentsHandler);

router.post("/:courseCode/assignment/delete", deleteAssignmentHandler);

export default router;
