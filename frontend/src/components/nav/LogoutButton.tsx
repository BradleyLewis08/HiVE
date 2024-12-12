import { Button } from "@chakra-ui/react";
import { useAuth } from "../../providers/AuthProvider";

export default function LogoutButton() {
  const { handleLogout } = useAuth();

  return (
    <Button
      display={{ base: "none", md: "inline-flex" }}
      fontSize={"sm"}
      fontWeight={600}
      color={"white"}
      bg={"blue.400"}
      _hover={{
        bg: "blue.300",
      }}
      onClick={handleLogout}
    >
      Logout
    </Button>
  );
}
