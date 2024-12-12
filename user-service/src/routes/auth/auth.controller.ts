import express, { Router } from "express";
import passport from "passport";
import jwt from "jsonwebtoken";
import { handleUserLogin } from "./auth.service";
import type { User } from "@prisma/client";
import { userService } from "../../app";

const router = Router();

interface CasPayload {
  netId: string;
}

export const authWithCas = function (
  req: express.Request,
  res: express.Response,
  next: express.NextFunction
) {
  passport.authenticate("cas", async function (err: any, payload: CasPayload) {
    if (err) {
      return next(err);
    }

    if (!payload) {
      return next(new Error("CAS auth but no user"));
    }

    const { netId } = payload;

    let user: User;

    try {
      user = await handleUserLogin(netId);
    } catch (error) {
      return next(error);
    }

    const token = jwt.sign({ data: user }, process.env.JWT_SECRET || "secret", {
      expiresIn: "10d",
    });

    req.logIn(user, function (err) {
      if (err) {
        return next(err);
      }

      const redirectUrl = req.query.redirect as string;
      const redirectWithToken = `${redirectUrl}?token=${token}`;

      if (req.query.redirect) {
        return res.redirect(redirectWithToken);
      }

      return res.redirect("/");
    });
  })(req, res, next);
};

router.get("/login", authWithCas);

router.get("/logout", (req: any, res) => {
  req.logOut((err: any) => {
    if (err) {
      return res.json({ success: false });
    }
    res.json({ success: true });
  });
});

router.get("/check", (req, res) => {
  if (req.user) {
    res.json({ auth: true, user: req.user });
  } else {
    res.json({ auth: false });
  }
});

export default router;
