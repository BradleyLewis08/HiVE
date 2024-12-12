import { useCallback, useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import {
  Box,
  Container,
  Heading,
  Button,
  VStack,
  HStack,
  Text,
  useToast,
  Card,
  CardHeader,
  CardBody,
  Badge,
  Skeleton,
  Alert,
  AlertIcon,
  IconButton,
  useDisclosure,
} from "@chakra-ui/react";
import { AddIcon, RepeatIcon } from "@chakra-ui/icons";
import APIClient from "../../api/APIClient";
import Assignment from "../../models/Assignment";
import AssignmentTable from "./components/AssignmentTable";
import CreateAssignmentModal from "./components/CreateAssignmentModal";
import { getFullURL } from "../../utils";

export interface AssignmentFormData {
  title: string;
  description: string;
  imageURL: string;
}

export default function Course() {
  const [assignments, setAssignments] = useState<Assignment[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [formData, setFormData] = useState<AssignmentFormData>({
    title: "",
    description: "",
    imageURL: "",
  });
  const { isOpen, onOpen, onClose } = useDisclosure();
  const { courseCode } = useParams<{ courseCode: string }>();
  const toast = useToast();

  const onDeleteAssignment = async (assignmentId: string) => {
    if (!courseCode) {
      return;
    }
    const client = new APIClient();
    await client.deleteAssignment(courseCode, assignmentId);
    await fetchCourseAssignments();
  };

  const fetchCourseAssignments = useCallback(async () => {
    if (!courseCode) {
      console.log("No course code found");
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const client = new APIClient();
      let assignments = await client.getCourseAssignments(courseCode);
      assignments = assignments.map((assignment: any) => {
        return {
          ...assignment,
          environmentURL: getFullURL(assignment.Environment[0]?.baseUrl ?? ""),
        };
      });
      setAssignments(assignments);
    } catch (error) {
      setError("Failed to fetch assignments. Please try again later.");
      console.error("Error fetching assignments:", error);
    } finally {
      setIsLoading(false);
    }
  }, [courseCode]);

  const handleCreateAssignment = async () => {
    if (!courseCode) return;

    try {
      const client = new APIClient();
      await client.createAssignment(
        courseCode,
        formData.title,
        formData.imageURL,
        ["bel25"],
        formData.description
      );
      await fetchCourseAssignments();
      toast({
        title: "Assignment created",
        description: "New assignment has been successfully created",
        status: "success",
        duration: 3000,
        isClosable: true,
      });

      // Reset form and close modal
      setFormData({ title: "", description: "", imageURL: "" });
      onClose();
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to create new assignment",
        status: "error",
        duration: 3000,
        isClosable: true,
      });
    }
  };

  const handleInputChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  useEffect(() => {
    fetchCourseAssignments();
  }, [courseCode, fetchCourseAssignments]);

  return (
    <Container maxW="container.xl" py={8}>
      <VStack spacing={6} align="stretch">
        {/* Header Section */}
        <Box bg="white" p={6} rounded="lg" shadow="sm">
          <HStack justify="space-between" align="center">
            <VStack align="start" spacing={1}>
              <Heading size="lg" color="gray.700">
                Course: {courseCode}
              </Heading>
              <Badge colorScheme="blue" fontSize="sm">
                {assignments.length} Assignment
                {assignments.length !== 1 ? "s" : ""}
              </Badge>
            </VStack>
            <HStack spacing={3}>
              <IconButton
                aria-label="Refresh assignments"
                icon={<RepeatIcon />}
                onClick={fetchCourseAssignments}
                size="sm"
              />
              <Button
                leftIcon={<AddIcon />}
                colorScheme="blue"
                onClick={onOpen}
              >
                Create Assignment
              </Button>
            </HStack>
          </HStack>
        </Box>

        <Card>
          <CardHeader>
            <Heading size="md">Assignments</Heading>
          </CardHeader>
          <CardBody>
            {error && (
              <Alert status="error" mb={4}>
                <AlertIcon />
                {error}
              </Alert>
            )}
            {isLoading ? (
              <VStack spacing={4}>
                <Skeleton height="50px" width="100%" />
                <Skeleton height="50px" width="100%" />
                <Skeleton height="50px" width="100%" />
              </VStack>
            ) : assignments.length === 0 ? (
              <Text color="gray.500" textAlign="center" py={8}>
                No assignments found. Create your first assignment to get
                started.
              </Text>
            ) : (
              <AssignmentTable
                assignments={assignments}
                onDeleteAssignment={onDeleteAssignment}
              />
            )}
          </CardBody>
        </Card>
      </VStack>
      <CreateAssignmentModal
        isOpen={isOpen}
        onClose={onClose}
        onInputChange={handleInputChange}
        formData={formData}
        handleCreateAssignment={handleCreateAssignment}
      />
    </Container>
  );
}
