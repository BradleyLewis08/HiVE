import { useEffect, useState } from "react";
import Assignment from "../models/Assignment";

export type EnvironmentStatus = "pending" | "available";

interface StatusState {
  [key: string]: EnvironmentStatus;
}

const TIMEOUT = 5000;

async function checkEndpointAvailability(endpoint: string) {
  console.log("Checking endpoint availability", endpoint);
  try {
    const response = await fetch(endpoint, {
      method: "HEAD",
      mode: "no-cors",
      signal: AbortSignal.timeout(TIMEOUT),
    });
    console.log(response);
    return true;
  } catch (error) {
    console.log("Error checking endpoint availability for", endpoint, error);
    return false;
  }
}

export default function useAssignmentStatus(assignments: Assignment[]) {
  const [statuses, setStatuses] = useState<StatusState>({});

  useEffect(() => {
    const pollingIntervals: { [key: string]: NodeJS.Timeout } = {};

    assignments.forEach((assignment) => {
      if (assignment.environmentURL) {
        console.log("Setting up polling for", assignment.environmentURL);
        const pollStatus = async () => {
          const isAvailable = await checkEndpointAvailability(
            assignment.environmentURL
          );
          setStatuses((prev) => ({
            ...prev,
            [assignment.id]: isAvailable ? "available" : "pending",
          }));
        };
        // Initial check
        pollStatus();
        // Set up interval for continuous checking
        pollingIntervals[assignment.id] = setInterval(pollStatus, 10000); // Check every 10 seconds
      }
    });

    // Cleanup intervals on unmount or when assignments change
    return () => {
      Object.values(pollingIntervals).forEach((interval) =>
        clearInterval(interval)
      );
    };
  }, [assignments]);

  return statuses;
}
