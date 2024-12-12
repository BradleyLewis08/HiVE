import express from "express";
import passport from "passport";
import { Strategy } from "passport-cas";

const CAS_VERSION = "CAS2.0";
const CAS_URL = "https://secure.its.yale.edu/cas";

export interface User {
  netId: string;
}

export default function initPassport(app: express.Express) {
  passport.use(
    new Strategy(
      {
        version: CAS_VERSION,
        ssoBaseURL: CAS_URL,
      },
      function (profile: any, done: any) {
        done(null, {
          netId: profile.user,
        });
      }
    )
  );

  passport.serializeUser<User>(function (user: any, done) {
    done(null, user.netId);
  });

  passport.deserializeUser(function (netId, done) {
    done(null, {
      netId,
    });
  });

  app.use(passport.initialize());
  app.use(passport.session());
}
