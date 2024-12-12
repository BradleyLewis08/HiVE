import axios, { AxiosInstance } from "axios";
import type User from "../models/User";
import { SystemRole } from "../models/User";

const BASE_URL = "http://localhost:8080/api";

export default class APIClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: BASE_URL,
      headers: {
        "Content-Type": "application/json",
      },
    });
    this.client.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response) {
          // The request was made and the server responded with a status code
          // that falls out of the range of 2xx
          console.error("Response error:", error.response.data);
        } else if (error.request) {
          // The request was made but no response was received
          console.error("Request error:", error.request);
        } else {
          // Something happened in setting up the request that triggered an Error
          console.error("Error:", error.message);
        }
        return Promise.reject(error);
      }
    );
  }

  async getUser(netid: string): Promise<User> {
    try {
      const response = await this.client.get(`/user/${netid}`);
      return response.data;
    } catch (error) {
      console.error("Error getting user:", error);
      throw error;
    }
  }

  async updateUserPermissions(
    netId: string,
    systemRole: SystemRole
  ): Promise<User> {
    try {
      const response = await this.client.post(`/user/permissions`, {
        netId,
        systemRole,
      });
      return response.data;
    } catch (error) {
      console.error("Error updating user permissions:", error);
      throw error;
    }
  }

  async createCourse(
    netId: string,
    courseCode: string,
    name: string,
    description: string
  ) {
    try {
      const response = await this.client.post(`/course/create`, {
        netId,
        courseCode,
        name,
        description,
      });
      return response.data;
    } catch (error) {
      console.error("Error creating course:", error);
      throw error;
    }
  }

  async getCourses(netId: string) {
    try {
      const response = await this.client.post(`/course/getCourses`, {
        netId,
      });

      return response.data;
    } catch (error) {
      console.error("Error getting courses:", error);
      throw error;
    }
  }

  async getCourseAssignments(courseCode: string) {
    try {
      const response = await this.client.get(
        `/course/${courseCode}/assignments`
      );
      return response.data;
    } catch (error) {
      console.error("Error getting course assignments:", error);
      throw error;
    }
  }

  async createAssignment(
    courseCode: string,
    title: string,
    imageURL: string,
    netIDs: string[],
    description: string
  ) {
    try {
      const response = await this.client.post(
        `/course/${courseCode}/assignment`,
        {
          title,
          netIDs,
          description,
          imageURL,
        }
      );
      return response.data;
    } catch (error) {
      console.error("Error creating assignment:", error);
      throw error;
    }
  }

  async deleteAssignment(courseCode: string, assignmentId: string) {
    console.log("Deleting assignment", assignmentId);
    const response = await this.client.post(
      `/course/${courseCode}/assignment/delete`,
      {
        assignmentId,
      }
    );
    return response.data;
  }
}
