import { Center, Spinner } from "@chakra-ui/react";

export default function CenteredLoader() {
  return (
    <Center h="100vh">
      <Center w="50%" h="50%">
        <Spinner size="xl" />
      </Center>
    </Center>
  );
}
