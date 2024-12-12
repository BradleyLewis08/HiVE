import { useEffect, useState } from "react";
import {
  Box,
  Container,
  Heading,
  Button,
  Text,
  SimpleGrid,
  VStack,
  HStack,
  useToast,
  Skeleton,
  Icon,
  Badge,
  useColorModeValue,
} from "@chakra-ui/react";
import { AddIcon } from "@chakra-ui/icons";
import { FaChalkboardTeacher, FaUserGraduate } from "react-icons/fa";
import APIClient from "../../api/APIClient";
import { useAuth } from "../../providers/AuthProvider";
import Course from "../../models/Course";
import { SystemRole } from "../../models/User";
import CourseCard from "./components/CourseCard";

export default function Home() {
  const { user } = useAuth();
  const [courses, setCourses] = useState<Course[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const toast = useToast();

  const bgGradient = useColorModeValue(
    "linear(to-br, blue.50, purple.50)",
    "linear(to-br, gray.800, purple.900)"
  );

  const cardBg = useColorModeValue("white", "gray.800");

  const fetchCourses = async () => {
    if (!user) return;

    setIsLoading(true);
    try {
      const client = new APIClient();
      const courses = await client.getCourses(user.netId);
      setCourses(courses);
    } catch (error) {
      console.error("Error fetching courses:", error);
      toast({
        title: "Error fetching courses",
        status: "error",
        duration: 3000,
        isClosable: true,
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateCourse = async () => {
    try {
      const client = new APIClient();
      await client.createCourse(
        user!.netId,
        "CPSC490",
        "Senior Project",
        "A course for senior students"
      );
      await fetchCourses();
      toast({
        title: "Course created successfully",
        status: "success",
        duration: 3000,
        isClosable: true,
      });
    } catch (error) {
      console.error("Error creating course:", error);
      toast({
        title: "Error creating course",
        status: "error",
        duration: 3000,
        isClosable: true,
      });
    }
  };

  useEffect(() => {
    fetchCourses();
  }, [user]);

  if (!user) {
    return null;
  }

  return (
    <Box minH="100vh" bgGradient={bgGradient} py={8}>
      <Container maxW="container.xl">
        <VStack spacing={8} align="stretch">
          {/* Header Section */}
          <Box
            bg={cardBg}
            p={6}
            rounded="xl"
            shadow="md"
            border="1px"
            borderColor="gray.200"
          >
            <VStack align="start" spacing={4}>
              <HStack spacing={4}>
                <Icon
                  as={
                    user.systemRole === SystemRole.FACULTY
                      ? FaChalkboardTeacher
                      : FaUserGraduate
                  }
                  boxSize={8}
                  color="purple.500"
                />
                <VStack align="start" spacing={1}>
                  <Heading size="lg">Welcome, {user.name}</Heading>
                  <Badge
                    colorScheme="purple"
                    fontSize="md"
                    px={3}
                    py={1}
                    rounded="full"
                  >
                    {user.systemRole}
                  </Badge>
                </VStack>
              </HStack>
            </VStack>
          </Box>

          {/* Courses Section */}
          <Box>
            <HStack justify="space-between" mb={6}>
              <Heading size="md">Your Courses</Heading>
              {user.systemRole === SystemRole.FACULTY && (
                <Button
                  leftIcon={<AddIcon />}
                  colorScheme="purple"
                  onClick={handleCreateCourse}
                  size="sm"
                >
                  Create Course
                </Button>
              )}
            </HStack>

            {isLoading ? (
              <SimpleGrid columns={{ base: 1, md: 2, lg: 3 }} spacing={6}>
                {[1, 2, 3].map((i) => (
                  <Skeleton key={i} height="200px" rounded="lg" />
                ))}
              </SimpleGrid>
            ) : courses.length === 0 ? (
              <Box
                p={8}
                bg={cardBg}
                rounded="xl"
                shadow="sm"
                textAlign="center"
                border="1px dashed"
                borderColor="gray.200"
              >
                <Text fontSize="lg" color="gray.600">
                  {user.systemRole === SystemRole.FACULTY
                    ? "Create your first course to get started!"
                    : "You haven't been enrolled in any courses yet."}
                </Text>
              </Box>
            ) : (
              <SimpleGrid columns={{ base: 1, md: 2, lg: 3 }} spacing={6}>
                {courses.map((course) => (
                  <CourseCard key={course.id} course={course} />
                ))}
              </SimpleGrid>
            )}
          </Box>
        </VStack>
      </Container>
    </Box>
  );
}
