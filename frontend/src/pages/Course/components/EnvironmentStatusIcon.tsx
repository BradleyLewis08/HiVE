import type { EnvironmentStatus } from "../../../hooks/useAssignmentStatus";
import { Badge, Tooltip } from "@chakra-ui/react";

export interface EnvironmentStatusIconProps {
  status: EnvironmentStatus;
}

const statusProps = {
  pending: { colorScheme: "yellow", label: "Initializing" },
  available: { colorScheme: "green", label: "Available" },
};

export const EnvironmentStatusIcon = ({
  status,
}: EnvironmentStatusIconProps) => {
  return (
    <Tooltip
      label={
        status === "pending"
          ? "Environment is being initialized"
          : "Environment is ready"
      }
    >
      <Badge
        colorScheme={
          status === "available"
            ? statusProps.available.colorScheme
            : statusProps.pending.colorScheme
        }
        variant="subtle"
        px={2}
        py={0.5}
      >
        {status === "available"
          ? statusProps.available.label
          : statusProps.pending.label}
      </Badge>
    </Tooltip>
  );
};
