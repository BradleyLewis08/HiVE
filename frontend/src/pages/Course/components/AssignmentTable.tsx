import Assignment from "../../../models/Assignment";
import {
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  TableContainer,
  IconButton,
} from "@chakra-ui/react";
import { FaTrash } from "react-icons/fa";
import { EnvironmentStatusIcon } from "./EnvironmentStatusIcon";
import { formatDate } from "../../../utils";
import useAssignmentStatus from "../../../hooks/useAssignmentStatus";

export interface AssignmentTableProps {
  assignments: Assignment[];
  onDeleteAssignment: (assignmentId: string) => void;
}

export default function AssignmentTable({
  assignments,
  onDeleteAssignment,
}: AssignmentTableProps) {
  // const statuses = useAssignmentStatus(assignments);
  return (
    <TableContainer>
      <Table variant="simple">
        <Thead>
          <Tr>
            <Th>Assignment</Th>
            <Th>Created</Th>
            <Th>Status</Th>
            <Th>Open</Th>
            <Th>Actions</Th>
          </Tr>
        </Thead>
        <Tbody>
          {assignments.map((assignment) => (
            <Tr key={assignment.id}>
              <Td>{assignment.title}</Td>
              <Td>{formatDate(assignment.createdAt)}</Td>
              {/* <Td>
                <EnvironmentStatusIcon status={statuses[assignment.id]} />
              </Td> */}
              <Td>
                <a
                  href={assignment.environmentURL}
                  target="_blank"
                  rel="noreferrer"
                >
                  {assignment.environmentURL}
                </a>
              </Td>
              <Td>
                <IconButton
                  aria-label="Delete assignment"
                  icon={<FaTrash size={16} />}
                  colorScheme="red"
                  variant="ghost"
                  size="sm"
                  onClick={() => onDeleteAssignment(assignment.id)}
                />
              </Td>
            </Tr>
          ))}
        </Tbody>
      </Table>
    </TableContainer>
  );
}
