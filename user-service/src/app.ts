// src/index.ts
import express from "express";
import session from "express-session";
import cors from "cors";
import { PrismaClient } from "@prisma/client";
import initPassport from "./routes/auth/passport";
import UserService from "./services/UserService";
import routes from "./routes/routes";
import CourseService from "./services/CourseService";
import EnvironmentService from "./services/EnvironmentService";

const app = express();
const port = process.env.PORT || 8080;
const prisma = new PrismaClient();

app.use(express.json());
app.use(
  session({
    secret: process.env.SESSION_SECRET || "secret",
    resave: false,
    saveUninitialized: false,
  })
);

app.use(cors());

initPassport(app);

export const userService = new UserService(prisma);
export const courseService = new CourseService(prisma);
export const environmentService = new EnvironmentService(prisma);

app.use("/", routes);

// // Basic health check endpoint
// app.get("/health", (req, res) => {
//   res.json({ status: "ok" });
// });

// // Example endpoint to create a course

// // Example endpoint to assign a role to a user in a course
// app.post("/course-roles", async (req, res) => {
//   try {
//     const { userId, courseId, role } = req.body;
//     const courseRole = await prisma.courseRole.create({
//       data: {
//         userId,
//         courseId,
//         role,
//       },
//     });
//     res.json(courseRole);
//   } catch (error) {
//     res.status(400).json({ error: "Failed to assign role" });
//   }
// });

// app.post("/assignment", async (req, res) => {
//   try {
//     const { courseId, title, description, dueDate } = req.body;
//     const assignment = await prisma.assignment.create({
//       data: {
//         courseId,
//         title,
//         description,
//         dueDate: new Date(dueDate),
//       },
//     });
//     res.json(assignment);
//   } catch (error) {
//     res.status(400).json({ error: "Failed to create assignment" });
//   }
// });

// app.post("/environment", async (req, res) => {
//   try {
//     const { assignmentId, name, baseUrl } = req.body;

//     const assignment = await prisma.assignment.findUnique({
//       where: {
//         id: assignmentId,
//       },
//     });

//     if (!assignment) {
//       res.status(404).json({ error: "Assignment not found" });
//       return;
//     }

//     const environment = await prisma.environment.create({
//       data: {
//         assignmentId: assignment.id,
//         name,
//         baseUrl,
//       },
//     });
//     res.json(environment);
//   } catch (error) {
//     res.status(400).json({ error: "Failed to create environment" });
//   }
// });

// app.get("/courses/:courseCode/assignments", async (req, res) => {
//   try {
//     const courseCode = req.params.courseCode;
//     const course = await prisma.course.findUnique({
//       where: {
//         courseCode,
//       },
//     });

//     if (!course) {
//       res.status(404).json({ error: "Course not found" });
//       return;
//     }

//     const assignments = await prisma.assignment.findMany({
//       where: {
//         courseId: course.id,
//       },
//       include: {
//         Environment: true,
//       },
//     });

//     res.json(assignments);
//   } catch (error) {
//     res.status(400).json({ error: "Failed to get assignments" });
//   }
// });

app.listen(port, () => {
  console.log(`Server running at http://localhost:${port}`);
});
