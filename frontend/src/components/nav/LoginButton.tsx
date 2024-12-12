import { Button } from "@chakra-ui/react";

const LOGIN_REDIRECT = `${process.env.REACT_APP_CLIENT_URL}/redirect`;
const LOGIN_URL = `${
  process.env.REACT_APP_USER_SERVICE_URL
}/auth/login?redirect=${encodeURIComponent(LOGIN_REDIRECT)}`;

export default function LoginButton() {
  const handleLogin = () => {
    window.location.href = LOGIN_URL;
  };

  return (
    <Button
      display={{ base: "none", md: "inline-flex" }}
      fontSize={"sm"}
      fontWeight={600}
      color={"white"}
      bg={"yellow.400"}
      _hover={{
        bg: "blue.300",
      }}
      onClick={handleLogin}
    >
      Login
    </Button>
  );
}
