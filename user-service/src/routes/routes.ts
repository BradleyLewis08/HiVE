import { Router } from "express";
import userController from "./user/user.controller";
import authController from "./auth/auth.controller";
import courseController from "./course/course.controller";
import environmentController from "./environment/environment.controller";

const api = Router();

api.use("/user", userController);
api.use("/auth", authController);
api.use("/course", courseController);
api.use("/environment", environmentController);
export default Router().use("/api", api);
