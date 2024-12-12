export default interface User {
  netId: string;
  name: string;
  systemRole: SystemRole;
}

export enum SystemRole {
  ADMIN = "ADMIN",
  FACULTY = "FACULTY",
  STUDENT = "STUDENT",
}
