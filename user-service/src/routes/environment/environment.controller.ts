import { Router } from "express";
import { createEnvironmentHandler } from "./environment.service";

const router = Router();

router.post("/create", createEnvironmentHandler);

export default router;
