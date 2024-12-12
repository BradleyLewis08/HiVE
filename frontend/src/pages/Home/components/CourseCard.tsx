import { Text, Card, CardBody, Heading, Stack } from "@chakra-ui/react";
import Course from "../../../models/Course";
import { useNavigate } from "react-router-dom";

interface CourseCardProps {
  course: Course;
}

export default function CourseCard({ course }: CourseCardProps) {
  const navigate = useNavigate();

  const handleCourseClick = () => {
    navigate(`/course/${course.courseCode}`);
  };

  return (
    <Card
      maxW="sm"
      onClick={handleCourseClick}
      _hover={{
        cursor: "pointer",
        bg: "gray.100",
      }}
    >
      <CardBody>
        <Stack spacing="3">
          <Heading size="md">{course.courseCode}</Heading>
          <Heading size="md">{course.name}</Heading>
          <Text>{course.description}</Text>
        </Stack>
      </CardBody>
    </Card>
  );
}
