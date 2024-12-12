import PrismaService from "./PrismaService";

const PROVISIONER_ENDPOINT = "http://localhost:8000/environment";

// type EnvironmentRequest struct {
// 	CourseName string `json:"courseName"`
// 	NetIDs   []string `json:"netIDs"`
// 	Image   string   `json:"image"`
// }

export interface EnvironmentCreationResponse {
  baseURL: string;
}

export default class EnvironmentService extends PrismaService {
  public async createEnvironment(
    assignmentName: string,
    courseName: string,
    imageURL: string,
    netIDs: string[]
  ) {
    // TODO: Don't hardcode NetID
    const body = {
      AssignmentName: assignmentName,
      CourseName: courseName,
      NetIDs: netIDs,
      Image: imageURL,
    };

    try {
      const response = await fetch(PROVISIONER_ENDPOINT, {
        method: "POST",
        body: JSON.stringify(body),
      });

      if (!response.ok) {
        throw new Error(
          `Failed to create environment: ${response.status} ${response.statusText}`
        );
      }

      const data = (await response.json()) as EnvironmentCreationResponse;

      return data;
    } catch (error) {
      throw new Error(`Failed to create environment: ${error}`);
    }
  }
}
