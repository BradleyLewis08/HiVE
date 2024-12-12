import { Router } from "express";
import { getUserHandler, changeUserPermissions } from "./user.service";

const router = Router();

router.get("/:netId", getUserHandler);

router.post("/permissions", changeUserPermissions);

export default router;
