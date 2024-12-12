import React from "react";
import {
  Box,
  Flex,
  Text,
  Button,
  Stack,
  Link,
  useColorModeValue,
  useDisclosure,
  IconButton,
  Collapse,
  Select,
} from "@chakra-ui/react";
import { HamburgerIcon, CloseIcon } from "@chakra-ui/icons";
import { useAuth } from "../../providers/AuthProvider";
import LoginButton from "./LoginButton";
import LogoutButton from "./LogoutButton";
import { SystemRole } from "../../models/User";

const NavLink = ({
  children,
  href,
}: {
  children: React.ReactNode;
  href: any;
}) => (
  <Link
    px={2}
    py={1}
    rounded={"md"}
    _hover={{
      textDecoration: "none",
      bg: useColorModeValue("gray.200", "gray.700"),
    }}
    href={href}
  >
    {children}
  </Link>
);

const StickyNavbar = () => {
  const { user, isAuthenticated, handleRoleChange } = useAuth();
  const { isOpen, onToggle } = useDisclosure();

  // Example navigation items
  return (
    <Box
      position="sticky"
      top={0}
      zIndex="sticky"
      bg={useColorModeValue("white", "gray.800")}
      borderBottom={1}
      borderStyle={"solid"}
      borderColor={useColorModeValue("gray.200", "gray.900")}
      px={4}
    >
      <Flex
        minH={"60px"}
        py={{ base: 2 }}
        px={{ base: 4 }}
        align={"center"}
        justify={"space-between"}
      >
        <Flex
          flex={{ base: 1, md: "auto" }}
          ml={{ base: -2 }}
          display={{ base: "flex", md: "none" }}
        >
          <IconButton
            onClick={onToggle}
            icon={
              isOpen ? <CloseIcon w={3} h={3} /> : <HamburgerIcon w={5} h={5} />
            }
            variant={"ghost"}
            aria-label={"Toggle Navigation"}
          />
        </Flex>

        <Text
          textAlign={"left"}
          fontFamily={"heading"}
          color={useColorModeValue("gray.800", "white")}
          fontWeight="bold"
        >
          Logo
        </Text>
        {user && isAuthenticated && (
          <Select
            display={{ base: "none", md: "flex" }}
            size="sm"
            w="auto"
            onChange={(e) => {
              const systemRole = e.target.value as SystemRole;
              handleRoleChange(systemRole);
            }}
          >
            {user.systemRole !== SystemRole.ADMIN && (
              <option value={SystemRole.ADMIN}>{SystemRole.ADMIN}</option>
            )}
            {user.systemRole !== SystemRole.FACULTY && (
              <option value={SystemRole.FACULTY}>{SystemRole.FACULTY}</option>
            )}
            {user.systemRole !== SystemRole.STUDENT && (
              <option value={SystemRole.STUDENT}>{SystemRole.STUDENT}</option>
            )}
          </Select>
        )}
        <Stack
          flex={{ base: 1, md: 0 }}
          justify={"flex-end"}
          direction={"row"}
          spacing={6}
        >
          {isAuthenticated ? <LogoutButton /> : <LoginButton />}
        </Stack>
      </Flex>
    </Box>
  );
};

export default StickyNavbar;
